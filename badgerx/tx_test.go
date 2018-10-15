package badgerx_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/maprost/testbox/is"
	"github.com/maprost/testbox/must"
	"github.com/maprost/testbox/should"
	"github.com/maprost/toolbox/badgerx"
	. "github.com/maprost/toolbox/badgerx/benchmark"
)

func newTx(t testing.TB, db badgerx.DB) (badgerx.Tx, func()) {
	tx := db.Tx()
	return tx, func() { must.BeNoError(t, tx.Close(true)) }
}

func TestBadger_cantFindKeyInOlderTransaction(t *testing.T) {
	key := []byte("hello")
	value := []byte("welt")

	// ------------------------
	// ----------- Setup ------
	// ------------------------
	rootPath, err := os.Getwd()
	must.BeNoError(t, err)

	path := "badger/TestBadger_cantFindKeyInOlderTransaction"

	// remove db
	fullPath := rootPath + "/" + path
	err = os.RemoveAll(fullPath)
	must.BeNoError(t, err)

	// create folder
	err = os.MkdirAll(fullPath, os.ModePerm)
	must.BeNoError(t, err)

	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path

	db, err := badger.Open(opts)
	must.BeNoError(t, err)

	beforeTx := db.NewTransaction(false)

	// write something
	writeTx := db.NewTransaction(true)
	writeTx.Set(key, value)

	err = writeTx.Commit(nil)
	must.BeNoError(t, err)

	// check --> can't find key in older transaction
	_, err = beforeTx.Get(key)
	must.BeError(t, err)
	must.BeEqual(t, err, badger.ErrKeyNotFound)
}

func TestMultiplyReadCallOnOneTransaction(t *testing.T) {
	db, cleanUp := newDB(t, "TestKeyIDWorkflow")
	defer cleanUp()

	// Start a transaction.
	tx := db.Tx()

	// create an entry
	err := tx.UpsertKeyID("answer", 12)
	must.BeNoError(t, err)

	// Commit the transaction and check for error.
	err = tx.Close(true)
	must.BeNoError(t, err)

	t.Run("check multiply read routines", func(t *testing.T) {
		// Start a new transaction.
		tx := db.Tx()

		goRoutines := 10
		done := make(chan error)

		for i := 0; i < goRoutines; i++ {
			go func(i int) {
				for j := 0; j < 10; j++ {
					val, found, err := tx.GetKeyID("answer")
					if is.Error(err) {
						done <- err
						return
					}

					if !found {
						done <- fmt.Errorf("not found")
						return
					}

					if notEqual, err := is.NotEqualf(val, uint64(42)); notEqual {
						done <- err
						return
					}
				}
				done <- nil
			}(i)
		}

		for i := 0; i < goRoutines; i++ {
			should.BeNoError(t, <-done)
		}
	})
}

func TestMultiplyWriteTransactions(t *testing.T) {
	db, cleanUp := newDB(t, "TestKeyIDWorkflow")
	defer cleanUp()

	goRoutines := 100
	errorChan := make(chan error)

	for i := 0; i < goRoutines; i++ {
		go func(i int) {
			tx := db.Tx()

			err := tx.UpsertKeyID(strconv.Itoa(i), uint64(i))
			if err != nil {
				errorChan <- err
			}

			err = tx.Close(true)
			if err != nil {
				errorChan <- err
			}

			errorChan <- nil
		}(i)
	}

	for i := 0; i < goRoutines; i++ {
		should.BeNoError(t, <-errorChan)
	}
}

func TestKeyIDWorkflow(t *testing.T) {
	db, cleanUpDB := newDB(t, "TestKeyIDWorkflow")
	defer cleanUpDB()

	tx, cleanUpTx := newTx(t, db)
	defer cleanUpTx()

	must.BeNoError(t, tx.UpsertKeyID("a", 1))
	must.BeNoError(t, tx.UpsertKeyID("b", 2))

	id, found, err := tx.GetKeyID("a")
	must.BeNoError(t, err)
	must.BeTrue(t, found)
	should.BeEqual(t, id, uint64(1))

	id, found, err = tx.GetKeyID("b")
	must.BeNoError(t, err)
	must.BeTrue(t, found)
	should.BeEqual(t, id, uint64(2))

	id, found, err = tx.GetKeyID("c")
	must.BeNoError(t, err)
	should.BeFalse(t, found)

	// update 'a'
	must.BeNoError(t, tx.UpsertKeyID("a", 42))

	// delete 'b'
	must.BeNoError(t, tx.DeleteKey("b"))

	id, found, err = tx.GetKeyID("a")
	must.BeNoError(t, err)
	must.BeTrue(t, found)
	should.BeEqual(t, id, uint64(42))

	id, found, err = tx.GetKeyID("b")
	must.BeNoError(t, err)
	should.BeFalse(t, found)
}

