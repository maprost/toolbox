package badgerx

import (
	"os"

	"github.com/dgraph-io/badger"
	"sync"
)

type DB struct {
	db    *badger.DB
	mutex *sync.Mutex
}

func NewDB(path string) (DB, error) {
	db, err := openDB(path)
	if err != nil {
		return DB{}, err
	}

	return DB{
		db:    db,
		mutex: &sync.Mutex{},
	}, nil
}

func (u DB) Tx() Tx {
	return Tx{
		db:          u,
		readableTx:  nil,
		writeableTx: nil,
	}
}

func (u DB) Close() error {
	return u.db.Close()
}

func (u DB) nextID() (uint64, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	seq, err := u.db.GetSequence([]byte("ID"), 1)
	if err != nil {
		return 0, err
	}

	return seq.Next()
}

func (u DB) newReadableTransaction() *badger.Txn {
	return u.db.NewTransaction(false)
}

func (u DB) newWriteableTransaction() *badger.Txn {
	return u.db.NewTransaction(true)
}

func openDB(path string) (*badger.DB, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path = "badger/" + path

	fullPath := rootPath + "/" + path
	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path

	return badger.Open(opts)
}
