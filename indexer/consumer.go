package indexer

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (a *IndexerApp) consumeEvents(client *tmclient.HTTP) error {
	blockEvents := subscribe(client, "tm.event = 'NewBlockHeader'")
	bondProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgBondProvider'")
	modProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgModProvider'")
	openContractEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")
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
			a.handleBlockEvent(data.Header.Height)
		case evt := <-openContractEvents:
			log.Debugf("received open contract event")
			converted := convertWebSocketEvent("open_contract", evt.Events)
			log.Infof("converted open_contract map: %#v", converted)
			handleOpenContractEvent(a, &converted)
		case evt := <-bondProviderEvents:
			converted := convertWebSocketEvent("provider_bond", evt.Events)
			handleBondProviderEvent(a, &converted)
		case evt := <-modProviderEvents:
			converted := convertWebSocketEvent("provider_mod", evt.Events)
			handleModProviderEvent(a, &converted)
		case <-quit:
			log.Infof("received os quit signal")
			return nil
		}
	}
}

func (a *IndexerApp) consumeHistoricalEvents(client *tmclient.HTTP) error {
	var currentBlock *ctypes.ResultBlock
	currentBlock, err := client.Block(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "error getting current block")
	}
	blocksToSync := currentBlock.Block.Height - int64(a.Height)
	log.Infof("Current block %d, syncing from block %d. %d blocks to go", currentBlock.Block.Height, a.Height, blocksToSync)
	retries := 10
	blocksSynced := 0
	for currentBlock.Block.Height > int64(a.Height) {
		nextHeight := int64(a.Height)
		nextBlock, err := client.Block(context.Background(), &nextHeight)
		if err != nil {
			retries = retries - 1
			log.Warnf("Getting next block results at height: %d failed, will retry %d more times", nextHeight, retries)
			if retries == 0 {
				log.Errorf("Getting next block results at height: %d failed with no additional retries", nextHeight)
				return errors.Wrapf(err, "error getting block results at height: %d", nextHeight)
			}
			continue
		}

		for _, transaction := range nextBlock.Block.Txs {
			txInfo, err := client.Tx(context.Background(), transaction.Hash(), false)
			if err != nil {
				log.Warnf("failed to get transaction data for %s", transaction.Hash())
				continue
			}

			for _, event := range txInfo.TxResult.Events {
				switch event.Type {
				case "open_contract":
					convertedEvent := convertHistoricalEvent(event, string(transaction.Hash()))
					handleOpenContractEvent(a, &convertedEvent)
				case "provider_bond":
					convertedEvent := convertHistoricalEvent(event, string(transaction.Hash()))
					handleBondProviderEvent(a, &convertedEvent)
				case "provider_mod":
					convertedEvent := convertHistoricalEvent(event, string(transaction.Hash()))
					handleModProviderEvent(a, &convertedEvent)
				}
			}
		}
		blocksSynced++
		if blocksSynced%500 == 0 {
			log.Debugf("synced %d of initial %d", blocksSynced, blocksToSync)
		}

		a.Height++
		if currentBlock.Block.Height == int64(a.Height) {
			// we should update to see if new blocks have become available while we were processing
			currentBlock, err = client.Block(context.Background(), nil)
			if err != nil {
				return errors.Wrap(err, "error getting current block")
			}
			blocksToSync = currentBlock.Block.Height - int64(a.Height)
			blocksSynced = 0
		}
	}
	return nil
}

// TODO: if there are multiple of the same type of event, this may be
// problematic, multiple events may get purged into one (not sure)
func convertWebSocketEvent(etype string, raw map[string][]string) map[string]string {
	newEvt := make(map[string]string, 0)
	if txID, ok := raw["tx.hash"]; ok {
		newEvt["txID"] = txID[0]
	} else {
		log.Warnf("no tx.hash in event attributes: %#v", raw)
	}

	for k, v := range raw {
		if strings.HasPrefix(k, etype+".") {
			parts := strings.SplitN(k, ".", 2)
			newEvt[parts[1]] = v[0]
		}
	}

	return newEvt
}

func convertHistoricalEvent(event abcitypes.Event, txHash string) map[string]string {
	newEvt := make(map[string]string, 0)
	for _, attr := range event.Attributes {
		newEvt[string(attr.Key)] = string(attr.Value)
	}
	newEvt["txID"] = txHash
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

func handleBondProviderEvent(a *IndexerApp, convertedEvent *map[string]string) {
	bondProviderEvent, err := parseBondProviderEvent(*convertedEvent)
	if err != nil {
		log.Errorf("error parsing bondProviderEvent: %+v", err)
		return
	}
	if err = a.handleBondProviderEvent(bondProviderEvent); err != nil {
		log.Errorf("error handling provider bond event: %+v", err)
		return
	}
}

func handleOpenContractEvent(a *IndexerApp, convertedEvent *map[string]string) {
	openContractEvent, err := parseOpenContractEvent(*convertedEvent)
	if err != nil {
		log.Errorf("error parsing openContractEvent: %+v", err)
		return
	}
	if err = a.handleOpenContractEvent(openContractEvent); err != nil {
		log.Errorf("error handling open contract event: %+v", err)
		return
	}
}

func handleModProviderEvent(a *IndexerApp, convertedEvent *map[string]string) {
	modProviderEvent, err := parseModProviderEvent(*convertedEvent)
	if err != nil {
		log.Errorf("error parsing modProviderEvent: %+v", err)
		return
	}
	if err = a.handleModProviderEvent(modProviderEvent); err != nil {
		log.Errorf("error storing provider bond event: %+v", err)
		return
	}
	log.Infof("providerModEvent: %#v", modProviderEvent)
}
