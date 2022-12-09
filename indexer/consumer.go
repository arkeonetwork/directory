package indexer

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ArkeoNetwork/directory/pkg/db"
	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/mitchellh/mapstructure"
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
	attribs := make(map[string]string, len(src.Events))
	for k, v := range src.Events {
		if len(v) > 0 {
			key := k
			if sl := strings.Split(k, "."); len(sl) > 1 {
				key = sl[1]
			}
			if _, ok := attribs[key]; ok {
				log.Debugf("key %s already in results with value %s, overwriting with %s", key, attribs[key], v[0])
			}
			attribs[key] = v[0]
		}
		if len(v) > 1 {
			log.Warnf("attrib %s has %d array values: %v", k, len(v), v)
		}
	}
	attribs["eventHeight"] = attribs["height"]
	return func() map[string]string { return attribs }
}

func tmAttributeSource(tx tmtypes.Tx, evt abcitypes.Event, height uint64) func() map[string]string {
	attribs := make(map[string]string, 0)
	for _, attr := range evt.Attributes {
		attribs[string(attr.Key)] = string(attr.Value)
		log.Infof("%s: %s", attr.Key, string(attr.Value))
	}

	if tx != nil {
		if _, ok := attribs["hash"]; !ok {
			attribs["hash"] = strings.ToUpper(hex.EncodeToString(tx.Hash()[:]))
		}
	}

	attribs["eventHeight"] = strconv.FormatUint(height, 10)
	if _, ok := attribs["height"]; !ok {
		attribs["height"] = attribs["eventHeight"]
	}

	return func() map[string]string { return attribs }
}

