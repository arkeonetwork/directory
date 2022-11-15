package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// swagger:model ArkeoProvider
type ArkeoProvider struct {
	Pubkey string
}

type ArkeoProviders []*ArkeoProvider

// swagger:route Get /provider/{pubkey}/ getProvider
//
// Get a specific ArkeoProvider by a unique id (pubkey+chain)
//
// parameters:
//   + name: pubkey
//     in: path
//     description: provider public key
//     required: true
//     type: string
//   + name: chain
//	   in: query
//     description: chain identifier
//     required: true
//     type: string
//
// responses:
//
//	200: ArkeoProvider
//	500: InternalServerError

func (a *ApiService) getProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pubkey := vars["pubkey"]
	chain := r.FormValue("chain")
	if pubkey == "" {
		respondWithError(w, http.StatusBadRequest, "pubkey is required")
		return
	}
	if chain == "" {
		respondWithError(w, http.StatusBadRequest, "chain is required")
		return
	}
	// "bitcoin-mainnet"
	provider, err := a.findProvider(pubkey, chain)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error finding provider with pubkey %s", pubkey))
		return
	}

	respondWithJSON(w, http.StatusOK, provider)
}

// search providers
func (a *ApiService) searchProviders(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, ArkeoProviders{})
}

// just returns provider w/ passed pubkey
func (a *ApiService) findProvider(pubkey, chain string) (*ArkeoProvider, error) {
	dbProvider, err := a.db.FindProvider(pubkey, chain)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding provider for %s %s", pubkey, chain)
	}
	if dbProvider == nil {
		return nil, nil
	}
	provider := &ArkeoProvider{Pubkey: dbProvider.Pubkey}
	return provider, nil
}
