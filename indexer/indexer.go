package indexer

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"os/signal"
	"strconv"
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
	modProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgModProvider'")
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
			bondProviderEvent, err := parseBondProviderEvent(converted)
			if err != nil {
				log.Errorf("error parsing bondProviderEvent: %+v", err)
				continue
			}
			if err = a.handleBondProviderEvent(bondProviderEvent); err != nil {
				log.Errorf("error handling provider bond event: %+v", err)
				continue
			}
		case evt := <-modProviderEvents:
			converted := convertEvent("provider_mod", evt.Events)
			modProviderEvent, err := parseModProviderEvent(converted)
			if err != nil {
				log.Errorf("error parsing modProviderEvent: %+v", err)
				continue
			}
			if err = a.handleModProviderEvent(modProviderEvent); err != nil {
				log.Errorf("error storing provider bond event: %+v", err)
				continue
			}
			log.Infof("providerModEvent: %#v", modProviderEvent)
		case <-quit:
			log.Infof("received os quit signal")
			return nil
		}
	}
}

func (a *IndexerApp) handleModProviderEvent(evt types.ModProviderEvent) error {
	provider, err := a.db.FindProvider(evt.Pubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.Pubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("cannot mod provider, DNE %s %s", evt.Pubkey, evt.Chain)
	}
	provider.MetadataURI = evt.MetadataURI
	provider.MetadataNonce = evt.MetadataNonce
	provider.Status = evt.Status
	provider.MinContractDuration = evt.MinContractDuration
	provider.MaxContractDuration = evt.MaxContractDuration
	provider.SubscriptionRate = evt.SubscriptionRate
	provider.PayAsYouGoRate = evt.PayAsYouGoRate

	if _, err = a.db.UpdateProvider(provider); err != nil {
		return errors.Wrapf(err, "error updating provider for mod event %s chain %s", provider.Pubkey, provider.Chain)
	}
	log.Infof("updated provider %s chain %s", provider.Pubkey, provider.Chain)
	if _, err = a.db.InsertModProviderEvent(provider.ID, evt); err != nil {
		return errors.Wrapf(err, "error inserting ModProviderEvent for %s chain %s", evt.Pubkey, evt.Chain)
	}
	return nil
}

func (a *IndexerApp) handleBondProviderEvent(evt types.BondProviderEvent) error {
	provider, err := a.db.FindProvider(evt.Pubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.Pubkey, evt.Chain)
	}
	if provider == nil {
		// new provider for chain, insert
		if provider, err = a.createProvider(evt); err != nil {
			return errors.Wrapf(err, "error creating provider %s chain %s", evt.Pubkey, evt.Chain)
		}
	} else {
		if evt.BondAbsolute != nil {
			provider.Bond = evt.BondAbsolute.String()
		}
		if _, err = a.db.UpdateProvider(provider); err != nil {
			return errors.Wrapf(err, "error updating provider for bond event %s chain %s", evt.Pubkey, evt.Chain)
		}
	}

	log.Debugf("handled bond provider event for %s chain %s", evt.Pubkey, evt.Chain)
	if _, err = a.db.InsertBondProviderEvent(provider.ID, evt); err != nil {
		return errors.Wrapf(err, "error inserting BondProviderEvent for %s chain %s", evt.Pubkey, evt.Chain)
	}
	return nil
}

func (a *IndexerApp) createProvider(evt types.BondProviderEvent) (*db.ArkeoProvider, error) {
	// new provider for chain, insert
	provider := &db.ArkeoProvider{Pubkey: evt.Pubkey, Chain: evt.Chain, Bond: evt.BondAbsolute.String()}
	entity, err := a.db.InsertProvider(provider)
	if err != nil {
		return nil, errors.Wrapf(err, "error inserting provider %s %s", evt.Pubkey, evt.Chain)
	}
	log.Debugf("inserted provider record %d for %s %s", entity.ID, evt.Pubkey, evt.Chain)
	return provider, nil
}

var validChains = map[string]struct{}{"arkeo-mainnet": {}, "eth-mainnet": {}, "btc-mainnet": {}}

func validateChain(chain string) (ok bool) {
	_, ok = validChains[chain]
	return
}

func validateMetadataURI(uri string) bool {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return false
	}
	return true
}

func validateProviderStatus(s string) bool {
	switch types.ProviderStatus(s) {
	case types.ProviderStatusOffline:
		return true
	case types.ProviderStatusOnline:
		return true
	default:
		return false
	}
}

func parseBondProviderEvent(input map[string]string) (types.BondProviderEvent, error) {
	// var err error
	var ok bool
	evt := types.BondProviderEvent{}

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

func parseModProviderEvent(input map[string]string) (types.ModProviderEvent, error) {
	var err error
	var ok bool
	evt := types.ModProviderEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.Pubkey = v
		case "chain":
			if ok = validateChain(v); !ok {
				return evt, fmt.Errorf("invalid chain %s", v)
			}
			evt.Chain = v
		case "metadata_uri":
			if ok = validateMetadataURI(v); !ok {
				return evt, fmt.Errorf("invalid metadata_uri %s", v)
			}
			evt.MetadataURI = v
		case "metadata_nonce":
			if evt.MetadataNonce, err = strconv.ParseUint(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing metadata nonce %s", v)
			}
		case "status":
			if ok = validateProviderStatus(v); !ok {
				return evt, fmt.Errorf("invalid status %s", v)
			}
			evt.Status = types.ProviderStatus(v)
		case "min_contract_duration":
			if evt.MinContractDuration, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing min-contract-duration %s", v)
			}
		case "max_contract_duration":
			if evt.MaxContractDuration, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing max-contract-duration %s", v)
			}
		case "subscription_rate":
			if evt.SubscriptionRate, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing subscription_rate %s", v)
			}
		case "pay-as-you-go_rate":
			if evt.PayAsYouGoRate, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing pay-as-you-go_rate %s", v)
			}
		default:
			log.Warnf("not a support attribute for mod-provider %s", k)
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
