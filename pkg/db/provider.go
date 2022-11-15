package db

import (
	"fmt"

	"github.com/ArkeoNetwork/directory/pkg/types"
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
