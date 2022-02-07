package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-pkgz/lgr"
	"github.com/gorilla/mux"
	"github.com/parMaster/ethtrx/internal/app/model"
	"github.com/parMaster/ethtrx/internal/app/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	logger *lgr.Logger
	router *mux.Router
	client *Client
	store  store.Storer
	cfg    Config
}

func newServer(store store.Storer, config Config) *server {
	s := &server{
		logger: lgr.New(lgr.Debug, lgr.Option(lgr.Secret(config.APIKey))),
		router: mux.NewRouter(),
		client: newClient(config.APIURL, config.APIKey),
		store:  store,
		cfg:    config,
	}
	s.configureRouter()

	s.startUp(config.LoadTransactions)

	if s.cfg.DaemonMode {
		go s.catchUp()
	}

	return s
}

// Identify the "empty db" scenario
// put current-offset block to have something to start with
func (s *server) startUp(offset int) {
	_, err := s.store.MostRecentBlock()
	if err != nil && err.Error() == "record not found" {
		s.logger.Logf("No transactions in DB. Initializing...")

		currentBlock, err := s.getCurrentBlockNumber()
		if err != nil {
			s.logger.Logf("FATAL Can't start up: %s", err.Error())
		}
		s.logger.Logf("INFO Current block: %s", currentBlock)

		currentBlockNumber, err := hexutil.DecodeUint64(currentBlock)
		if err != nil {
			s.logger.Logf("FATAL Can't decode block number: %s", err.Error())
		}

		startingBlock, err := s.getBlockByNumber(hexutil.EncodeUint64(currentBlockNumber-uint64(offset)), true)
		if err != nil {
			s.logger.Logf("FATAL Can't get starting block: %s", err.Error())
		}
		s.logger.Logf("INFO Starting with block: %s", startingBlock.Number)

		for _, trx := range startingBlock.Transactions {
			trx.BlockTime = startingBlock.Timestamp
			trx.BlockHeight = hexutil.MustDecodeUint64(trx.BlockNumber)
			s.store.Create(&trx)
		}
	}
}

// Cyclic process of getting current block and saving transactions
func (s *server) catchUp() {

	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {

		startingBlock, err := s.store.MostRecentBlock()
		if err != nil {
			s.logger.Logf("WARN Failed to find any blocks in DB: %s", err.Error())
			continue
		}

		currentBlock, err := s.getCurrentBlockNumber()
		if err != nil {
			s.logger.Logf("WARN Failed to get last block number: %s", err.Error())
			continue
		}

		if startingBlock == currentBlock {
			continue
		}

		var block *model.Block
		for i := hexutil.MustDecodeUint64(startingBlock) + 1; i <= hexutil.MustDecodeUint64(currentBlock); i++ {
			block, err = s.getBlockByNumber(hexutil.EncodeUint64(i), true)
			if err != nil {
				s.logger.Logf("WARN Failed to get block: %s", err.Error())
				continue
			}

			s.logger.Logf("INFO Saving new block %d (%s)", i, hexutil.EncodeUint64(i))
			for _, trx := range block.Transactions {
				trx.BlockTime = block.Timestamp
				trx.BlockHeight = hexutil.MustDecodeUint64(trx.BlockNumber)
				s.store.Create(&trx)
			}
			time.Sleep(200 * time.Millisecond)
		}

	}
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/getTxList", s.handleGetTxList()).Methods("POST")
}

