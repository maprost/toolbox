package badgerx_test

import (
	"testing"

	"github.com/maprost/testbox/must"
	"github.com/maprost/toolbox/badgerx"
)

func newDB(t testing.TB, path string) (badgerx.DB, func()) {
	db, err := badgerx.NewDB(path)
	must.BeNoError(t, err)
	must.NotBeNil(t, db)

	return db, func() { must.BeNoError(t, db.Close()) }
}

func TestDeepPath(t *testing.T) {
	db, err := badgerx.NewDB("TestDeepPath/blob/drop/mopp")
	must.BeNoError(t, err)
	must.NotBeNil(t, db)

	db.Close()
}

func TestTwoOpenDB(t *testing.T) {
	db1, err := badgerx.NewDB("TestTwoOpenDB")
	must.BeNoError(t, err)
	must.NotBeNil(t, db1)
	defer db1.Close()

	// can't open a db twice
	_, err = badgerx.NewDB("TestTwoOpenDB")
	must.BeError(t, err)
}
