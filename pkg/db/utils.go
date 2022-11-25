package db

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func insert(conn *pgxpool.Conn, sql string, params ...interface{}) (*Entity, error) {
	var (
		id      int64
		created time.Time
		updated time.Time
		err     error
	)
	log.Debugf("sql: %s", sql)
	row := conn.QueryRow(context.Background(), sql, params...)
	if err = row.Scan(&id, &created, &updated); err != nil {
		return nil, errors.Wrap(err, "error inserting")
	}

	return &Entity{ID: id, Created: created, Updated: updated}, nil
}

func update(conn *pgxpool.Conn, sql string, params ...interface{}) (*Entity, error) {
	var (
		id      int64
		created time.Time
		updated time.Time
		err     error
	)
	log.Debugf("sql: %s", sql)
	row := conn.QueryRow(context.Background(), sql, params...)
	if err = row.Scan(&id, &created, &updated); err != nil {
		return nil, errors.Wrap(err, "error inserting")
	}

	return &Entity{ID: id, Created: created, Updated: updated}, nil
}

// if the query returns no rows, the passed target remains unchanged. target must be a pointer
func selectOne(conn *pgxpool.Conn, sql string, target interface{}, params ...interface{}) error {
	if err := pgxscan.Get(context.Background(), conn, target, sqlFindProvider, params...); err != nil {
		unwrapped := errors.Unwrap(err)
		if unwrapped != nil && unwrapped.Error() == "no rows in result set" {
			return nil
		}
		return errors.Wrapf(err, "error selecting with params: %v", params)
	}
	return nil
}

func upsert(conn *pgxpool.Conn, sql string, params ...interface{}) (*Entity, error) {
	row := conn.QueryRow(context.Background(), sql, params...)

	var (
		id      int64
		created time.Time
		updated time.Time
		err     error
	)

	if err = row.Scan(&id, &created, &updated); err != nil {
		return nil, fmt.Errorf("error upserting: %+v", err)
	}

	entity := &Entity{
		ID:      id,
		Created: created,
		Updated: updated,
	}

	return entity, nil
}
