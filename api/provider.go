package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/ArkeoNetwork/directory/pkg/utils"
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
		log.Errorf("error finding provider for %s chain %s: %+v", pubkey, chain, err)
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
//   + name: maxDistance
//     in: query
//     description: maximum distance in kilometers from provided coordinates
//     required: false
//     type: integer
//   + name: coordinates
//	   description: latitude and longitude (required when providing distance filter, example 40.7127837,-74.0059413)
//     in: query
//     required: false
//     type: string
//   + name: min-validator-payments
//	   description: minimum amount the provider has paid to validators
//     in: query
//     required: false
//	   type: integer
//   + name: min-provider-age
//	   description: minimum age of provider
//     in: query
//     required: false
//     type: integer
//   + name: min-rate-limit
//	   description: min rate limit of provider in requests per seconds
//     in: query
//     required: false
//	   type: integer
//   + name: min-open-contracts
//	   description: minimum number of contracts open with proivder
//     in: query
//     required: false
//	   type: integer
// Responses:
//
//	200: ArkeoProviders
//	500: InternalServerError

func (a *ApiService) searchProviders(response http.ResponseWriter, request *http.Request) {
	sort := request.FormValue("sort")
	maxDistanceInput := request.FormValue("maxDistance")
	coordinatesInput := request.FormValue("coordinates")
	minValidatorPaymentsInput := request.FormValue("min-validator-payments")
	minProviderAgeInput := request.FormValue("min-provider-age")
	minRateLimitInput := request.FormValue("min-rate-limit")
	minOpenContractsInput := request.FormValue("min-open-contracts")

	if sort == "" &&
		maxDistanceInput == "" &&
		coordinatesInput == "" &&
		minValidatorPaymentsInput == "" &&
		minProviderAgeInput == "" &&
		minRateLimitInput == "" &&
		minOpenContractsInput == "" {
		respondWithError(response, http.StatusBadRequest, "must supply atleast one valid search parameter")
		return
	}

	if maxDistanceInput != "" && coordinatesInput == "" || coordinatesInput != "" && maxDistanceInput == "" {
		respondWithError(response, http.StatusBadRequest, "max distance must accompany coordinates when supplied")
		return
	}

	// TODO: create enum for sort keys
	var maxDistance int64
	var coordinates types.Coordinates
	if maxDistanceInput != "" {
		var err error
		maxDistance, err = strconv.ParseInt(maxDistanceInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "max distance can not be parse")
			return
		}
		coordinates, err = utils.ParseCoordinates(coordinatesInput)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "coordinates can not be parsed")
			return
		}
	}

	var minValidatorPayments int64
	if minValidatorPaymentsInput != "" {
		var err error
		minValidatorPayments, err = strconv.ParseInt(minValidatorPaymentsInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-validator-payments can not be parsed")
			return
		}
	}

	var minProviderAge int64
	if minProviderAgeInput != "" {
		var err error
		minProviderAge, err = strconv.ParseInt(minProviderAgeInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-provider-age can not be parsed")
			return
		}
	}

	var minRateLimit int64
	if minRateLimitInput != "" {
		var err error
		minRateLimit, err = strconv.ParseInt(minRateLimitInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-rate-limit can not be parsed")
			return
		}
	}

	var minOpenContracts int64
	if minOpenContractsInput != "" {
		var err error
		minOpenContracts, err = strconv.ParseInt(minOpenContractsInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-rate-limit can not be parsed")
			return
		}
	}

	respondWithJSON(response, http.StatusOK, ArkeoProviders{})
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
