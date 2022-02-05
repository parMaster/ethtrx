package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-pkgz/lgr"
	"github.com/gorilla/mux"
	"github.com/parMaster/ethtrx/internal/app/model"
	"github.com/parMaster/ethtrx/internal/app/store"
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

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/getTx", s.handleGetTx())
	s.router.HandleFunc("/getBlockNumber", s.handleGetBlockNumber())
}

func (s *server) handleGetTx() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		lastBlockNumber, _ := s.getBlockNumber()

		lastBlock, _ := s.getBlockByNumber(lastBlockNumber)

		for _, v := range lastBlock.Transactions {
			tx, _ := s.getTxByHash(v, true)
			s.logger.Logf(tx.Hash)
		}

		s.respond(w, r, http.StatusOK, "test")
	}
}

func (s *server) handleGetBlockNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		result, _ := s.getBlockNumber()
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) getBlockNumber() (string, error) {

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

	s.logger.Logf("DEBUG Block: %+v", record)
	return record.Result, nil
}

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

func (s *server) getBlockByNumber(blockNumber string) (*model.Block, error) {
	var record model.BlockResponse

	_, err := s.client.get(map[string]string{
		"module":  "proxy",
		"action":  "eth_getBlockByNumber",
		"boolean": "false",
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
