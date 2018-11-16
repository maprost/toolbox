package pqx_test

import (
	"testing"
	"time"

	"github.com/maprost/testbox/must"
	"github.com/maprost/testbox/should"
	"github.com/maprost/toolbox/pqx"
	. "github.com/maprost/toolbox/pqx/benchmark"
)

func TestOpenDatabaseConnection(t *testing.T) {
	err := pqx.OpenSingleDatabaseConnection(pqx.NewDefaultConnectionInfo("db"))
	must.BeNoError(t, err)
}

func TestInsert(t *testing.T) {
	err := pqx.OpenSingleDatabaseConnection(pqx.NewDefaultConnectionInfo("db"))
	must.BeNoError(t, err)

	CreateTestEntityTableIfNotExists(t, pqx.Query)

	t0 := time.Now()
	testEntity := TestEntity{
		ID:      0,
		Message: "hello world",
		Changed: time.Time{},
		Created: time.Time{},
	}
	result := pqx.Insert(&testEntity)
	must.BeNoError(t, result.Err)

	t1 := time.Now()
	should.NotBeEqual(t, testEntity.ID, 0)
	should.BeEqual(t, testEntity.Message, "hello world")

	// changed is in time range [t0:t1]
	should.NotBeEqual(t, testEntity.Changed, time.Time{})
	should.BeTrue(t, testEntity.Changed.Equal(t0) || testEntity.Changed.After(t0))
	should.BeTrue(t, testEntity.Changed.Equal(t1) || testEntity.Changed.Before(t1))

	// created is in time range [t0:t1]
	should.NotBeEqual(t, testEntity.Created, time.Time{})
	should.BeTrue(t, testEntity.Created.Equal(t0) || testEntity.Created.After(t0))
	should.BeTrue(t, testEntity.Created.Equal(t1) || testEntity.Created.Before(t1))
}

func TestContains(t *testing.T) {
	err := pqx.OpenSingleDatabaseConnection(pqx.NewDefaultConnectionInfo("db"))
	must.BeNoError(t, err)

	CreateTestEntityTableIfNotExists(t, pqx.Query)

	testEntity := TestEntity{Message: "hello world"}

	// insert
	result := pqx.Insert(&testEntity)
	must.BeNoError(t, result.Err)

	// contains
	found, result := pqx.Contains(&TestEntity{ID: testEntity.ID})
	must.BeNoError(t, result.Err)
	should.BeTrue(t, found)

	// not found
	found, result = pqx.Contains(&TestEntity{ID: testEntity.ID + 1})
	must.BeNoError(t, result.Err)
	should.BeFalse(t, found)
}

func TestDelete(t *testing.T) {
	err := pqx.OpenSingleDatabaseConnection(pqx.NewDefaultConnectionInfo("db"))
	must.BeNoError(t, err)

	CreateTestEntityTableIfNotExists(t, pqx.Query)

	testEntity := &TestEntity{Message: "hello world"}

	// insert
	result := pqx.Insert(testEntity)
	must.BeNoError(t, result.Err)

	// contains
	found, result := pqx.Contains(testEntity)
	must.BeNoError(t, result.Err)
	should.BeTrue(t, found)

	// delete
	result = pqx.Delete(testEntity)
	must.BeNoError(t, result.Err)

	// not found
	found, result = pqx.Contains(testEntity)
	must.BeNoError(t, result.Err)
	should.BeFalse(t, found)
}

func TestSelect(t *testing.T) {
	err := pqx.OpenSingleDatabaseConnection(pqx.NewDefaultConnectionInfo("db"))
	must.BeNoError(t, err)

	CreateTestEntityTableIfNotExists(t, pqx.Query)

	testEntity := TestEntity{Message: "hello world"}

	// insert
	t0 := time.Now()
	result := pqx.Insert(&testEntity)
	t1 := time.Now()
	must.BeNoError(t, result.Err)

	// select
	{
		selectedEntity := TestEntity{ID: testEntity.ID}
		found, result := pqx.FillEntity(&selectedEntity)
		must.BeNoError(t, result.Err)
		should.BeTrue(t, found)

		should.BeEqual(t, selectedEntity.ID, testEntity.ID)
		should.BeEqual(t, selectedEntity.Message, "hello world")

		should.NotBeEqual(t, selectedEntity.Changed, time.Time{})
		should.BeTrue(t, selectedEntity.Changed.Equal(t0) || selectedEntity.Changed.After(t0))
		should.BeTrue(t, selectedEntity.Changed.Equal(t1) || selectedEntity.Changed.Before(t1))

		should.NotBeEqual(t, selectedEntity.Created, time.Time{})
		should.BeTrue(t, selectedEntity.Created.Equal(t0) || selectedEntity.Created.After(t0))
		should.BeTrue(t, selectedEntity.Created.Equal(t1) || selectedEntity.Created.Before(t1))
	}

	// not found
	{
		notFoundEntity := TestEntity{ID: testEntity.ID + 1}
		found, result := pqx.FillEntity(&notFoundEntity)
		must.BeNoError(t, result.Err)
		should.BeFalse(t, found)
	}
}

func TestUpdate(t *testing.T) {
	err := pqx.OpenSingleDatabaseConnection(pqx.NewDefaultConnectionInfo("db"))
	must.BeNoError(t, err)

	CreateTestEntityTableIfNotExists(t, pqx.Query)

	testEntity := &TestEntity{Message: "hello world"}

	// insert
	t0 := time.Now()
	result := pqx.Insert(testEntity)
	t1 := time.Now()
	must.BeNoError(t, result.Err)

	// select
	{
		selectedEntity := TestEntity{ID: testEntity.ID}
		found, result := pqx.FillEntity(&selectedEntity)
		must.BeNoError(t, result.Err)
		should.BeTrue(t, found)

		should.BeEqual(t, selectedEntity.ID, testEntity.ID)
		should.BeEqual(t, selectedEntity.Message, "hello world")
	}
	// to see a change in the times
	time.Sleep(10 * time.Millisecond)

	// update
	updateEntity := &TestEntity{ID: testEntity.ID, Message: "rule the world"}
	t2 := time.Now()
	result = pqx.Update(updateEntity)
	t3 := time.Now()
	must.BeNoError(t, result.Err)

	// select again
	{
		selectedAgainEntity := TestEntity{ID: testEntity.ID}
		found, result := pqx.FillEntity(&selectedAgainEntity)
		must.BeNoError(t, result.Err)
		should.BeTrue(t, found)

		should.BeEqual(t, selectedAgainEntity.ID, testEntity.ID)
		should.BeEqual(t, selectedAgainEntity.Message, "rule the world")

		// changed is in the time range: [t2:t3]
		should.NotBeEqual(t, selectedAgainEntity.Changed, time.Time{})
		should.BeTrue(t, selectedAgainEntity.Changed.Equal(t2) || selectedAgainEntity.Changed.After(t2))
		should.BeTrue(t, selectedAgainEntity.Changed.Equal(t3) || selectedAgainEntity.Changed.Before(t3))

		// created stays in the time range [t0:t1]
		should.NotBeEqual(t, selectedAgainEntity.Created, time.Time{})
		should.BeTrue(t, selectedAgainEntity.Created.Equal(t0) || selectedAgainEntity.Created.After(t0))
		should.BeTrue(t, selectedAgainEntity.Created.Equal(t1) || selectedAgainEntity.Created.Before(t1))
	}
}
