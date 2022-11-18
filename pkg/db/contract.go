package db

import (
	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/pkg/errors"
)

type ArkeoContract struct {
	Entity
	ProviderID     int64              `db:"provider_id"`
	DelegatePubkey string             `db:"bond"`
	ClientPubkey   string             `db:"bond"`
	Height         int64              `db:"bond"`
	ContractType   types.ContractType `db:"bond"`
	Duration       int64              `db:"bond"`
	Rate           int64              `db:"bond"`
	OpenCost       int64              `db:"bond"`
}

/*
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
*/

func (d *DirectoryDB) FindContract(providerID int64, delegatePubkey string) (*ArkeoContract, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	contract := ArkeoContract{}
	if err = selectOne(conn, sqlFindContract, &contract, providerID, delegatePubkey); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}

	// not found
	if contract.ID == 0 {
		return nil, nil
	}
	return &contract, nil
}

func (d *DirectoryDB) UpsertContract(providerID int64, evt types.OpenContractEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertContract, providerID, evt.DelegatePubkey, evt.ClientPubkey, evt.ContractType,
		evt.Duration, evt.Rate, evt.OpenCost, evt.Height)
}

func (d *DirectoryDB) InsertOpenContractEvent(contractID int64, evt types.OpenContractEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return insert(conn, sqlInsertOpenContractEvent, contractID, evt.ClientPubkey, evt.ContractType, evt.Height, evt.TxID,
		evt.Duration, evt.Rate, evt.OpenCost)
}
