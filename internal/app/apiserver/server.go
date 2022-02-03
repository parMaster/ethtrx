package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-pkgz/lgr"
	"github.com/gorilla/mux"
	"github.com/parMaster/ethtrx/internal/app/model"
	"github.com/parMaster/ethtrx/internal/app/store"
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
	s.router.HandleFunc("/getTrx", s.handleGetTrx())
	s.router.HandleFunc("/getBlockNumber", s.handleGetBlockNumber())
}

func (s *server) handleGetTrx() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		s.respond(w, r, http.StatusOK, "test")
	}
}

func (s *server) handleGetBlockNumber() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		result, _ := s.getBlockNumber()
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) getBlockNumber() (*model.BlockNumber, error) {
	var record model.BlockNumber

	_, err := s.client.get(map[string]string{
		"module": "proxy",
		"action": "eth_blockNumber",
	}, &record)
	if err != nil {
		s.logger.Logf("ERROR calling eth_blockNumber %s", err.Error())
	}

	s.logger.Logf("DEBUG Block Id: %d", record.Id)
	s.logger.Logf("DEBUG Block Result: %s", record.Result)
	return &record, nil
}

func (s *server) getBlockByNumber(blockNumber string) (*model.BlockInfo, error) {
	var record model.BlockInfo

	_, err := s.client.get(map[string]string{
		"module":  "proxy",
		"action":  "eth_getBlockByNumber",
		"boolean": "true",
		"tag":     blockNumber,
	}, &record)
	if err != nil {
		s.logger.Logf("ERROR calling eth_getBlockByNumber %s", err.Error())
	}

	for k, v := range record.Result.Transactions {
		s.logger.Logf("DEBUG transaction %d: %s", k, v)
	}

	return &record, nil
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
