package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-pkgz/lgr"
	"github.com/gorilla/mux"
	"github.com/parMaster/ethtrx/internal/app/model"
	"github.com/parMaster/ethtrx/internal/app/store"
)

type server struct {
	logger *lgr.Logger
	router *mux.Router
	client *http.Client
	store  store.Storer
	cfg    Config
}

func newServer(store store.Storer, config Config) *server {
	s := &server{
		logger: lgr.New(lgr.Debug, lgr.Option(lgr.Secret(config.APIKey))),
		router: mux.NewRouter(),
		client: &http.Client{},
		store:  store,
		cfg:    config,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
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
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=%s", s.cfg.APIKey)
	s.logger.Logf("DEBUG Client called: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.logger.Logf("DEBUG NewRequest: %s", url)
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Logf("ERROR %s", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	var record model.BlockNumber

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		s.logger.Logf("ERROR decoding BlockNumber %s", err.Error())
		return nil, err
	}

	s.logger.Logf("DEBUG Block Id: %d", record.Id)
	s.logger.Logf("DEBUG Block Result: %s", record.Result)
	return &record, nil
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
