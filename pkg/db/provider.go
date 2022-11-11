package db

import (
	"fmt"

	"github.com/pkg/errors"
)

type ArkeoProvider struct {
	Entity
	Pubkey string `db:"pubkey"`
	Chain  string `db:"chain"`
	Bond   string `db:"bond"` // this is a DECIMAL type in the db
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
