package indexer

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/ArkeoNetwork/directory/pkg/db"
	"github.com/ArkeoNetwork/directory/pkg/logging"
	"github.com/pkg/errors"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
)

var log = logging.WithoutFields()

type IndexerAppParams struct {
	ArkeoApi            string
	TendermintApi       string
	TendermintWs        string
	ChainID             string
	Bech32PrefixAccAddr string
	Bech32PrefixAccPub  string
	IndexerID           int64
	db.DBConfig
}

type IndexerApp struct {
	Height   uint64
	IsSynced atomic.Bool
	params   IndexerAppParams
	db       *db.DirectoryDB
	done     chan struct{}
}

func NewIndexer(params IndexerAppParams) *IndexerApp {
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

func makeTMClient(baseURL string) (*tmclient.HTTP, error) {
	client, err := tmclient.New(baseURL, "/websocket")
	if err != nil {
		return nil, errors.Wrapf(err, "error creating websocket client")
	}
	logger := tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
	client.SetLogger(logger)

	return client, nil
}

const numClients = 3

func (a *IndexerApp) start() {
	log.Infof("starting realtime indexing using /websocket at %s", a.params.TendermintWs)
	clients := make([]*tmclient.HTTP, numClients)
	for i := 0; i < numClients; i++ {
		client, err := makeTMClient(a.params.TendermintWs)
		if err != nil {
			panic(fmt.Sprintf("error creating tm client for %s: %+v", a.params.TendermintWs, err))
		}
		if err = client.Start(); err != nil {
			panic(fmt.Sprintf("error starting ws client: %s: %+v", a.params.TendermintWs, err))
		}
		defer client.Stop()
		clients[i] = client
	}

	// determine last seen block from db
	indexerStatus, err := a.db.FindIndexerStatus(a.params.IndexerID)
	if err != nil {
		panic(fmt.Sprintf("error getting indexer state from db: %+v", err))
	}

	a.IsSynced.Store(false)
	if indexerStatus == nil {
		// this is the first time we have started the indexer with this db, sync from heigh of 1 (0 will fail)
		indexerStatus = &db.IndexerStatus{ID: 0, Height: 0}
		if _, err = a.db.UpdateIndexerStatus(indexerStatus); err != nil {
			log.Errorf("error writing initial indexer status: %+v", err)
		}
	}
	a.Height = indexerStatus.Height

	go func() {
		// We kick off new thread to background historical sync while we keep consuming events that are coming in real time.
		// this should ensure that we don't miss any events since we will end up overlapping with out real time WS subscriptions
		log.Infof("Starting historical syncing from block height: %d", a.Height)
		err := a.consumeHistoricalEvents(clients[0])
		if err != nil {
			log.Infof("Historical syncing failed")
			panic(fmt.Sprintf("Historical syncing failed: %+v", err))
		}
		log.Infof("Historical syncing completed")
		a.IsSynced.Store(true)
	}()

	a.consumeEvents(clients)
	a.done <- struct{}{}
}

func (a *IndexerApp) handleBlockEvent(height int64) error {
	if !a.IsSynced.Load() {
		return nil // when we are syncing we don't want to update the DB until we are fully up to date.
	}

	indexerStatus := db.IndexerStatus{
		ID:     a.params.IndexerID,
		Height: uint64(height),
	}
	_, err := a.db.UpsertIndexerStatus(&indexerStatus)
	if err != nil {
		return errors.Wrapf(err, "error updating indexer status for %d height %d", indexerStatus.ID, indexerStatus.Height)
	}
	a.Height = uint64(height)
	return nil
}
