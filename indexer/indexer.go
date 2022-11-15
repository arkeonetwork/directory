package indexer

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ArkeoNetwork/directory/internal/logging"
	"github.com/ArkeoNetwork/directory/pkg/db"
	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/pkg/errors"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
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

func (a *IndexerApp) consumeEvents(client *tmclient.HTTP) error {
	blockEvents := subscribe(client, "tm.event = 'NewBlockHeader'")
	bondProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgBondProvider'")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case evt := <-blockEvents:
			data, ok := evt.Data.(tmtypes.EventDataNewBlockHeader)
			if !ok {
				log.Errorf("event not block header: %T", evt.Data)
				continue
			}
			log.Debugf("received block: %d", data.Header.Height)
		case evt := <-bondProviderEvents:
			converted := convertEvent("provider_bond", evt.Events)
			providerBondEvent, err := parseProviderBondEvent(converted)
			if err != nil {
				log.Errorf("error parsing providerBondEvent: %+v", err)
				continue
			}
			if err = a.storeProviderBondEvent(providerBondEvent); err != nil {
				log.Errorf("error storing provider bond event: %+v", err)
				continue
			}
		case <-quit:
			log.Infof("received os quit signal")
			return nil
		}
	}
}

func (a *IndexerApp) storeProviderBondEvent(evt types.ProviderBondEvent) error {
	provider, err := a.db.FindProvider(evt.Pubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.Pubkey, evt.Chain)
	}
	if provider == nil {
		// new provider for chain, insert
		provider := &db.ArkeoProvider{Pubkey: evt.Pubkey, Chain: evt.Chain, Bond: evt.BondAbsolute.String()}
		entity, err := a.db.InsertProvider(provider)
		if err != nil {
			return errors.Wrapf(err, "error inserting provider %s %s", evt.Pubkey, evt.Chain)
		}
		log.Debugf("inserted provider record %d for %s %s", entity.ID, evt.Pubkey, evt.Chain)
	}

	// now store bond event for provider
	return nil
}

var validChains = map[string]struct{}{"arkeo-mainnet": {}, "eth-mainnet": {}, "btc-mainnet": {}}

func validateChain(chain string) (ok bool) {
	_, ok = validChains[chain]
	return
}

func parseProviderBondEvent(input map[string]string) (types.ProviderBondEvent, error) {
	// var err error
	var ok bool
	evt := types.ProviderBondEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.Pubkey = v
		case "chain":
			if ok = validateChain(v); !ok {
				return evt, fmt.Errorf("invalid chain %s", v)
			}
			evt.Chain = v
		case "bond_rel":
			evt.BondRelative, ok = new(big.Int).SetString(v, 10)
			if !ok {
				return evt, fmt.Errorf("cannot parse %s as int", v)
			}
		case "bond_abs":
			evt.BondAbsolute, ok = new(big.Int).SetString(v, 10)
			if !ok {
				return evt, fmt.Errorf("cannot parse %s as int", v)
			}
		}
	}

	return evt, nil
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
