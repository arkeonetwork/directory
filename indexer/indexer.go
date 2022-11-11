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
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var log = logging.WithoutFields()

type JsonRpcMsg struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

type SubscribeMsg struct {
	JsonRpcMsg
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type IndexerAppParams struct{}
type IndexerApp struct {
	params IndexerAppParams
	done   chan struct{}
}

const (
	// todo configs
	chainID             = "arkeo"
	bech32PrefixAccAddr = "rko"
	bech32PrefixAccPub  = "rkopub"
	tendermintWsUrl     = "tcp://localhost:26657"
)

func NewIndexer(params IndexerAppParams) *IndexerApp {
	configure()
	return &IndexerApp{params: params}
}

func (a *IndexerApp) Run() (done <-chan struct{}, err error) {
	// initialize by reading all existing providers?
	a.done = make(chan struct{})
	go a.start()
	return a.done, nil
}

func (a *IndexerApp) start() {
	log.Infof("starting realtime indexer")
	client, err := tmclient.New(tendermintWsUrl, "/websocket")
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
			log.Debugf("received bondProviderEvent: %#v", evt)
			converted := convertEvent("provider_bond", evt.Events)
			providerBondEvent, err := parseProviderBondEvent(converted)
			if err != nil {
				log.Errorf("error parsing providerBondEvent: %+v", err)
				continue
			}
			storeProviderBondEvent(providerBondEvent)
		case <-quit:
			log.Infof("received os quit signal")
			return nil
		}
	}
}

func storeProviderBondEvent(evt ProviderBondEvent) error {
	log.Infof("storing ProviderBondEvent %#v", evt)

	return nil
}

var validChains = map[string]struct{}{"arkeo-mainnet": {}, "eth-mainnet": {}, "btc-mainnet": {}}

func validateChain(chain string) (ok bool) {
	_, ok = validChains[chain]
	return
}

type ProviderBondEvent struct {
	PubKey       string
	Chain        string
	BondRelative *big.Int
	BondAbsolute *big.Int
}

func parseProviderBondEvent(input map[string]string) (ProviderBondEvent, error) {
	// var err error
	var ok bool
	evt := ProviderBondEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.PubKey = v
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
