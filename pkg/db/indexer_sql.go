package db

var (
	sqlInsertIndexerStatus = `insert into indexer_status(id,height) values ($1,$2) returning id, created, updated`
	sqlUpdateIndexerStatus = `update indexer_status set height = $2, updated = now() where id = $1 returning id, created, updated`
	sqlFindIndexerStatus   = `select id,height from indexer_status where id = $1`
)
