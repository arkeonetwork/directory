package api

import (
	"net/http"

	"github.com/ArkeoNetwork/directory/internal/logging"
	"github.com/gorilla/mux"
)

type ApiService struct {
	router *mux.Router
	params ApiServiceParams
}

type ApiServiceParams struct {
	ListenAddr string
	// DbHost string
	// DbPort string
	// DbUser string
	// DbPass string
	// DbName string
}

const DefaultListenAddress = "localhost:7777"

var log = logging.WithoutFields()

func NewApiService(params ApiServiceParams) *ApiService {
	if params.ListenAddr == "" {
		params.ListenAddr = DefaultListenAddress
	}
	return &ApiService{params: params, router: buildRouter()}
}

func (a *ApiService) Start() (chan struct{}, error) {
	doneChan := make(chan struct{})
	go a.start(doneChan)
	return doneChan, nil
}

func (a *ApiService) start(doneChan chan struct{}) {
	log.Infof("starting http service on %s", a.params.ListenAddr)
	log.Fatal(http.ListenAndServe(a.params.ListenAddr, a.router))
}

func buildRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", handleHealth).Methods(http.MethodGet)
	router.HandleFunc("/stats", getStats).Methods(http.MethodGet)

	providerRouter := router.PathPrefix("/provider").Subrouter()
	providerRouter.HandleFunc("/{pubkey}", getProvider).Methods(http.MethodGet)
	providerRouter.HandleFunc("/search", searchProviders).Methods(http.MethodGet)

	// router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// 	tpl, _ := route.GetPathTemplate()
	// 	log.Infof("walk: %s", tpl)
	// 	return nil
	// })

	return router
}