// Handles getTxList endpoint
func (s *server) handleGetTxList() http.HandlerFunc {
	type request struct {
		TxHash      *string `json:"txhash,omitempty"`
		From        *string `json:"from,omitempty"`
		To          *string `json:"to,omitempty"`
		BlockHeight *int    `json:"blockheight,omitempty"`
		BlockDate   *string `json:"date,omitempty"`
		Page        *int    `json:"page,omitempty"`
		Limit       *int    `json:"limit,omitempty"`
	}

	type response struct {
		Transactions []model.Transaction `json:"transactions,omitempty"`
		Error        string              `json:"error,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		response := &response{}

		s.logger.Logf("%+v", r.Body)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil && err.Error() != "EOF" {
			response.Error = fmt.Sprintf("Error decoding request: %s", err.Error())
			s.respond(w, r, http.StatusBadRequest, response)
			return
		}

		var transactions []model.Transaction

		filter := bson.M{}
		if req.TxHash != nil && len(*req.TxHash) > 0 {
			_, err := hexutil.Decode(*req.TxHash)
			if err != nil {
				response.Error = fmt.Sprintf("Invalid 'txhash' parameter: %s", err.Error())
				s.respond(w, r, http.StatusBadRequest, response)
				return
			}
			filter["hash"] = req.TxHash
		}

		if req.From != nil {
			if !common.IsHexAddress(*req.From) {
				response.Error = "'from' parameter is not a valid hex address"
				s.respond(w, r, http.StatusBadRequest, response)
				return
			}
			filter["from"] = req.From
		}

		if req.To != nil {
			if !common.IsHexAddress(*req.To) {
				response.Error = "'to' parameter is not a valid hex address"
				s.respond(w, r, http.StatusBadRequest, response)
				return
			}
			filter["to"] = req.To
		}

		if req.BlockHeight != nil {
			filter["blockNumber"] = hexutil.EncodeUint64(uint64(*req.BlockHeight))
		}

		var page, limit int64
		if req.Page != nil &&
			req.Limit != nil &&
			int64(*req.Page) >= 0 &&
			int64(*req.Limit) > 0 {

			page = int64(*req.Page)
			limit = int64(*req.Limit)
		} else {
			// Default values - page = 1, limit = 10
			page = int64(1)
			limit = int64(10)
		}

		if req.BlockDate != nil {
			// filter["blockDate"] = Validate(*req.BlockDate)
		}

		transactions, err := s.store.FindTx(filter, page, limit)
		if err != nil {
			s.logger.Logf("ERROR finding transactions: %s", err.Error())
			s.respond(w, r, http.StatusInternalServerError, nil)
			return
		}
		if len(transactions) == 0 {
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		mostRecentBlock, err := s.store.MostRecentBlock()
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil)
			return
		}
		mostRecentBlockNumber := hexutil.MustDecodeUint64(mostRecentBlock)

		// Preparing output data
		conv := BigHexToStr()
		for k, v := range transactions {

			// Calculating Confirmations instead of storing and updating
			transactions[k].Confirmations = mostRecentBlockNumber - hexutil.MustDecodeUint64(v.BlockNumber)

			// Same with human-readable values: providing both fields - raw and string with human-readable float
			transactions[k].ValueNumber = conv(v.Value)
			transactions[k].BlockHeight = hexutil.MustDecodeUint64(v.BlockNumber)
		}

		s.respond(w, r, http.StatusOK, transactions)
	}
}

// Returns current Eth block number (hex string)
func (s *server) getCurrentBlockNumber() (string, error) {

	type BlockNumber struct {
		Id     int64  `json:"id"`
		Result string `json:"result"`
	}

	var record BlockNumber

	_, err := s.client.get(map[string]string{
		"module": "proxy",
		"action": "eth_blockNumber",
	}, &record)
	if err != nil {
		s.logger.Logf("ERROR calling eth_blockNumber %s", err.Error())
	}

	return record.Result, nil
}

// Returns information about a transaction by hash
func (s *server) getTxByHash(txhash string, saveUpdate bool) (*model.TransactionResponse, error) {

	var record model.TransactionResponse

	_, err := s.client.get(map[string]string{
		"module":  "proxy",
		"action":  "eth_getTransactionByHash",
		"boolean": "true",
		"txhash":  txhash,
	}, &record)
	if err != nil {
		s.logger.Logf("ERROR calling eth_getTransactionByHash %s", err.Error())
		return nil, err
	}

	if saveUpdate {
		_, err := s.store.Find(record.Hash)
		if err == mongo.ErrNoDocuments {
			s.store.Create(&record.Transaction)
		} else {
			s.store.Update(&record.Transaction)
		}
	}

	return &record, nil
}

// Returns information about a block by block number
// full - the boolean value to show full transaction objects.
func (s *server) getBlockByNumber(blockNumber string, full bool) (*model.Block, error) {
	var record model.BlockResponse

	_, err := s.client.get(map[string]string{
		"module":  "proxy",
		"action":  "eth_getBlockByNumber",
		"boolean": strconv.FormatBool(full),
		"tag":     blockNumber,
	}, &record)
	if err != nil {
		s.logger.Logf("ERROR calling eth_getBlockByNumber %s", err.Error())
		return nil, err
	}

	return &record.Block, nil
}

// Maintain compatibility with Handler interface
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
