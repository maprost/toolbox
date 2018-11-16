package speed

import (
	"testing"

	"github.com/maprost/testbox/must"
	"github.com/maprost/toolbox/pqx"
	. "github.com/maprost/toolbox/pqx/benchmark"
)

// 1000	   2223374 ns/op	    3328 B/op	      97 allocs/op
// 5000	    320293 ns/op	    2009 B/op	      58 allocs/op
func BenchmarkSpeed(b *testing.B) {
	err := pqx.OpenSingleDatabaseConnection(pqx.NewDefaultConnectionInfo("db"))
	must.BeNoError(b, err)

	CreateTestEntityTableIfNotExists(b, pqx.Query)

	testEntity := TestEntity{Message: "hello world"}

	// test insert
	result := pqx.Insert(&testEntity)
	must.BeNoError(b, result.Err)

	startID := testEntity.ID
	lastID := testEntity.ID

	// test select
	selectedEntity := TestEntity{ID: testEntity.ID}
	found, result := pqx.FillEntity(&selectedEntity)
	must.BeNoError(b, result.Err)
	must.BeTrue(b, found)

	b.Run("Insert", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			testEntity := TestEntity{Message: "hello world"}
			pqx.Insert(&testEntity)

			lastID = testEntity.ID
		}
	})

	b.Run("Select", func(b *testing.B) {
		b.ReportAllocs()
		id := startID
		for i := 0; i < b.N; i++ {
			id++
			if id > lastID {
				id = startID
			}

			selectedEntity := TestEntity{ID: id}
			pqx.FillEntity(&selectedEntity)
		}
	})
}
