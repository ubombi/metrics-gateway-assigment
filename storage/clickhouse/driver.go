package clickhouse

import (
	"database/sql"

	"github.com/pkg/errors"
)

func connect(dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("clickhouse", dsn)
	if err != nil {
		err = errors.Wrap(err, "can not connect to clickhouse")
		return
	}

	err = errors.Wrap(db.Ping(), "failed to ping clickhouse")
	return
}

func makeInsert(db *sql.DB) (tx *sql.Tx, stmt *sql.Stmt, err error) {
	if tx, err = db.Begin(); err != nil {
		return nil, nil, errors.Wrap(err, "can not start transaction")
	}
	if stmt, err = tx.Prepare(insertQuery); err != nil {
		return nil, nil, errors.Wrap(err, "can not prepare insert query")
	}

	return tx, stmt, nil
}
