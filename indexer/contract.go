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
	if _, err = a.db.UpsertOpenContractEvent(ent.ID, evt); err != nil {
		return errors.Wrapf(err, "error inserting open contract event")
	}

	log.Infof("update finished for contract %d", ent.ID)
	return nil
}

func (a *IndexerApp) handleCloseContractEvent(evt types.CloseContractEvent) error {
	contract, err := a.db.FindContractByPubKeys(evt.Chain, evt.ProviderPubkey, evt.GetDelegatePubkey())
	if err != nil {
		return errors.Wrapf(err, "error finding contract for %s:%s %s", evt.ProviderPubkey, evt.Chain, evt.GetDelegatePubkey())
	}
	if contract == nil {
		return fmt.Errorf("no contract found: %s:%s %s", evt.ProviderPubkey, evt.Chain, evt.GetDelegatePubkey())
	}
	if _, err = a.db.UpsertCloseContractEvent(contract.ID, evt); err != nil {
		return errors.Wrapf(err, "error upserting open contract event")
	}

	log.Infof("update finished for close contract %d", contract.ID)
	return nil
}

func (a *IndexerApp) handleContractSettlementEvent(evt types.ContractSettlementEvent) error {
	log.Infof("receieved contractSettlementEvent %#v", evt)
	provider, err := a.db.FindProvider(evt.ProviderPubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.ProviderPubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("cannot claim income provider %s on chain %s DNE", evt.ProviderPubkey, evt.Chain)
	}
	contract, err := a.db.FindContract(provider.ID, evt.ClientPubkey)
	if err != nil {
		return errors.Wrapf(err, "error finding contract provider %s chain %s", evt.ProviderPubkey, evt.Chain)
	}
	if contract == nil {
		return fmt.Errorf("no contract found for %s:%s delegPub: %s", evt.ProviderPubkey, evt.Chain, evt.ClientPubkey)
	}
	if _, err = a.db.UpsertContractSettlementEvent(contract.ID, evt); err != nil {
		return errors.Wrapf(err, "error upserting contract settlement event")
	}

	return nil
}
