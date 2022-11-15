package api

import (
	"net/http"
)

// swagger:model ArkeoStats
type ArkeoStats struct {
	ContractsOpen                   int64
	ContractsTotal                  int64
	ContractsMedianDuration         int64
	ContractsMedianRatePayPer       int64
	ContractsMedianRateSubscription int64
	ChainStats                      []*ChainStats // should this be a map instead?
}

// swagger:model ChainStats
type ChainStats struct {
	Chain              string
	ProviderCount      int64
	QueryCount         int64
	QueryCountLastDay  int64
	TotalIncome        int64
	TotalIncomeLastDay int64
}

// swagger:route Get /stats getStatsArkeo
//
// get Arkeo network stats
//
// Responses:
//
//	200: ArkeoStats
//	500: InternalServerError
func getStatsArkeo(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, &ArkeoStats{})
}

// swagger:route Get /stats/{chain} getStatsChain
//
// get chain specific network stats
// Parameters:
//   + name: chain
//     in: path
//     description: chain identifier
//     required: true
//     type: string
//
// Responses:
//
//	200: ChainStats
//	500: InternalServerError

func getStatsChain(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, &ChainStats{})
}
