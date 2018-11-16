package pqx

import (
	"database/sql"
	"time"
)

type Result struct {
	SQL      string
	Duration time.Duration
	Err      error
}

type RowsResult struct {
	Result
	rows *sql.Rows
}

func (r *RowsResult) Close() {
	if r.rows != nil {
		r.rows.Close()
	}
}

func (r *RowsResult) Next() bool {
	if r.rows != nil {
		return r.rows.Next()
	}
	return false
}

func (r *RowsResult) Scan(dest ...interface{}) error {
	if r.rows != nil {
		return r.rows.Scan(dest...)
	}
	return nil
}

func (r *RowsResult) ScanEntity(entity Entity) (bool, error) {
	if r.rows == nil {
		return false, nil
	}

	config := entity.EntityConfig()

	// prepare scan
	dest := make([]interface{}, len(config.Columns))
	for i, column := range config.Columns {
		value, err := entity.PtrValue(column)
		if err != nil {
			return false, err
		}

		dest[i] = value
	}

	// scan
	if r.Next() {
		err := r.Scan(dest...)
		if err != nil {
			return false, err
		}
	} else {
		// no error!
		return false, nil
	}

	return true, nil
}

func resultError(result Result, err error) Result {
	return Result{
		SQL:      result.SQL,
		Duration: result.Duration,
		Err:      err,
	}
}