func (a *IndexerApp) consumeEvents(client *tmclient.HTTP) error {
	blockEvents := subscribe(client, "tm.event = 'NewBlock'")
	bondProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgBondProvider'")
	modProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgModProvider'")
	openContractEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")
	closeContractEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'")
	claimContractIncomeEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgClaimContractIncome'")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case evt := <-blockEvents:
			data, ok := evt.Data.(tmtypes.EventDataNewBlock)
			if !ok {
				log.Errorf("event not block: %T", evt.Data)
				continue
			}
			log := log.WithField("height", strconv.FormatInt(data.Block.Height, 10))
			log.Debugf("received block: %d", data.Block.Height)

			if err := a.handleBlockEvent(data.Block.Height); err != nil {
				log.Errorf("error handling block event %d: %+v", data.Block.Height, err)
			}

			endBlockEvents := data.ResultEndBlock.Events
			log.Infof("block %d has %d endBlock events", data.Block.Height, len(endBlockEvents))
			for _, evt := range endBlockEvents {
				switch evt.GetType() {
				case "validator_payout":
					log.Infof("validator_payout event")
				case "contract_settlement":
					log.Debugf("contract_settlement")
					contractSettlementEvent := types.ContractSettlementEvent{}
					if err := convertEvent(tmAttributeSource(nil, evt, uint64(data.Block.Height)), &contractSettlementEvent); err != nil {
						log.Errorf("error converting contract_settlement event: %+v", err)
						break
					}
					if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
						log.Errorf("error handling open_contract event: %+v", err)
					}
				}
				attribs := tmAttributeSource(nil, evt, uint64(data.Block.Height))()
				log.Infof("%s: %#v", evt.GetType(), attribs)
			}

		case evt := <-openContractEvents:
			log.Debugf("received open contract event")
			openContractEvent := types.OpenContractEvent{}
			if err := convertEvent(wsAttributeSource(evt), &openContractEvent); err != nil {
				log.Errorf("error converting open_contract event: %+v", err)
				break
			}
			if err := a.handleOpenContractEvent(openContractEvent); err != nil {
				log.Errorf("error handling open_contract event: %+v", err)
			}
		case evt := <-bondProviderEvents:
			log.Debugf("received bond provider event")
			bondProviderEvent := types.BondProviderEvent{}
			if err := convertEvent(wsAttributeSource(evt), &bondProviderEvent); err != nil {
				log.Errorf("error converting bond_provider event: %+v", err)
				break
			}
			if err := a.handleBondProviderEvent(bondProviderEvent); err != nil {
				log.Errorf("error handling bond_provider event: %+v", err)
			}
		case evt := <-modProviderEvents:
			log.Debugf("received mod provider event")
			modProviderEvent := types.ModProviderEvent{}
			if err := convertEvent(wsAttributeSource(evt), &modProviderEvent); err != nil {
				log.Errorf("error converting mod_provider event: %+v", err)
				break
			}
			if err := a.handleModProviderEvent(modProviderEvent); err != nil {
				log.Errorf("error handling mod_provider event: %+v", err)
			}
		case evt := <-claimContractIncomeEvents:
			log.Debugf("received claim contract income event")
			contractSettlementEvent := types.ContractSettlementEvent{}
			attribs := wsAttributeSource(evt)
			if err := convertEvent(attribs, &contractSettlementEvent); err != nil {
				log.Errorf("error converting open_contract event: %+v", err)
				break
			}
			if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
				log.Errorf("error handling claim contract income event: %+v", err)
			}
		case evt := <-closeContractEvents:
			log.Debugf("received close_contract event")
			closeContractEvent := types.CloseContractEvent{}
			if err := convertEvent(wsAttributeSource(evt), &closeContractEvent); err != nil {
				log.Errorf("error converting close_contract event: %+v", err)
				break
			}
			if err := a.handleCloseContractEvent(closeContractEvent); err != nil {
				log.Errorf("error handling close contract event: %+v", err)
			}
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
		log := log.WithField("height", strconv.FormatUint(a.Height, 10))
		for _, transaction := range nextBlock.Block.Txs {
			txInfo, err := client.Tx(context.Background(), transaction.Hash(), false)
			if err != nil {
				log.Warnf("failed to get transaction data for %s", transaction.Hash())
				continue
			}

			for _, event := range txInfo.TxResult.Events {
				log.Debugf("received %s event", event.Type)
				if err := a.handleAbciEvent(event, transaction); err != nil {
					log.Errorf("error handling abci event %#v\n%+v", event, err)
				}
			}
		}
		blocksSynced++
		if blocksSynced%500 == 0 {
			log.Debugf("synced %d of initial %d", blocksSynced, blocksToSync)

			// update DB periodically to avoid having to sync all over
			indexerStatus := db.IndexerStatus{
				ID:     a.params.IndexerID,
				Height: uint64(nextHeight),
			}
			_, err := a.db.UpsertIndexerStatus(&indexerStatus)
			if err != nil {
				log.Warnf("error writing block status to db %#v", err)
			}
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

func (a *IndexerApp) handleAbciEvent(event abcitypes.Event, transaction tmtypes.Tx) error {
	var err error
	switch event.Type {
	case "provider_bond":
		bondProviderEvent := types.BondProviderEvent{}
		if err = convertEvent(tmAttributeSource(transaction, event, a.Height), &bondProviderEvent); err != nil {
			log.Errorf("error converting %s event: %+v", event.Type, err)
			break
		}
		if err = a.handleBondProviderEvent(bondProviderEvent); err != nil {
			log.Errorf("error handling %s event: %+v", event.Type, err)
		}
	case "provider_mod":
		modProviderEvent := types.ModProviderEvent{}
		if err = convertEvent(tmAttributeSource(transaction, event, a.Height), &modProviderEvent); err != nil {
			log.Errorf("error converting %s event: %+v", event.Type, err)
			break
		}
		if err = a.handleModProviderEvent(modProviderEvent); err != nil {
			log.Errorf("error handling %s event: %+v", event.Type, err)
		}
	case "open_contract":
		openContractEvent := types.OpenContractEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, a.Height), &openContractEvent); err != nil {
			log.Errorf("error converting %s event: %+v", event.Type, err)
			break
		}
		if err = a.handleOpenContractEvent(openContractEvent); err != nil {
			log.Errorf("error handling %s event: %+v", event.Type, err)
		}
	case "claim_contract_income":
		contractSettlementEvent := types.ContractSettlementEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, a.Height), &contractSettlementEvent); err != nil {
			log.Errorf("error converting claim_contract_income event: %+v", err)
			break
		}
		if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
			log.Errorf("error handling claim contract income event: %+v", err)
		}
	case "validator_payout":
		validatorPayoutEvent := types.ValidatorPayoutEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, a.Height), &validatorPayoutEvent); err != nil {
			log.Errorf("error converting validatorPayoutEvent event: %+v", err)
			break
		}
		if err := a.handleValidatorPayoutEvent(validatorPayoutEvent); err != nil {
			log.Errorf("error handling claim contract income event: %+v", err)
		}
	case "contract_settlement":
		contractSettlementEvent := types.ContractSettlementEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, a.Height), &contractSettlementEvent); err != nil {
			log.Errorf("error converting contractSettlementEvent: %+v", err)
			break
		}
		if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
			log.Errorf("error handling contractSettlementEvent: %+v", err)
		}
	case "close_contract":
		log.Debugf("received close_contract event")
		closeContractEvent := types.CloseContractEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, a.Height), &closeContractEvent); err != nil {
			log.Errorf("error converting close_contract event: %+v", err)
			break
		}
		if err := a.handleCloseContractEvent(closeContractEvent); err != nil {
			log.Errorf("error handling close contract event: %+v", err)
		}
	default:
		log.Warnf("received event %s", event.Type)
	}
	return nil
}

// copy attributes of map given by attributeFunc() to target which must be a pointer (map/slice implicitly ptr)
func convertEvent(attributeFunc attributes, target interface{}) error {
	return mapstructure.WeakDecode(attributeFunc(), target)
}

func subscribe(client *tmclient.HTTP, query string) <-chan ctypes.ResultEvent {
	out, err := client.Subscribe(context.Background(), "", query)
	if err != nil {
		log.Errorf("failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}
	return out
}

func (a *IndexerApp) handleValidatorPayoutEvent(event types.ValidatorPayoutEvent) error {
	return fmt.Errorf("not impl")
}
