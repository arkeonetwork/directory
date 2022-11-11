package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ArkeoNetwork/directory/indexer"
	"github.com/ArkeoNetwork/directory/internal/logging"
	"github.com/ArkeoNetwork/directory/pkg/db"
)

var log = logging.WithoutFields()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Info("starting indexer")
	app := indexer.NewIndexer(indexer.IndexerAppParams{
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
	done, err := app.Run()
	if err != nil {
		panic(fmt.Sprintf("error starting indexer: %+v", err))
	}
	<-done
	log.Info("indexer complete")
}
