package db

import (
	"github.com/ArkeoNetwork/directory/pkg/types"
	"github.com/pkg/errors"
)

func (d *DirectoryDB) UpsertContract(providerID int64, evt types.OpenContractEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertContract, providerID, evt.DelegatePubkey, evt.ClientPubkey, evt.ContractType,
		evt.Duration, evt.Rate, evt.OpenCost)
}
