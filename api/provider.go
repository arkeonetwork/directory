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

// Contains info about a 500 Internal Server Error response
// swagger:model InternalServerError
type InternalServerError struct {
	Message string `json:"message"`
}

// swagger:model ArkeoProviders
type ArkeoProviders []*ArkeoProvider

// swagger:route Get /provider/{pubkey} getProvider
//
// Get a specific ArkeoProvider by a unique id (pubkey+chain)
//
// Parameters:
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
// Responses:
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

// swagger:route Get /search searchProviders
//
// queries the service for a list of providers
//
// Parameters:
//   + name: sort
//     in: query
//     description: defines how to sort the list of providers
//     required: false
//     schema:
//      type: string
//      enum: age, conract_count, amount_paid
//   + name: distance
//     in: query
//     description: maximum distance in kilometers from provided coordinates
//     required: false
//     type: integer
//   + name: coordinates
//	   description: latitude and longitude (required when providing distance filter, example 40.7127837,-74.0059413)
//     in: query
//     required: false
//   + name: min-validator-payments
//	   description: minimum amount the provider has paid to validators
//     in: query
//     required: false
//   + name: min-provider-age
//	   description: minimum age of provider
//     in: query
//     required: false
//     type: integer
//   + name: min-rate-lmit
//	   description: min rate limit of provider in requests per seconds
//     in: query
//     required: false
//	   type: integer
//   + name: min-contracts
//	   description: minimum number of contracts open with proivder
//     in: query
//     required: false
//	   type: integer
// Responses:
//
//	200: ArkeoProviders
//	500: InternalServerError

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
