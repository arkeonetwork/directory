package api

import (
	"net/http"
)

type ArkeoStats struct {
	ContractsOpen  int64
	ContractsTotal int64
}

func getStats(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, &ArkeoStats{})
}
