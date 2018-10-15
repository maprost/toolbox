package badgerx_test

import (
	"testing"

	"github.com/maprost/testbox/must"
	"github.com/maprost/toolbox/badgerx"
)

func TestMapWorkflow(t *testing.T) {
	db1 := "TestMapWorkflow1"
	db2 := "TestMapWorkflow2"
	db3 := "TestMapWorkflow3"

	dbMap, err := badgerx.NewDBMap([]string{db1, db2, db3})
	must.BeNoError(t, err)

	{
		// insert something
		txMap := dbMap.TxMap()
		must.BeNoError(t, txMap[db1].UpsertKeyID("blob", 1))
		txMap.Close(true)
	}

	{
		// get it back
		txMap := dbMap.TxMap()
		id, found, err := txMap[db1].GetKeyID("blob")
		must.BeNoError(t, err)
		must.BeTrue(t, found)
		must.BeEqual(t, id, uint64(1))
		txMap.Close(false)
	}

	must.BeNoError(t, dbMap.Close())
}
