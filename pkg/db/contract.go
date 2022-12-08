package db

import (
	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/pkg/errors"
)

type ArkeoContract struct {
	Entity
	ProviderID     int64              `db:"provider_id"`
	DelegatePubkey string             `db:"delegate_pubkey"`
	ClientPubkey   string             `db:"client_pubkey"`
	Height         int64              `db:"height"`
	ContractType   types.ContractType `db:"contract_type"`
	Duration       int64              `db:"duration"`
	Rate           int64              `db:"rate"`
	OpenCost       int64              `db:"open_cost"`
}

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

func (d *DirectoryDB) UpsertContractSettlementEvent(contractID int64, evt types.ContractSettlementEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertContractSettlementEvent, contractID, evt.TxID, evt.ClientPubkey, evt.Height,
		evt.Nonce, evt.Paid, evt.Reserve)
}

func (d *DirectoryDB) UpsertOpenContractEvent(contractID int64, evt types.OpenContractEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return insert(conn, sqlUpsertOpenContractEvent, contractID, evt.ClientPubkey, evt.ContractType, evt.Height, evt.TxID,
		evt.Duration, evt.Rate, evt.OpenCost)
}
