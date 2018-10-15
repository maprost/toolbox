// generate easyjson str.go
package benchmark

import (
	"encoding/json"
	"math"
	"unsafe"

	"github.com/json-iterator/go"
	"github.com/ugorji/go/codec"
)

type EncodeFunc func(b TestStructWrapper) ([]byte, error)
type DecodeFunc func(in []byte, out *TestStructWrapper) error

type Codec struct {
	Name     string
	Encoding EncodeFunc
	Decoding DecodeFunc
}

var (
	DefaultEncoder = func(b TestStructWrapper) ([]byte, error) { return json.Marshal(b.TestStruct) }
	DefaultDecoder = func(in []byte, out *TestStructWrapper) error { return json.Unmarshal(in, &out.TestStruct) }

	UnsafePointerEncoder = func(b TestStructWrapper) ([]byte, error) {
		size := unsafe.Sizeof(b.TestStruct)
		bytes := GetBytesStoredAt(unsafe.Pointer(&b.TestStruct), size)
		//fmt.Println("Encode: ", b.TestStruct, " --> ", bytes)
		return bytes, nil
	}
	UnsafePointerDecoder = func(in []byte, out *TestStructWrapper) error {
		testStruct := (*TestStruct)(unsafe.Pointer(&in[0]))
		out.TestStruct = *testStruct
		//fmt.Println("Decode: ", in, " --> ", out.TestStruct)
		return nil
	}

	EasyJsonEncoder = func(b TestStructWrapper) ([]byte, error) { return b.TestStruct.MarshalJSON() }
	EasyJsonDecoder = func(in []byte, out *TestStructWrapper) error { return out.TestStruct.UnmarshalJSON(in) }

	MsgPackEncoder = func(b TestStructWrapper) (bytes []byte, err error) {
		var msgPack codec.MsgpackHandle
		enc := codec.NewEncoderBytes(&bytes, &msgPack)
		err = enc.Encode(b)
		return
	}
	MsgPackDecoder = func(in []byte, out *TestStructWrapper) error {
		var msgPack codec.MsgpackHandle
		dec := codec.NewDecoderBytes(in, &msgPack)
		return dec.Decode(out)
	}

	JsoniterEncoder = func(b TestStructWrapper) ([]byte, error) { return jsoniter.Marshal(b.TestStruct) }
	JsoniterDecoder = func(in []byte, out *TestStructWrapper) error { return jsoniter.Unmarshal(in, &out.TestStruct) }

	ProtoEncoder = func(b TestStructWrapper) ([]byte, error) {
		proto := TestStruct2{
			Id:      b.ID(),
			Key:     b.Key(),
			Msg:     string(b.Msg),
			Counter: int32(b.Counter),
			Flag:    b.Flag,
		}
		return proto.Marshal()
	}
	ProtoDecoder = func(in []byte, out *TestStructWrapper) error {
		var proto TestStruct2
		err := proto.Unmarshal(in)
		if err != nil {
			return err
		}
		out.TestStruct = TestStruct{
			ID:      proto.Id,
			Key:     []byte(proto.Key),
			Msg:     []byte(proto.Msg),
			Counter: int(proto.Counter),
			Flag:    proto.Flag,
		}
		return nil
	}
)

var CodecList = []Codec{
	{"default", DefaultEncoder, DefaultDecoder},
	{"easyJson", EasyJsonEncoder, EasyJsonDecoder},
	{"Jsoniter", JsoniterEncoder, JsoniterDecoder},
	{"msgPack", MsgPackEncoder, MsgPackDecoder},
	{"unsafePointer", UnsafePointerEncoder, UnsafePointerDecoder},
	{"proto", ProtoEncoder, ProtoDecoder},
}

//easyjson:json
type TestStruct struct {
	ID      uint64 `json:"id"`
	Key     []byte `json:"key"`
	Msg     []byte `json:"msg"`
	Counter int    `json:"counter"`
	//Liste   []int  `json:"liste"`
	Flag bool `json:"flag"`
}

type TestStructWrapper struct {
	TestStruct

	EncodeFunc EncodeFunc `json:"-"`
	DecodeFunc DecodeFunc `json:"-"`
}

func (ts TestStructWrapper) Key() string {
	return string(ts.TestStruct.Key)
}

func (ts TestStructWrapper) ID() uint64 {
	return ts.TestStruct.ID
}

func (ts *TestStructWrapper) SetID(id uint64) {
	ts.TestStruct.ID = id
}

func (ts TestStructWrapper) Encode() (b []byte, err error) {
	if ts.EncodeFunc == nil {
		b, err = DefaultEncoder(ts)
	} else {
		b, err = ts.EncodeFunc(ts)
	}

	//if err == nil {
	//	fmt.Println("Encode: ", string(b), "  ", ts.ID(), " '", ts.Key(), "'")
	//}

	return b, err
}

func (ts *TestStructWrapper) Decode(bytes []byte) (err error) {
	if ts.DecodeFunc != nil {
		err = ts.DecodeFunc(bytes, ts)
	} else {
		err = DefaultDecoder(bytes, ts)
	}

	//if err == nil {
	//	fmt.Printf("Decode: %+v\n", ts)
	//}

	return err
}

func NewTestStruct() *TestStructWrapper {
	return &TestStructWrapper{
		TestStruct: TestStruct{
			Msg:     []byte(".."),
			Counter: 1,
			//Liste:   []int{1, 2, 3, 4, 5, 6, math.MaxInt32},
			Flag: true,
		},
	}
}

func NewKeyTestStruct(key string) *TestStructWrapper {
	b := NewTestStruct()
	b.TestStruct.Key = []byte(key)

	return b
}

// getBytesStoredAt converts a pointer to the raw bytes go stores in memory at
// the pointer's position.
// This function will only work reliably for pointers to structs that only contain
// value types
func GetBytesStoredAt(ptr unsafe.Pointer, size uintptr) []byte {
	// get a byte array of the memory stored at the pointer's position
	pointerToRawBytes := (*[math.MaxInt32]byte)(ptr)
	return (*pointerToRawBytes)[:size]
}
