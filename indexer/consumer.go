package indexer

import (
	"context"
	"encoding/hex"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/pkg/errors"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (a *IndexerApp) consumeEvents(client *tmclient.HTTP) error {
	blockEvents := subscribe(client, "tm.event = 'NewBlockHeader'")
	bondProviderEvents := make(chan int) // := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgBondProvider'")
	modProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgModProvider'")
	openContractEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")
	closeContractEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'")
	claimContractIncomeEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgClaimContractIncome'")
	// openContractEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")

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
			handleOpenContractEvent(a, &converted)
		case evt := <-bondProviderEvents:
			log.Debug(evt)
			// converted := convertWebSocketEvent("provider_bond", evt.Events)
			// handleBondProviderEvent(a, &converted)
		case evt := <-modProviderEvents:
			converted := convertWebSocketEvent("provider_mod", evt.Events)
			handleModProviderEvent(a, &converted)
		case evt := <-claimContractIncomeEvents:
			converted := convertWebSocketEvent("claim_contract_income", evt.Events)
			a.handleClaimContractIncomeEvent(converted)
		case evt := <-closeContractEvents:
			converted := convertWebSocketEvent("close_contract", evt.Events)
			log.Infof("close_contract: %#v", converted)
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
			log.Errorf("error reading block %d: %+v", nextHeight, err)
			retries = retries - 1
			log.Warnf("Getting next block results at height: %d failed, will retry %d more times", nextHeight, retries)
			time.Sleep(time.Second)
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
					convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					handleOpenContractEvent(a, &convertedEvent)
				case "provider_bond":
					convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					handleBondProviderEvent(a, &convertedEvent)
				case "provider_mod":
					convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					handleModProviderEvent(a, &convertedEvent)
				case "validator_payout":
					convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					a.handleValidatorPayoutEvent(&convertedEvent)
				case "contract_settlement":
					convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					a.handleContractSettlement(&convertedEvent)
				case "close_contract":
					convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					log.Warnf("close_contract event: %#v", convertedEvent)
				case "claim_contract_income":
					convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					log.Warnf("claim_contract_income event: %#v", convertedEvent)
				default:
					log.Warnf("received event %s", event.Type)
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

	if height, ok := raw["tx.height"]; ok && len(height) > 0 {
		newEvt["height"] = height[0]
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

func convertHistoricalEvent(event abcitypes.Event, height int64, txHash string) map[string]string {
	newEvt := make(map[string]string, 0)
	for _, attr := range event.Attributes {
		newEvt[string(attr.Key)] = string(attr.Value)
	}
	newEvt["height"] = strconv.FormatInt(height, 10)
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

func (a *IndexerApp) handleClaimContractIncomeEvent(event map[string]string) {
	contractSettlement := types.ContractSettlementEvent{}
	err := parseEvent(event, &contractSettlement)
	if err != nil {
		log.Errorf("error parsing contractSettlement: %+v", err)
		return
	}
	log.Infof("captured contractSettlement %#v", contractSettlement)
}

func (a *IndexerApp) handleContractSettlement(event *map[string]string) {

	contractSettlement, err := parseContractSettlementEvent(*event)
	if err != nil {
		log.Errorf("error parsing contractSettlement: %+v", err)
		return
	}
	log.Infof("captured contractSettlement %#v", contractSettlement)
}

func (a *IndexerApp) handleValidatorPayoutEvent(event *map[string]string) {

	payoutEvent, err := parseValidatorPayoutEvent(*event)
	if err != nil {
		log.Errorf("error parsing validatorPayoutEvent: %+v", err)
		return
	}
	log.Infof("captured payoutEvent %#v", payoutEvent)
	// openContractEvent, err := parseOpenContractEvent(*convertedEvent)

	// if err = a.handleOpenContractEvent(openContractEvent); err != nil {
	// 	log.Errorf("error handling open contract event: %+v", err)
	// 	return
	// }
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
		log.Errorf("error storing provider mod event: %+v", err)
		return
	}
	log.Infof("providerModEvent: %#v", modProviderEvent)
}