func TestIDStructWorkflow(t *testing.T) {
	db, cleanUpDB := newDB(t, "TestIDStructWorkflow")
	defer cleanUpDB()

	tx, cleanUpTx := newTx(t, db)
	defer cleanUpTx()

	s1 := NewTestStruct()
	s1.Msg = []byte("hello")
	must.BeNoError(t, tx.InsertIDStruct(s1))
	must.NotBeEqual(t, s1.ID(), 0)

	s2 := NewTestStruct()
	s2.Msg = []byte("World")
	must.BeNoError(t, tx.InsertIDStruct(s2))
	must.NotBeEqual(t, s2.ID(), 0)

	{
		var out TestStructWrapper
		found, err := tx.GetIDStruct(s1.ID(), &out)
		must.BeNoError(t, err)
		must.BeTrue(t, found)
		should.BeEqual(t, &out, s1)
	}

	{
		var out TestStructWrapper
		found, err := tx.GetIDStruct(s2.ID(), &out)
		must.BeNoError(t, err)
		must.BeTrue(t, found)
		should.BeEqual(t, &out, s2)
	}

	{
		var out TestStructWrapper
		found, err := tx.GetIDStruct(s2.ID()+1, &out)
		must.BeNoError(t, err)
		should.BeFalse(t, found)
	}

	// update 's1'
	s1.Msg = []byte("Hallo Welt")
	must.BeNoError(t, tx.UpdateIDStruct(s1))

	// delete 's2'
	must.BeNoError(t, tx.DeleteID(s2.ID()))

	{
		var out TestStructWrapper
		found, err := tx.GetIDStruct(s1.ID(), &out)
		must.BeNoError(t, err)
		must.BeTrue(t, found)
		should.BeEqual(t, &out, s1)
	}

	{
		var out TestStructWrapper
		found, err := tx.GetIDStruct(s2.ID(), &out)
		must.BeNoError(t, err)
		should.BeFalse(t, found)
	}
}

func TestIDStructUniqueIDs(t *testing.T) {
	db, cleanUpDB := newDB(t, "TestIDStructWorkflow")
	defer cleanUpDB()

	goRoutines := 100
	idChan := make(chan int64)

	for i := 0; i < goRoutines; i++ {
		go func(i int) {
			tx := db.Tx()

			s := NewTestStruct()
			s.Msg = []byte(fmt.Sprintf("%d", i))
			err := tx.InsertIDStruct(s)
			if err != nil {
				idChan <- -1
			}

			err = tx.Close(true)
			if err != nil {
				idChan <- -2
			}

			idChan <- int64(s.ID())
		}(i)
	}

	// get all ids
	ids := make([]int64, 0, goRoutines)
	for i := 0; i < goRoutines; i++ {
		id := <-idChan
		if id < 0 {
			should.Fail(t, "Error inside go routine: ", id)
		}

		ids = append(ids, id)
	}

	// check duplicates
	for i := 0; i < goRoutines; i++ {
		for j := i + 1; j < goRoutines; j++ {
			should.NotBeEqual(t, ids[i], ids[j], "Error duplicate ids created ")
		}
	}
}

func TestKeyStructWorkflow(t *testing.T) {
	db, cleanUpDB := newDB(t, "TestKeyStructWorkflow")
	defer cleanUpDB()

	tx, cleanUpTx := newTx(t, db)
	defer cleanUpTx()

	s1 := NewKeyTestStruct("a")
	must.BeNoError(t, tx.UpsertKeyStruct(s1))

	s2 := NewKeyTestStruct("b")
	must.BeNoError(t, tx.UpsertKeyStruct(s2))

	{
		var out TestStructWrapper
		found, err := tx.GetKeyStruct(s1.Key(), &out)
		must.BeNoError(t, err)
		must.BeTrue(t, found)
		must.BeEqual(t, &out, s1)
	}

	{
		var out TestStructWrapper
		found, err := tx.GetKeyStruct(s2.Key(), &out)
		must.BeNoError(t, err)
		must.BeTrue(t, found)
		must.BeEqual(t, &out, s2)
	}

	{
		var out TestStructWrapper
		found, err := tx.GetKeyStruct("c", &out)
		must.BeNoError(t, err)
		must.BeFalse(t, found)
	}

	// update 's1'
	s1.Msg = []byte("Hallo Welt")
	tx.UpsertKeyStruct(s1)

	// delete 's2'
	tx.DeleteKey(s2.Key())

	{
		var out TestStructWrapper
		found, err := tx.GetKeyStruct(s1.Key(), &out)
		must.BeNoError(t, err)
		must.BeTrue(t, found)
		must.BeEqual(t, &out, s1)
	}

	{
		var out TestStructWrapper
		found, err := tx.GetKeyStruct(s2.Key(), &out)
		must.BeNoError(t, err)
		must.BeFalse(t, found)
	}
}

func TestTTL(t *testing.T) {
	db, cleanUpDB := newDB(t, "TestKeyStructWorkflow")
	defer cleanUpDB()

	// create item
	createTx := db.Tx()
	createTx.SetTTL(time.Second)

	createTx.UpsertKeyStruct(NewKeyTestStruct("blob"))
	must.BeNoError(t, createTx.Close(true))

	// check item
	checkTx := db.Tx()
	var check TestStructWrapper
	found, err := checkTx.GetKeyStruct("blob", &check)
	must.BeNoError(t, err)
	must.BeTrue(t, found)

	time.Sleep(1 * time.Second)

	t.Run("check with new transaction", func(t *testing.T) {
		checkTx := db.Tx()
		found, err = checkTx.GetKeyStruct("blob", nil)
		must.BeNoError(t, err)
		must.BeFalse(t, found)

		must.BeNoError(t, checkTx.Close(false))
	})

	t.Run("check with old transaction", func(t *testing.T) {
		found, err := checkTx.GetKeyStruct("blob", nil)
		must.BeNoError(t, err)
		must.BeFalse(t, found)
	})

	must.BeNoError(t, checkTx.Close(false))
}
