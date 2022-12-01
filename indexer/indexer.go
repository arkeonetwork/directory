package indexer

import (
	"fmt"
	"os"

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
	IsSynced bool
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
	// determine last seen block from db
	indexerStatus, err := a.db.FindIndexerStatus(a.params.IndexerID)
	if err != nil {
		panic(fmt.Sprintf("error getting indexer state from db: %+v", err))
	}

	a.IsSynced = false
	if indexerStatus == nil {
		// this is the first time we have started the indexer wit this db, sync from heigh of 0
		a.Height = 0
	} else {
		// we have existing records, roll back 60 blocks on startup to ensure we have not missed any events due to crashing, etc.
		var rollbackBlocks uint64 = 60
		if indexerStatus.Height > rollbackBlocks {
			a.Height = indexerStatus.Height - rollbackBlocks
		} else {
			a.Height = 0
		}
	}
	log.Infof("Starting historical syncing from block height: %d", a.Height)
	// unsure how to deal with syncronization here. Ideally we kick off new thread to background historical sync
	// while we keep consuming events that are coming in real time.
	a.consumeHistoricalEvents(client)
	a.IsSynced = true
	a.consumeEvents(client)
	a.done <- struct{}{}
}

func (a *IndexerApp) handleBlockEvent(height int64) error {
	if !a.IsSynced {
		return nil // when we are syncing we don't want to update the DB until we are fully up to date.
	}

	indexerStatus := db.IndexerStatus{
		ID:     a.params.IndexerID,
		Height: uint64(height),
	}
	_, err := a.db.UpdateIndexerStatus(&indexerStatus)
	if err != nil {
		return errors.Wrapf(err, "error updating indexer status for %d height %d", indexerStatus.ID, indexerStatus.Height)
	}
	a.Height = uint64(height)
	return nil
}
