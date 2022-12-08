package indexer

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/ArkeoNetwork/directory/pkg/db"
	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/ArkeoNetwork/directory/pkg/utils"
	"github.com/pkg/errors"
)

func (a *IndexerApp) handleModProviderEvent(evt types.ModProviderEvent) error {
	provider, err := a.db.FindProvider(evt.Pubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.Pubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("cannot mod provider, DNE %s %s", evt.Pubkey, evt.Chain)
	}

	isMetaDataUpdated := provider.MetadataNonce != evt.MetadataNonce

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

	if isMetaDataUpdated {
		log.Debugf("updating provider metadata for provider %s", provider.Pubkey)
		if !validateMetadataURI(provider.MetadataURI) {
			log.Warnf("updating provider metadata for provider %s failed due to bad MetadataURI %s", provider.MetadataURI)
			return nil
		}
		providerMetadata, err := utils.DownloadProviderMetadata(provider.MetadataURI, 5, 1e6)
		if err != nil {
			log.Warnf("updating provider metadata for provider %s failed %v", err)
			return nil
		}
		if _, err = a.db.UpsertProviderMetadata(provider.ID, *providerMetadata); err != nil {
			return errors.Wrapf(err, "error updating provider metadta for mod event %s chain %s", provider.Pubkey, provider.Chain)
		}
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
		if evt.BondAbsolute != "" {
			provider.Bond = evt.BondAbsolute
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
	provider := &db.ArkeoProvider{Pubkey: evt.Pubkey, Chain: evt.Chain, Bond: evt.BondAbsolute}
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

// func parseBondProviderEvent(input map[string]string) (types.BondProviderEvent, error) {
// 	var err error
// 	var ok bool
// 	evt := types.BondProviderEvent{}

// 	for k, v := range input {
// 		switch k {
// 		case "pubkey":
// 			evt.Pubkey = v
// 		case "chain":
// 			if ok = utils.ValidateChain(v); !ok {
// 				return evt, fmt.Errorf("invalid chain %s", v)
// 			}
// 			evt.Chain = v
// 		case "txID":
// 			evt.TxID = v
// 		case "height":
// 			if evt.Height, err = strconv.ParseInt(v, 10, 64); err != nil {
// 				return evt, errors.Wrapf(err, "error parsing height %s", v)
// 			}
// 		case "bond_rel":
// 			evt.BondRelative, ok = new(big.Int).SetString(v, 10)
// 			if !ok {
// 				return evt, fmt.Errorf("cannot parse %s as int", v)
// 			}
// 		case "bond_abs":
// 			evt.BondAbsolute, ok = new(big.Int).SetString(v, 10)
// 			if !ok {
// 				return evt, fmt.Errorf("cannot parse %s as int", v)
// 			}
// 		}
// 	}

// 	return evt, nil
// }

func parseModProviderEvent(input map[string]string) (types.ModProviderEvent, error) {
	var err error
	var ok bool
	evt := types.ModProviderEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.Pubkey = v
		case "chain":
			if ok = utils.ValidateChain(v); !ok {
				return evt, fmt.Errorf("invalid chain %s", v)
			}
			evt.Chain = v
		case "height":
			if evt.Height, err = strconv.ParseInt(v, 10, 64); err != nil {
				return evt, errors.Wrapf(err, "error parsing height %s", v)
			}
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
			log.Warnf("not a supported attribute for mod-provider %s", k)
		}
	}

	return evt, nil
}
