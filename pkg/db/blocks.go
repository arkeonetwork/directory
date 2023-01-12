package db

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type Block struct {
	Entity
	Height    int64     `db:"height"`
	Hash      string    `db:"hash"`
	BlockTime time.Time `db:"block_time"`
}

func (d *DirectoryDB) InsertBlock(b *Block) (*Entity, error) {
	if b == nil {
		return nil, fmt.Errorf("nil block")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	return insert(conn, sqlInsertBlock, b.Height, b.Hash, b.BlockTime)
}

func (d *DirectoryDB) FindLatestBlock() (*Block, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	block := &Block{} // used to designate not found... need a better way!
	if err = selectOne(conn, sqlFindLatestBlock, block); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}
	// not found
	// if block.Height == math.MaxUint64 {
	// 	return nil, nil
	// }
	return block, nil
}
