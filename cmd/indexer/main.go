package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ArkeoNetwork/directory/indexer"
	"github.com/ArkeoNetwork/directory/internal/logging"
)

var log = logging.WithoutFields()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Info("starting indexer")
	app := indexer.NewIndexer(indexer.IndexerAppParams{})
	done, err := app.Run()
	if err != nil {
		panic(fmt.Sprintf("error starting indexer: %+v", err))
	}
	<-done
	log.Info("indexer complete")
}
