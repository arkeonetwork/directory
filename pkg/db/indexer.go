package db

import (
	"fmt"

	"github.com/pkg/errors"
)

type IndexerStatus struct {
	ID     string `db:"id"`
	Height uint64 `db:"height"`
}

func (d *DirectoryDB) InsertIndexerStatus(indexerStatus *IndexerStatus) (*Entity, error) {
	if indexerStatus == nil {
		return nil, fmt.Errorf("nil IndexerStatus")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return insert(conn, sqlInsertIndexerStatus, indexerStatus.ID, indexerStatus.Height)
}

func (d *DirectoryDB) UpdateIndexerStatus(indexerStatus *IndexerStatus) (*Entity, error) {
	if indexerStatus == nil {
		return nil, fmt.Errorf("nil IndexerStatus")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return update(conn,
		sqlUpdateIndexerStatus,
		indexerStatus.ID,
		indexerStatus.Height,
	)
}

func (d *DirectoryDB) FindIndexerStatus(id string) (*IndexerStatus, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	indexerStatus := IndexerStatus{}
	if err = selectOne(conn, sqlFindIndexerStatus, &indexerStatus, id); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}
	// not found
	if indexerStatus.ID == "" {
		return nil, nil
	}
	return &indexerStatus, nil
}
