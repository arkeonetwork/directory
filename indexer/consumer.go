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

// type attributeProvider interface {
// 	attributes() map[string]string
// }

type attributes func() map[string]string

func wsAttributeSource(src ctypes.ResultEvent) func() map[string]string {
	results := make(map[string]string, len(src.Events))
	for k, v := range src.Events {
		if len(v) > 0 {
			results[k] = v[0]
		}
		if len(v) > 1 {
			log.Warnf("attrib %s has %d values in array", k, len(v))
		}
	}
	return func() map[string]string { return results }
}

func tmAttributeSource(src abcitypes.Event, height int64, txHash string) func() map[string]string {
	newEvt := make(map[string]string, 0)
	for _, attr := range src.Attributes {
		newEvt[string(attr.Key)] = string(attr.Value)
	}
	newEvt["height"] = strconv.FormatInt(height, 10)
	newEvt["txID"] = txHash
	return func() map[string]string { return newEvt }
}

func (a *IndexerApp) consumeEvents(client *tmclient.HTTP) error {
	blockEvents := subscribe(client, "tm.event = 'NewBlockHeader'")
	bondProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgBondProvider'")
	modProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgModProvider'")
	openContractEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")
	closeContractEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'")
	claimContractIncomeEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgClaimContractIncome'")
	// openContractEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		txHash := "abc23"
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
			openContractEvent := types.OpenContractEvent{}
			if err := convertEvent(wsAttributeSource(evt), openContractEvent, "open_contract", 0, txHash); err != nil {
				log.Errorf("error converting open_contract event: %+v", err)
				break
			}
			if err := a.handleOpenContractEvent(openContractEvent); err != nil {
				log.Errorf("error handling open_contract event: %+v", err)
				break
			}
		case evt := <-bondProviderEvents:
			log.Debug(evt)
			bondProviderEvent := types.BondProviderEvent{}
			if err := convertEvent(wsAttributeSource(evt), bondProviderEvent, "bond_provider", 0, txHash); err != nil {
				log.Errorf("error converting bond_provider event: %+v", err)
				break
			}
			if err := a.handleBondProviderEvent(bondProviderEvent); err != nil {
				log.Errorf("error handling bond_provider event: %+v", err)
			}
		case evt := <-modProviderEvents:
			log.Debug(evt)
			// converted := convertEvent("provider_mod", evt.Events)
			// converted := convertEvent(wsAttributeSource(evt), "provider_mod", 0, hash)
			// handleModProviderEvent(a, &converted)

		case evt := <-claimContractIncomeEvents:
			log.Debug(evt)
		// converted := convertEvent("claim_contract_income", evt.Events)
		// converted := convertEvent(wsAttributeSource(evt), "claim_contract_income", 0, hash)
		// a.handleClaimContractIncomeEvent(converted)
		case evt := <-closeContractEvents:
			log.Debug(evt)
		// converted := convertEvent("close_contract", evt.Events)
		// converted := convertEvent(wsAttributeSource(evt), "close_contract", 0, hash)
		// log.Infof("close_contract: %#v", converted)
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

			txHash := "hist123"
			for _, event := range txInfo.TxResult.Events {
				switch event.Type {
				case "open_contract":
					openContractEvent := types.OpenContractEvent{}
					if err := convertEvent(tmAttributeSource(event, currentBlock.Block.Height, txHash), openContractEvent, event.Type, currentBlock.Block.Height, txHash); err != nil {
						log.Errorf("error converting %s event: %+v", event.Type, err)
						break
					}
					// convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					// a.handleOpenContractEvent(openContractEvent)
					if err = a.handleOpenContractEvent(openContractEvent); err != nil {
						log.Errorf("error handling %s event: %+v", event.Type, err)
					}
				case "provider_bond":
					bondProviderEvent := types.BondProviderEvent{}
					if err = convertEvent(tmAttributeSource(event, currentBlock.Block.Height, txHash), &bondProviderEvent, event.Type, currentBlock.Block.Height, txHash); err != nil {
						log.Errorf("error converting %s event: %+v", event.Type, err)
						break
					}
					if err = a.handleBondProviderEvent(bondProviderEvent); err != nil {
						log.Errorf("error handling %s event: %+v", event.Type, err)
					}
				case "provider_mod":
					convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					handleModProviderEvent(a, &convertedEvent)
				case "validator_payout":
					convertedEvent := convertHistoricalEvent(event, txInfo.Height, strings.ToUpper(hex.EncodeToString(transaction.Hash()[:])))
					a.handleValidatorPayoutEvent(convertedEvent)
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

func convertEvent(attributeFunc attributes, target interface{}, etype string, height int64, txHash string) error {
	m := attributeFunc()
	if err := parseEvent(m, &target); err != nil {
		log.Errorf("error parsingEvent: %+v", err)
	}
	return parseEvent(m, target)
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

// func handleBondProviderEvent(a *IndexerApp, convertedEvent *map[string]string) {
// 	bondProviderEvent, err := parseBondProviderEvent(*convertedEvent)
// 	if err != nil {
// 		log.Errorf("error parsing bondProviderEvent: %+v", err)
// 		return
// 	}
// 	if err = a.handleBondProviderEvent(bondProviderEvent); err != nil {
// 		log.Errorf("error handling provider bond event: %+v", err)
// 		return
// 	}
// }

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

	// contractSettlement, err := parseContractSettlementEvent(*event)
	contractSettlement := types.ContractSettlementEvent{}
	err := parseEvent(*event, &contractSettlement)
	if err != nil {
		log.Errorf("error parsing contractSettlement: %+v", err)
		return
	}
	log.Infof("captured contractSettlement %#v", contractSettlement)
}

func (a *IndexerApp) handleValidatorPayoutEvent(event map[string]string) {

	payoutEvent := types.ValidatorPayoutEvent{}
	err := parseEvent(event, payoutEvent)
	if err != nil {
		log.Errorf("error parsing validatorPayoutEvent: %+v", err)
		return
	}
	log.Infof("captured payoutEvent %#v", payoutEvent)
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
