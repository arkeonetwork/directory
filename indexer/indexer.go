package indexer

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ArkeoNetwork/directory/pkg/db"
	"github.com/ArkeoNetwork/directory/pkg/logging"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var log = logging.WithoutFields()

type IndexerAppParams struct {
	ArkeoApi            string
	TendermintApi       string
	TendermintWs        string
	ChainID             string
	Bech32PrefixAccAddr string
	Bech32PrefixAccPub  string
	db.DBConfig
}

type IndexerApp struct {
	params IndexerAppParams
	db     *db.DirectoryDB
	done   chan struct{}
}

func NewIndexer(params IndexerAppParams) *IndexerApp {
	configure()
	d, err := db.New(params.DBConfig)
	if err != nil {
		panic(fmt.Sprintf("error connecting to the db: %+v", err))
	}
	return &IndexerApp{params: params, db: d}
}

func (a *IndexerApp) Run() (done <-chan struct{}, err error) {
	// initialize by reading all existing providers?
	a.done = make(chan struct{})
	go a.start()
	return a.done, nil
}

func (a *IndexerApp) start() {
	log.Infof("starting realtime indexing using /websocket at %s", a.params.TendermintWs)
	client, err := tmclient.New(a.params.TendermintWs, "/websocket")
	if err != nil {
		log.Errorf("failure to create websocket client: %+v", err)
		panic(err)
	}
	logger := tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
	client.SetLogger(logger)
	err = client.Start()
	if err != nil {
		panic(fmt.Sprintf("error starting client: %+v", err))
	}
	defer client.Stop()
	a.consumeEvents(client)
	a.done <- struct{}{}
}

// TODO: if there are multiple of the same type of event, this may be
// problematic, multiple events may get purged into one (not sure)
func convertEvent(etype string, raw map[string][]string) map[string]string {
	newEvt := make(map[string]string, 0)

	for k, v := range raw {
		if strings.HasPrefix(k, etype+".") {
			parts := strings.SplitN(k, ".", 2)
			newEvt[parts[1]] = v[0]
		}
	}

	return newEvt
}

func subscribe(client *tmclient.HTTP, query string) <-chan ctypes.ResultEvent {
	out, err := client.Subscribe(context.Background(), "", query)
	if err != nil {
		log.Errorf("failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}
	return out
}

var validChains = map[string]struct{}{"arkeo-mainnet": {}, "eth-mainnet": {}, "btc-mainnet": {}}

func validateChain(chain string) (ok bool) {
	_, ok = validChains[chain]
	return
}
