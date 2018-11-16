package pqx

import (
	"database/sql"
)

type Transaction struct {
	tx *sql.Tx
}

func NewTx() (Transaction, error) {
	return newTx(singleDB)
}

func (db *DB) NewTx() (Transaction, error) {
	return newTx(db.db)
}

func newTx(db *sql.DB) (Transaction, error) {
	tx, err := db.Begin()

	return Transaction{
		tx: tx,
	}, err
}

func (tx *Transaction) Query(sql SQL) RowsResult {
	return query(tx.tx, sql)
}

func (tx *Transaction) Commit() error {
	if tx.tx == nil {
		return nil
	}

	err := tx.tx.Commit()
	if err != nil {
		return err
	}

	tx.tx = nil
	return nil
}

func (tx *Transaction) Rollback() error {
	if tx.tx == nil {
		return nil
	}

	err := tx.tx.Rollback()
	if err != nil {
		return err
	}

	tx.tx = nil
	return nil
}
