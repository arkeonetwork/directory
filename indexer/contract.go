package indexer

import (
	"fmt"
	"strconv"

	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/ArkeoNetwork/directory/pkg/utils"
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
			if evt.ContractType, err = utils.ParseContractType(v); err != nil {
				return evt, fmt.Errorf("unexpected contract type %s", v)
			}
		}
	}
	if evt.DelegatePubkey == "" {
		evt.DelegatePubkey = evt.ClientPubkey
	}
	return evt, nil
}

/*
	switch k {
	case "pubkey":
		evt.Contract.ProviderPubKey, err = common.NewPubKey(v)
		if err != nil {
			return evt, err
		}
	case "chain":
		evt.Contract.Chain, err = common.NewChain(v)
		if err != nil {
			return evt, err
		}
	case "client":
		evt.Contract.Client, err = common.NewPubKey(v)
		if err != nil {
			return evt, err
		}
	case "delegate":
		evt.Contract.Delegate, err = common.NewPubKey(v)
		if err != nil {
			return evt, err
		}
	case "type":
		evt.Contract.Type = types.ContractType(types.ContractType_value[v])
		if err != nil {
			return evt, err
		}
	case "height":
		evt.Contract.Height, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return evt, err
		}
	case "nonce":
		evt.Contract.Nonce, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return evt, err
		}
	case "paid":
		evt.Paid, ok = cosmos.NewIntFromString(v)
		if !ok {
			return evt, fmt.Errorf("cannot parse %s as int", v)
		}
	case "reserve":
		evt.Reserve, ok = cosmos.NewIntFromString(v)
		if !ok {
			return evt, fmt.Errorf("cannot parse %s as int", v)
		}
	}
*/
