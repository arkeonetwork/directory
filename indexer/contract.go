package indexer

import (
	"fmt"

	"github.com/ArkeoNetwork/directory/pkg/types"
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

func (a *IndexerApp) handleContractSettlementEvent(evt types.ContractSettlementEvent) error {
	log.Infof("receieved contractSettlementEvent %#v", evt)
	provider, err := a.db.FindProvider(evt.Pubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.Pubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("cannot claim income provider %s on chain %s DNE", evt.Pubkey, evt.Chain)
	}
	contract, err := a.db.FindContract(provider.ID, evt.ClientPubkey)
	if err != nil {
		return errors.Wrapf(err, "error finding contract provider %s chain %s", evt.Pubkey, evt.Chain)
	}

	if _, err = a.db.UpsertContractSettlementEvent(contract.ID, evt); err != nil {
		return errors.Wrapf(err, "error upserting contract settlement event")
	}

	return nil
}
