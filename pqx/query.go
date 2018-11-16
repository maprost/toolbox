package pqx

import (
	"database/sql"
	"time"
)

type originExecutor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func query(exe originExecutor, sql SQL) RowsResult {
	t0 := time.Now()
	rows, err := exe.Query(sql.sql, sql.args...)

	return RowsResult{
		Result: Result{
			SQL:      sql.String(),
			Duration: time.Since(t0),
			Err:      err,
		},
		rows: rows,
	}
}
