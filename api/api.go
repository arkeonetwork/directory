// Package api Directory API.
// Version: 0.0.1
//
// swagger:meta
package api

import (
	"fmt"
	"net/http"

	"github.com/ArkeoNetwork/directory/internal/logging"
	"github.com/ArkeoNetwork/directory/pkg/db"
	"github.com/gorilla/mux"
)

type ApiService struct {
	router *mux.Router
	params ApiServiceParams
	db     *db.DirectoryDB
}

type ApiServiceParams struct {
	ListenAddr string
	DBConfig   db.DBConfig
}

const DefaultListenAddress = "localhost:7777"

var log = logging.WithoutFields()

func NewApiService(params ApiServiceParams) *ApiService {
	if params.ListenAddr == "" {
		params.ListenAddr = DefaultListenAddress
	}
	database, err := db.New(params.DBConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to instantiate db: %+v", err))
	}
	a := &ApiService{params: params, db: database}
	a.router = buildRouter(a)

	return a
}

func (a *ApiService) Start() (chan struct{}, error) {
	doneChan := make(chan struct{})
	go a.start(doneChan)
	return doneChan, nil
}

func (a *ApiService) start(doneChan chan struct{}) {
	log.Infof("starting http service on %s", a.params.ListenAddr)
	if err := http.ListenAndServe(a.params.ListenAddr, a.router); err != nil {
		log.Errorf("error from http listener: %+v", err)
	}
	doneChan <- struct{}{}
}

func buildRouter(a *ApiService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", handleHealth).Methods(http.MethodGet)
	router.HandleFunc("/stats", getStatsArkeo).Methods(http.MethodGet)
	router.HandleFunc("/stats/{chain}", getStatsChain).Methods(http.MethodGet)

	providerRouter := router.PathPrefix("/provider").Subrouter()
	providerRouter.HandleFunc("/{pubkey}", a.getProvider).Methods(http.MethodGet)
	providerRouter.HandleFunc("/search", a.searchProviders).Methods(http.MethodGet)

	// router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// 	tpl, _ := route.GetPathTemplate()
	// 	log.Infof("walk: %s", tpl)
	// 	return nil
	// })

	return router
}
