package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ArkeoProvider struct {
	Pubkey string
}

type ArkeoProviders []*ArkeoProvider

// find a provider by unique id (pubkey)
func getProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pubkey := vars["pubkey"]
	provider, err := findProvider(pubkey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error finding provider with pubkey %s", pubkey))
	}

	respondWithJSON(w, http.StatusOK, provider)
}

// search providers
func searchProviders(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, ArkeoProviders{})
}

// just returns provider w/ passed pubkey
func findProvider(pubkey string) (*ArkeoProvider, error) {
	// TODO provider, err := db.findProvider(pubkey)
	return &ArkeoProvider{Pubkey: pubkey}, nil
}
