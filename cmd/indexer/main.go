package main

import (
	"github.com/ArkeoNetwork/directory/internal/logging"
)

func main() {
	var log = logging.WithoutFields()
	log.Info("start the indexer")
}
