package indexer

import (
	"fmt"
	"strconv"

	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/ArkeoNetwork/directory/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func (a *IndexerApp) handleOpenContractEvent(evt types.OpenContractEvent) error {
	provider, err := a.db.FindProvider(evt.ProviderPubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.ProviderPubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("no provider found: DNE %s %s", evt.ProviderPubkey, evt.Chain)
	}
	ent, err := a.db.UpsertContract(provider.ID, evt)
	if err != nil {
		return errors.Wrapf(err, "error upserting contract")
	}
	if _, err = a.db.InsertOpenContractEvent(ent.ID, evt); err != nil {
		return errors.Wrapf(err, "error inserting open contract event")
	}

	log.Infof("update finished for contract %d", ent.ID)
	return nil
}

func parseContractSettlementEvent(input map[string]string) (types.ContractSettlementEvent, error) {
	ve := types.ContractSettlementEvent{}
	if err := mapstructure.Decode(input, &ve); err != nil {
		return ve, errors.Wrapf(err, "error reflecting properties to contract settlement event")
	}
	log.Infof("ve: %#v", ve)
	return ve, nil
}

func parseEvent(input map[string]string, target interface{}) error {
	if err := mapstructure.Decode(input, target); err != nil {
		return errors.Wrapf(err, "error reflecting properties to target")
	}

	return nil
}

func parseValidatorPayoutEvent(input map[string]string) (types.ValidatorPayoutEvent, error) {
	ve := types.ValidatorPayoutEvent{}
	if err := mapstructure.Decode(input, &ve); err != nil {
		return ve, errors.Wrapf(err, "error reflecting properties to validator payout event")
	}
	log.Infof("ve: %#v", ve)
	return ve, nil
}

func parseOpenContractEvent(input map[string]string) (types.OpenContractEvent, error) {
	var ok bool
	var err error
	evt := types.OpenContractEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.ProviderPubkey = v
		case "chain":
			if ok = utils.ValidateChain(v); !ok {
				return evt, fmt.Errorf("invalid chain %s", v)
			}
			evt.Chain = v
		case "delegate":
			evt.DelegatePubkey = v
		case "client":
			evt.ClientPubkey = v
		case "txID":
			evt.TxID = v
		case "duration":
			if evt.Duration, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing duration %s", v)
			}
		case "height":
			if evt.Height, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing height %s", v)
			}
		case "open_cost":
			if evt.OpenCost, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing open_cost %s", v)
			}
		case "rate":
			if evt.Rate, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing rate %s", v)
			}
		case "type":
			if types.ContractType(v) == types.ContractTypePayAsYouGo {
				evt.ContractType = types.ContractType(v)
			} else if types.ContractType(v) == types.ContractTypeSubscription {
				evt.ContractType = types.ContractType(v)
			} else {
				return evt, fmt.Errorf("unexpected contract type %s", v)
			}
		}
	}
	if evt.DelegatePubkey == "" {
		evt.DelegatePubkey = evt.ClientPubkey
	}
	return evt, nil
}
