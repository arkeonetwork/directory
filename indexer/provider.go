package indexer

import (
	"fmt"
	"math/big"
	"net/url"
	"strconv"

	"github.com/ArkeoNetwork/directory/pkg/db"
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
	// insert event
	log.Infof("update finished for contract %d", ent.ID)
	return nil
}

func (a *IndexerApp) handleModProviderEvent(evt types.ModProviderEvent) error {
	provider, err := a.db.FindProvider(evt.Pubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.Pubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("cannot mod provider, DNE %s %s", evt.Pubkey, evt.Chain)
	}
	provider.MetadataURI = evt.MetadataURI
	provider.MetadataNonce = evt.MetadataNonce
	provider.Status = evt.Status
	provider.MinContractDuration = evt.MinContractDuration
	provider.MaxContractDuration = evt.MaxContractDuration
	provider.SubscriptionRate = evt.SubscriptionRate
	provider.PayAsYouGoRate = evt.PayAsYouGoRate

	if _, err = a.db.UpdateProvider(provider); err != nil {
		return errors.Wrapf(err, "error updating provider for mod event %s chain %s", provider.Pubkey, provider.Chain)
	}
	log.Infof("updated provider %s chain %s", provider.Pubkey, provider.Chain)
	if _, err = a.db.InsertModProviderEvent(provider.ID, evt); err != nil {
		return errors.Wrapf(err, "error inserting ModProviderEvent for %s chain %s", evt.Pubkey, evt.Chain)
	}
	return nil
}

func (a *IndexerApp) handleBondProviderEvent(evt types.BondProviderEvent) error {
	provider, err := a.db.FindProvider(evt.Pubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.Pubkey, evt.Chain)
	}
	if provider == nil {
		// new provider for chain, insert
		if provider, err = a.createProvider(evt); err != nil {
			return errors.Wrapf(err, "error creating provider %s chain %s", evt.Pubkey, evt.Chain)
		}
	} else {
		if evt.BondAbsolute != nil {
			provider.Bond = evt.BondAbsolute.String()
		}
		if _, err = a.db.UpdateProvider(provider); err != nil {
			return errors.Wrapf(err, "error updating provider for bond event %s chain %s", evt.Pubkey, evt.Chain)
		}
	}

	log.Debugf("handled bond provider event for %s chain %s", evt.Pubkey, evt.Chain)
	if _, err = a.db.InsertBondProviderEvent(provider.ID, evt); err != nil {
		return errors.Wrapf(err, "error inserting BondProviderEvent for %s chain %s", evt.Pubkey, evt.Chain)
	}
	return nil
}

func (a *IndexerApp) createProvider(evt types.BondProviderEvent) (*db.ArkeoProvider, error) {
	// new provider for chain, insert
	provider := &db.ArkeoProvider{Pubkey: evt.Pubkey, Chain: evt.Chain, Bond: evt.BondAbsolute.String()}
	entity, err := a.db.InsertProvider(provider)
	if err != nil {
		return nil, errors.Wrapf(err, "error inserting provider %s %s", evt.Pubkey, evt.Chain)
	}
	if entity == nil {
		return nil, fmt.Errorf("nil entity after inserting provider")
	}
	log.Debugf("inserted provider record %d for %s %s", entity.ID, evt.Pubkey, evt.Chain)
	provider.Entity = *entity
	return provider, nil
}

func validateMetadataURI(uri string) bool {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return false
	}
	return true
}

func validateProviderStatus(s string) bool {
	switch types.ProviderStatus(s) {
	case types.ProviderStatusOffline:
		return true
	case types.ProviderStatusOnline:
		return true
	default:
		return false
	}
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
			if ok = validateChain(v); !ok {
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
			if evt.Height, err = strconv.ParseInt(v, 10, 64); err != nil {
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

	return evt, nil
}

func parseBondProviderEvent(input map[string]string) (types.BondProviderEvent, error) {
	// var err error
	var ok bool
	evt := types.BondProviderEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.Pubkey = v
		case "chain":
			if ok = validateChain(v); !ok {
				return evt, fmt.Errorf("invalid chain %s", v)
			}
			evt.Chain = v
		case "txID":
			evt.TxID = v
		case "bond_rel":
			evt.BondRelative, ok = new(big.Int).SetString(v, 10)
			if !ok {
				return evt, fmt.Errorf("cannot parse %s as int", v)
			}
		case "bond_abs":
			evt.BondAbsolute, ok = new(big.Int).SetString(v, 10)
			if !ok {
				return evt, fmt.Errorf("cannot parse %s as int", v)
			}
		}
	}

	return evt, nil
}

func parseModProviderEvent(input map[string]string) (types.ModProviderEvent, error) {
	var err error
	var ok bool
	evt := types.ModProviderEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.Pubkey = v
		case "chain":
			if ok = validateChain(v); !ok {
				return evt, fmt.Errorf("invalid chain %s", v)
			}
			evt.Chain = v
		case "txID":
			evt.TxID = v
		case "metadata_uri":
			if ok = validateMetadataURI(v); !ok {
				return evt, fmt.Errorf("invalid metadata_uri %s", v)
			}
			evt.MetadataURI = v
		case "metadata_nonce":
			if evt.MetadataNonce, err = strconv.ParseUint(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing metadata nonce %s", v)
			}
		case "status":
			if ok = validateProviderStatus(v); !ok {
				return evt, fmt.Errorf("invalid status %s", v)
			}
			evt.Status = types.ProviderStatus(v)
		case "min_contract_duration":
			if evt.MinContractDuration, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing min-contract-duration %s", v)
			}
		case "max_contract_duration":
			if evt.MaxContractDuration, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing max-contract-duration %s", v)
			}
		case "subscription_rate":
			if evt.SubscriptionRate, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing subscription_rate %s", v)
			}
		case "pay-as-you-go_rate":
			if evt.PayAsYouGoRate, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing pay-as-you-go_rate %s", v)
			}
		default:
			log.Warnf("not a support attribute for mod-provider %s", k)
		}
	}

	return evt, nil
}
