package db

import (
	"context"
	"fmt"

	"github.com/ArkeoNetwork/directory/pkg/sentinel"
	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
)

type ArkeoProvider struct {
	Entity
	Pubkey string `db:"pubkey"`
	Chain  string `db:"chain"`
	// this is a DECIMAL type in the db
	Bond                string               `db:"bond"`
	MetadataURI         string               `db:"metadata_uri"`
	MetadataNonce       uint64               `db:"metadata_nonce"`
	Status              types.ProviderStatus `db:"status,text"`
	MinContractDuration int64                `db:"min_contract_duration"`
	MaxContractDuration int64                `db:"max_contract_duration"`
	SubscriptionRate    int64                `db:"subscription_rate"`
	PayAsYouGoRate      int64                `db:"paygo_rate"`
}

func (d *DirectoryDB) InsertProvider(provider *ArkeoProvider) (*Entity, error) {
	if provider == nil {
		return nil, fmt.Errorf("nil provider")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return insert(conn, sqlInsertProvider, provider.Pubkey, provider.Chain, provider.Bond)
}

func (d *DirectoryDB) UpdateProvider(provider *ArkeoProvider) (*Entity, error) {
	if provider == nil {
		return nil, fmt.Errorf("nil provider")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return update(conn,
		sqlUpdateProvider,
		provider.Pubkey,
		provider.Chain,
		provider.Bond,
		provider.MetadataURI,
		provider.MetadataNonce,
		provider.Status,
		provider.MinContractDuration,
		provider.MaxContractDuration,
		provider.SubscriptionRate,
		provider.PayAsYouGoRate,
	)
}

func (d *DirectoryDB) FindProvider(pubkey string, chain string) (*ArkeoProvider, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	provider := ArkeoProvider{}
	if err = selectOne(conn, sqlFindProvider, &provider, pubkey, chain); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}
	// not found
	if provider.Pubkey == "" {
		return nil, nil
	}
	return &provider, nil
}

func (d *DirectoryDB) SearchProviders(criteria types.ProviderSearchParams) ([]*ArkeoProvider, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	sb := sqlbuilder.NewSelectBuilder()

	sb.Select("id", "created", "pubkey", "chain", "status", "metadata_uri", "metadata_nonce",
		"subscription_rate", "paygo_rate", "min_contract_duration", "max_contract_duration", "bond").
		From("providers")

	if criteria.Pubkey != "" {
		sb = sb.Where(sb.Equal("pubkey", criteria.Pubkey))
	}
	if criteria.Chain != "" {
		sb = sb.Where(sb.Equal("chain", criteria.Chain))
	}

	sql, params := sb.BuildWithFlavor(getFlavor())
	log.Debugf("sql: %s\n%v", sql, params)

	providers := make([]*ArkeoProvider, 0, 512)
	if err := pgxscan.Select(context.Background(), conn, &providers, sql, params...); err != nil {
		return nil, errors.Wrapf(err, "error selecting many")
	}

	return providers, nil
}

func (d *DirectoryDB) InsertBondProviderEvent(providerID int64, evt types.BondProviderEvent) (*Entity, error) {
	if evt.BondAbsolute == nil {
		return nil, fmt.Errorf("nil BondAbsolute")
	}
	if evt.BondRelative == nil {
		return nil, fmt.Errorf("nil BondRelative")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return insert(conn, sqlInsertBondProviderEvent, providerID, evt.TxID, evt.BondRelative.String(), evt.BondAbsolute.String())
}

func (d *DirectoryDB) InsertModProviderEvent(providerID int64, evt types.ModProviderEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return insert(conn, sqlInsertModProviderEvent, providerID, evt.TxID, evt.MetadataURI, evt.MetadataNonce, evt.Status,
		evt.MinContractDuration, evt.MaxContractDuration, evt.SubscriptionRate, evt.PayAsYouGoRate)
}

func (d *DirectoryDB) UpsertProviderMetadata(providerID int64, data sentinel.Metadata) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	c := data.Configuration
	return upsert(conn, sqlUpsertProviderMetadata, providerID, data.Version, c.Moniker, c.Website, c.Description, c.Location,
		c.Port, c.ProxyHost, c.SourceChain, c.EventStreamHost, c.ClaimStoreLocation, c.FreeTierRateLimit, c.FreeTierRateLimitDuration,
		c.SubTierRateLimit, c.SubTierRateLimitDuration, c.AsGoTierRateLimit, c.AsGoTierRateLimitDuration)
}
