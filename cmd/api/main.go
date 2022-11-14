package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ArkeoNetwork/directory/api"
	"github.com/ArkeoNetwork/directory/internal/logging"
	"github.com/ArkeoNetwork/directory/pkg/db"
)

var log = logging.WithoutFields()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Info("starting api")
	// TODO determine config mechanism
	api := api.NewApiService(api.ApiServiceParams{
		DBConfig: db.DBConfig{
			Host:         "localhost",
			Port:         5432,
			User:         "arkeo",
			Pass:         "arkeo123",
			DBName:       "arkeo_directory",
			PoolMaxConns: 2,
			PoolMinConns: 1,
			SSLMode:      "prefer",
		},
	})
	done, err := api.Start()
	if err != nil {
		panic(fmt.Sprintf("error starting api service: %+v", err))
	}
	<-done
	log.Info("api complete")
}
