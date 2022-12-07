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

func (a *IndexerApp) handleCloseContractEvent(evt types.CloseContractEvent) error {
	provider, err := a.db.FindProvider(evt.ProviderPubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.ProviderPubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("no provider found: DNE %s %s", evt.ProviderPubkey, evt.Chain)
	}
	// TODO: imeplement data model and insert event
	// ent, err := a.db.UpsertContract(provider.ID, evt)
	// if err != nil {
	// 	return errors.Wrapf(err, "error upserting contract")
	// }
	// if _, err = a.db.InsertOpenContractEvent(ent.ID, evt); err != nil {
	// 	return errors.Wrapf(err, "error inserting open contract event")
	// }

	//log.Infof("update finished for contract settlement event %d", ent.ID)
	return nil
}

func (a *IndexerApp) handleContractSettlementEvent(evt types.ContractSettlementEvent) error {
	provider, err := a.db.FindProvider(evt.ProviderPubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.ProviderPubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("no provider found: DNE %s %s", evt.ProviderPubkey, evt.Chain)
	}
	// TODO: imeplement data model and insert event
	// ent, err := a.db.UpsertContract(provider.ID, evt)
	// if err != nil {
	// 	return errors.Wrapf(err, "error upserting contract")
	// }
	// if _, err = a.db.InsertOpenContractEvent(ent.ID, evt); err != nil {
	// 	return errors.Wrapf(err, "error inserting open contract event")
	// }

	//log.Infof("update finished for contract settlement event %d", ent.ID)
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
			evt.ContractType, err = utils.ParseContractType(v)
			return evt, fmt.Errorf("unexpected contract type error %+v", err)
		}
	}
	if evt.DelegatePubkey == "" {
		evt.DelegatePubkey = evt.ClientPubkey
	}
	return evt, nil
}

func parseCloseContractEvent(input map[string]string) (types.CloseContractEvent, error) {
	var ok bool
	evt := types.CloseContractEvent{}

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
		}
	}
	if evt.DelegatePubkey == "" {
		evt.DelegatePubkey = evt.ClientPubkey
	}
	return evt, nil
}

func parseContractSettlementEvent(input map[string]string) (types.ContractSettlementEvent, error) {
	var err error
	var ok bool
	evt := types.ContractSettlementEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.ProviderPubkey = v
		case "chain":
			if ok = utils.ValidateChain(v); !ok {
				return evt, fmt.Errorf("invalid chain %s", v)
			}
			evt.Chain = v
		case "client":
			evt.ClientPubkey = v
		case "delegate":
			evt.DelegatePubkey = v
		case "txID":
			evt.TxID = v
		case "type":
			evt.ContractType, err = utils.ParseContractType(v)
			return evt, fmt.Errorf("unexpected contract type error %+v", err)
		case "nonce":
			if evt.ContractNonce, err = strconv.ParseUint(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing nonce %s", v)
			}
		case "height":
			if evt.Height, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing height %s", v)
			}
		case "paid":
			if evt.Paid, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing paid %s", v)
			}
		case "reserve":
			if evt.Reserve, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing reserve %s", v)
			}
		default:
			log.Warnf("not a support attribute for mod-provider %s", k)
		}
	}

	return evt, nil
}
