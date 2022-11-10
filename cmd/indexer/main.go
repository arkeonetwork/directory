package main

import (
	"math/rand"
	"time"

	"github.com/ArkeoNetwork/directory/internal/logging"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var log = logging.WithoutFields()
	log.Info("start the indexer")
}
