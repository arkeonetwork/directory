package indexer

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
			openContractEvent, err := parseOpenContractEvent(converted)
			if err != nil {
				log.Errorf("error parsing openContractEvent: %+v", err)
				continue
			}
			if err = a.handleOpenContractEvent(openContractEvent); err != nil {
				log.Errorf("error handling open contract event: %+v", err)
				continue
			}
		case evt := <-bondProviderEvents:
			converted := convertWebSocketEvent("provider_bond", evt.Events)
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
			converted := convertWebSocketEvent("provider_mod", evt.Events)
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

func convertHistoricalEvent(event abcitypes.Event) map[string]string {
	newEvt := make(map[string]string, 0)
	for _, attr := range event.Attributes {
		newEvt[string(attr.Key)] = string(attr.Value)
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

func (a *IndexerApp) consumeHistoricalEvents(client *tmclient.HTTP) error {
	var currentBlock *ctypes.ResultBlock
	currentBlock, err := client.Block(context.Background(), nil)
	if err != nil {
		// how to handle this?
		return err
	}
	log.Infof("Current block %d, syncing from block %d", currentBlock.Block.Height, a.Height)

	for currentBlock.Block.Height > int64(a.Height) {
		currentHeight := int64(a.Height)
		nextBlock, err := client.BlockResults(context.Background(), &currentHeight)
		for _, result := range nextBlock.TxsResults {
			for _, event := range result.Events {
				switch event.Type {
				case "open_contract":
				case "provider_bond":
					convertedEvent := convertHistoricalEvent(event)
					// currently stuck here trying to figure out how to get the transaction hash
					parseBondProviderEvent(convertedEvent)
				case "provider_mod":
				}
			}
		}

		a.Height = a.Height + 1
		if currentBlock.Block.Height == int64(a.Height) {
			// we should update to see if new blocks are available.
			currentBlock, err = client.Block(context.Background(), nil)
			if err != nil {
				// how to handle this?
				return err
			}
		}
	}
	return nil
}
