package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ArkeoNetwork/directory/api"
	"github.com/ArkeoNetwork/directory/internal/logging"
)

var log = logging.WithoutFields()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Info("starting api")
	api := api.NewApiService(api.ApiServiceParams{ /* determine config mechanism */ })
	done, err := api.Start()
	if err != nil {
		panic(fmt.Sprintf("error starting api service: %+v", err))
	}
	<-done
	log.Info("api complete")
}
