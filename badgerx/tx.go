package badgerx

import (
	"encoding/binary"
	"time"

	"github.com/dgraph-io/badger"
)

type ID = uint64
type Key = string

type Encoder interface {
	Encode() ([]byte, error)
}

type KeyEncoder interface {
	Encoder
	Key() Key
}

type IDEncoder interface {
	Encoder
	SetID(ID)
	ID() ID
}

type Decoder interface {
	Decode([]byte) error
}

type KeyIDTx interface {
	UpsertKeyID(key Key, id ID) error
	DeleteKey(key Key) error
	GetKeyID(key Key) (ID, bool, error)
}

type KeyIDsTx interface {
	UpsertKeyIDs(key Key, id []ID) error
	DeleteKey(key Key) error
	GetKeyIDs(key Key) ([]ID, bool, error)
}

type KeyStructTx interface {
	UpsertKeyStruct(val KeyEncoder) error
	DeleteKey(key Key) error
	GetKeyStruct(key Key, val Decoder) (bool, error)
}

type IDStructTx interface {
	InsertIDStruct(val IDEncoder) error
	UpdateIDStruct(val IDEncoder) error
	DeleteID(id ID) error
	GetIDStruct(id ID, val Decoder) (bool, error)
}

type Tx struct {
	// borrowed from base unit (will be closed there)
	db DB

	// created on demand!
	readableTx  *badger.Txn
	writeableTx *badger.Txn

	ttl time.Duration
}

func (t Tx) Close(commit bool) error {
	if t.readableTx != nil {
		t.readableTx.Discard()
	}

	if t.writeableTx != nil {
		if commit {
			err := t.writeableTx.Commit(nil)
			if err != nil {
				return err
			}
		} else {
			t.writeableTx.Discard()
		}
	}
	return nil
}

func (t *Tx) getReadableTx() *badger.Txn {
	if t.writeableTx != nil {
		return t.writeableTx
	}
	if t.readableTx == nil {
		t.readableTx = t.db.newReadableTransaction()
	}
	return t.readableTx
}

func (t *Tx) getWriteableTx() *badger.Txn {
	if t.writeableTx == nil {
		t.writeableTx = t.db.newWriteableTransaction()

		// close read tx
		if t.readableTx != nil {
			t.readableTx.Discard()
			t.readableTx = nil
		}
	}

	return t.writeableTx
}

func (t *Tx) SetTTL(ttl time.Duration) {
	t.ttl = ttl
}

// ---------------------------------------------------------------------
// ----------------------- Upsert --------------------------------------
// ---------------------------------------------------------------------

func (t *Tx) UpsertKeyStruct(val KeyEncoder) error {
	return t.upsertStruct([]byte(val.Key()), val)
}

func (t *Tx) UpsertKeyIDs(key Key, ids []ID) error {
	return t.upsert([]byte(key), idsToBytes(ids))
}

func (t *Tx) InsertIDStruct(val IDEncoder) error {
	id, err := t.addID(val)
	if err != nil {
		return err
	}

	return t.upsertStruct(idToBytes(id), val)
}

func (t *Tx) UpdateIDStruct(val IDEncoder) error {
	return t.upsertStruct(idToBytes(val.ID()), val)
}

func (t *Tx) UpsertKeyID(key Key, id ID) error {
	return t.upsert([]byte(key), idToBytes(id))
}

func (t *Tx) addID(val IDEncoder) (ID, error) {
	id, err := t.db.nextID()
	if err != nil {
		return 0, err
	}

	val.SetID(id)
	return id, nil
}

func (t *Tx) upsertStruct(key []byte, val Encoder) error {
	jVal, err := val.Encode()
	if err != nil {
		return err
	}

	return t.upsert(key, jVal)
}

func (t *Tx) upsert(key []byte, val []byte) error {
	if t.ttl > 0 {
		return t.getWriteableTx().SetWithTTL(key, val, t.ttl)
	}
	return t.getWriteableTx().Set(key, val)
}

// ---------------------------------------------------------------------
// ----------------------- Delete --------------------------------------
// ---------------------------------------------------------------------

func (t *Tx) DeleteKey(key Key) error {
	return t.delete([]byte(key))
}

func (t *Tx) DeleteID(id ID) error {
	return t.delete(idToBytes(id))
}

func (t *Tx) delete(key []byte) error {
	return t.getWriteableTx().Delete(key)
}

// ---------------------------------------------------------------------
// ----------------------- GET -----------------------------------------
// ---------------------------------------------------------------------

func (t *Tx) GetKeyStruct(key Key, val Decoder) (bool, error) {
	return t.getStruct([]byte(key), val)
}

func (t *Tx) GetIDStruct(id ID, val Decoder) (bool, error) {
	return t.getStruct(idToBytes(id), val)
}

func (t *Tx) GetKeyID(key Key) (ID, bool, error) {
	idBytes, found, err := t.getBytes([]byte(key))
	if err != nil {
		return 0, false, err
	}
	if !found {
		return 0, false, nil
	}

	return bytesToID(idBytes), true, nil
}

func (t *Tx) GetKeyIDs(key Key) ([]ID, bool, error) {
	idsBytes, found, err := t.getBytes([]byte(key))
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	return bytesToIDs(idsBytes), true, nil
}

func (t *Tx) getStruct(key []byte, val Decoder) (bool, error) {
	valBytes, found, err := t.getBytes(key)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	err = val.Decode(valBytes)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (t *Tx) getBytes(key []byte) ([]byte, bool, error) {
	item, err := t.getReadableTx().Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	if item.IsDeletedOrExpired() {
		return nil, false, nil
	}

	valBytes, err := item.Value()
	if err != nil {
		return nil, false, err
	}

	return valBytes, true, nil
}

// ---------------------------------------------------------------------
// ----------------------- Helper --------------------------------------
// ---------------------------------------------------------------------

func bytesToID(b []byte) ID {
	return binary.BigEndian.Uint64(b)
}

func idToBytes(id ID) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], id)
	return buf[:]
}

func idsToBytes(ids []ID) []byte {
	byteSize := 8
	sliceSize := byteSize * len(ids)
	buf := make([]byte, sliceSize)

	for i, id := range ids {
		byteIdx := i * byteSize
		binary.BigEndian.PutUint64(buf[byteIdx:byteIdx+byteSize], id)
	}
	return buf[:]
}

func bytesToIDs(byteIds []byte) []ID {
	byteSize := 8
	sliceSize := len(byteIds) / byteSize
	buf := make([]ID, 0, sliceSize)

	for i := 0; i < len(byteIds); i += byteSize {
		buf = append(buf, binary.BigEndian.Uint64(byteIds[i:i+byteSize]))
	}
	return buf
}
