package speed

import (
	"fmt"
	"testing"

	"github.com/maprost/testbox/must"
	"github.com/maprost/toolbox/badgerx"
	"github.com/maprost/toolbox/badgerx/benchmark"
)

// #### BenchmarkDatabaseJson/default
// 500	   		3793245 ns/op	    2748 B/op	      71 allocs/op
// 1000000	    2017 ns/op	     	536 B/op	       8 allocs/op
//
// #### BenchmarkDatabaseJson/Jsoniter
// 500	   		3687072 ns/op	    2460 B/op	      67 allocs/op
// 1000000	    1699 ns/op	     	360 B/op	      12 allocs/op
//
// #### BenchmarkDatabaseJson/msgPack
// 500	   		3683173 ns/op	    2947 B/op	      70 allocs/op
// 1000000	    1357 ns/op	     	536 B/op	       6 allocs/op
//
// #### BenchmarkDatabaseJson/easyJson
// 500	   		3742288 ns/op	    2129 B/op	      65 allocs/op
// 2000000	    881 ns/op	     	200 B/op	       3 allocs/op
//
// #### BenchmarkDatabaseJson/unsafePointer
// 500	   		3713792 ns/op	    2018 B/op	      65 allocs/op
// 10000000	    125 ns/op	     	184 B/op	       2 allocs/op
func BenchmarkDatabaseJson(b *testing.B) {
	db, err := badgerx.NewDB("BenchmarkDatabaseJson")
	must.BeNoError(b, err)
	must.NotBeNil(b, db)

	tx := db.Tx()

	for _, codec := range benchmark.CodecList {
		b.Run(codec.Name, func(b *testing.B) {
			fmt.Println("####", b.Name())
			in := benchmark.NewTestStruct()
			in.EncodeFunc = codec.Encoding

			err := tx.InsertIDStruct(in)
			must.BeNoError(b, err)

			var out benchmark.TestStructWrapper
			out.DecodeFunc = codec.Decoding

			found, err := tx.GetIDStruct(in.ID(), &out)
			must.BeNoError(b, err)
			must.BeTrue(b, found)
			must.BeEqual(b, in.TestStruct, out.TestStruct)

			b.Run("Insert", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					tx.InsertIDStruct(in)
				}
			})

			b.Run("Select", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					tx.GetIDStruct(in.ID(), &out)
				}
			})
		})
	}

	must.BeNoError(b, tx.Close(true))
	must.BeNoError(b, db.Close())
}

func BenchmarkJson(b *testing.B) {
	for _, codec := range benchmark.CodecList {
		b.Run(codec.Name, func(b *testing.B) {
			fmt.Println("####", b.Name())
			in := benchmark.NewTestStruct()

			bytes, err := codec.Encoding(*in)
			must.BeNoError(b, err)

			var out benchmark.TestStructWrapper
			err = codec.Decoding(bytes, &out)
			must.BeNoError(b, err)
			must.BeEqual(b, out, *in)

			b.Run("Marshal", func(t *testing.B) {
				t.ReportAllocs()
				for i := 0; i < t.N; i++ {
					codec.Encoding(*in)
				}
			})

			b.Run("Unmarshal", func(t *testing.B) {
				t.ReportAllocs()
				for i := 0; i < t.N; i++ {
					codec.Decoding(bytes, &out)
				}
			})
		})
	}
}

func TestClosedDB(t *testing.T) {
	dbName := "TestClosedDB"

	for _, codec := range benchmark.CodecList {
		t.Run(codec.Name, func(t *testing.T) {
			fmt.Println("####", t.Name())

			// open db
			db, err := badgerx.NewDB(dbName)
			must.BeNoError(t, err)
			must.NotBeNil(t, db)
			tx := db.Tx()

			// insert item
			in := benchmark.NewTestStruct()
			in.EncodeFunc = codec.Encoding

			err = tx.InsertIDStruct(in)
			must.BeNoError(t, err)

			// close db
			must.BeNoError(t, tx.Close(true))
			must.BeNoError(t, db.Close())

			// open again
			db, err = badgerx.NewDB(dbName)
			must.BeNoError(t, err)
			must.NotBeNil(t, db)
			tx = db.Tx()

			// select iten
			var out benchmark.TestStructWrapper
			out.DecodeFunc = codec.Decoding

			found, err := tx.GetIDStruct(in.ID(), &out)
			must.BeNoError(t, err)
			must.BeTrue(t, found)
			must.BeEqual(t, in.TestStruct, out.TestStruct)

			// close db
			must.BeNoError(t, tx.Close(false))
			must.BeNoError(t, db.Close())
		})
	}
}
