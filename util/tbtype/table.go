package tbtype

import "sync"

// Table can be used as in-memory database
type Table struct {
	data  map[uint64]interface{}
	mutex *sync.RWMutex
}

// NewTable creates a new initialized 'Table'
func NewTable() Table {
	return Table{
		data:  make(map[uint64]interface{}),
		mutex: &sync.RWMutex{},
	}
}

// Insert a new value into the table, the function get a new key and returns the value to insert
func (t *Table) Insert(valueFunc func(key uint64) interface{}) {
	// to avoid duplicate key, need a write look
	t.mutex.Lock()

	// search for ID
	key := uint64(1)
	for k := range t.data {
		if uint64(k) >= key {
			key = uint64(k) + 1
		}
	}

	// store
	t.data[key] = valueFunc(key)
	t.mutex.Unlock()
}

// Update deletes the value inside the 'Table' and insert it afterward
func (t *Table) Update(key uint64, value interface{}) {
	t.Delete(key)

	t.mutex.Lock()
	t.data[key] = value
	t.mutex.Unlock()
}

// Delete deletes the value of the given key
func (t *Table) Delete(key uint64) {
	t.mutex.Lock()
	delete(t.data, key)
	t.mutex.Unlock()
}

// Get returns the value of the given key
func (t *Table) Get(key uint64) (interface{}, bool) {
	t.mutex.RLock()
	v, ok := t.data[key]
	t.mutex.RUnlock()

	return v, ok
}

// Exists checks if the 'Table' contains the key
func (t *Table) Exists(key uint64) bool {
	t.mutex.RLock()
	_, ok := t.data[key]
	t.mutex.RUnlock()

	return ok
}

// FindFirst returns the first element that passes the 'checkFunc'
func (t *Table) FindFirst(checkFunc func(value interface{}) bool) (interface{}, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for _, value := range t.data {
		if checkFunc(value) {
			return value, true
		}
	}

	return nil, false
}

// FindAll returns all ids that passes the 'checkFunc'
func (t *Table) FindAll(checkFunc func(value interface{}) bool) []uint64 {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	ids := make([]uint64, 0, 2)
	for key, value := range t.data {
		if checkFunc(value) {
			ids = append(ids, key)
		}
	}

	return ids
}
