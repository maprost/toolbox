package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/maprost/toolbox/badgerx"
	"github.com/maprost/toolbox/badgerx/benchmark"
	"github.com/maprost/toolbox/util/tbstring"
)

const (
	//dbName  = "BenchmarkDatabaseSizeUnsafe3"
	dbName     = "BenchmarkDatabaseSizeProto2"
	inserts    = 10000
	goRoutines = 10
)

var (
	Encoder = benchmark.ProtoEncoder
	Decoder = benchmark.ProtoDecoder
)

func main() {
	rand.Seed(time.Now().UnixNano())

	maxID := fillDB()
	checkDB(maxID)
}

type TestStructWrapper struct {
	benchmark.TestStructWrapper
}

func (ts *TestStructWrapper) SetID(id uint64) {
	ts.TestStruct.ID = id
	ts.TestStruct.Msg = []byte(strconv.FormatUint(id, 10) + tbstring.RandomString(10))
	//ts.Liste = []int{int(id)}
	ts.Flag = true

	//fmt.Println("##", unsafe.Sizeof(ts.TestStruct))
}

func (ts *TestStructWrapper) Check() bool {
	return strings.HasPrefix(string(ts.Msg), strconv.FormatUint(ts.ID(), 10))
}

func fillDB() uint64 {
	db, err := badgerx.NewDB(dbName)
	if err != nil {
		panic(err)
	}

	idChan := make(chan uint64)

	for g := 0; g < goRoutines; g++ {
		go func() {
			var maxID uint64
			tx := db.Tx()

			for i := 1; i <= inserts; i++ {
				var testStruct TestStructWrapper
				testStruct.EncodeFunc = Encoder

				err := tx.InsertIDStruct(&testStruct)
				if err != nil {
					panic(err)
				}

				incCounter(inserts * goRoutines)
				maxID = testStruct.ID()

				// refresh connection
				if i%100 == 0 {
					err = tx.Close(true)
					if err != nil {
						panic(err)
					}
					tx = db.Tx()
				}
			}

			err = tx.Close(true)
			if err != nil {
				panic(err)
			}

			idChan <- maxID
		}()
	}

	maxID := uint64(0)
	for g := 0; g < goRoutines; g++ {
		routineID := <-idChan
		if maxID < routineID {
			maxID = routineID
		}
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}

	return maxID
}

func checkDB(maxID uint64) {
	db, err := badgerx.NewDB(dbName)
	if err != nil {
		panic(err)
	}
	tx := db.Tx()

	foundCounter := 0
	correctCounter := 0
	t0 := time.Now()
	for i := uint64(0); i < maxID; i++ {
		fmt.Print(".")
		if i%50 == 0 {
			fmt.Println(" ", i, "/", maxID)
		}

		var toCheck TestStructWrapper
		toCheck.DecodeFunc = Decoder

		found, err := tx.GetIDStruct(i, &toCheck)
		if err != nil {
			panic(err)
		}
		if found {
			foundCounter++
		}

		//fmt.Print("Check: actual(", toCheck.Msg, ") expected(a)")
		if toCheck.Check() {
			correctCounter++
		} else {
			fmt.Printf("Incorrect! %+v\n", toCheck.TestStruct)
		}
		//fmt.Println()
	}

	err = tx.Close(true)
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}

	duration := time.Since(t0)
	fmt.Println("\nN:", maxID, " Found:", foundCounter, " Correct:", correctCounter, " Duration:", duration.Seconds(), "s - ", uint64(duration.Nanoseconds())/maxID, "ns/op")
}

var (
	counter      = 0
	counterMutex = sync.Mutex{}
)

func incCounter(max int) {
	counterMutex.Lock()
	counter++

	fmt.Print(".")
	if counter%100 == 0 {
		fmt.Println(" ", counter, "/", max)
	}

	counterMutex.Unlock()
}

/*

Encode:  {0  a 0 [0] true}  -->  [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 56 225 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
.Encode:  {1  a 0 [1] true}  -->  [1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 8 226 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
.Encode:  {2  a 0 [2] true}  -->  [2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 216 226 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
.Encode:  {3  a 0 [3] true}  -->  [3 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 152 227 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
.Encode:  {4  a 0 [4] true}  -->  [4 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 88 228 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
.Encode:  {5  a 0 [5] true}  -->  [5 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 24 229 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
.Encode:  {6  a 0 [6] true}  -->  [6 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 216 229 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
.Encode:  {7  a 0 [7] true}  -->  [7 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 152 230 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
.Encode:  {8  a 0 [8] true}  -->  [8 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 88 231 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
.Encode:  {9  a 0 [9] true}  -->  [9 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 116 35 165 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 24 232 22 32 196 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0]
*/
