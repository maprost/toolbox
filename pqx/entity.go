package pqx

import (
	"errors"
	"reflect"
	"strings"
)

type Entity interface {
	EntityConfig() *EntityConfig
	PtrValue(column string) (interface{}, error)
}

type EntityConfig struct {
	Table         string   // !
	Columns       []string // !
	PKName        string   // (optional)
	Autoincrement bool     // (optional) is the PK an autoincrement unit? if so: the PK type should be an int8/uint8/int32/uint32/int/int64, uint64 can't be used
	ChangedName   string   // (optional) name of a time unit, that will set every insert/update
	CreatedName   string   // (optional) name of a time unit, that will set every insert
}

func (c EntityConfig) ColumnList() string {
	return strings.Join(c.Columns, ",")
}

func setPrtValue(entity Entity, column string, toInsert interface{}) error {
	ptrValue, err := entity.PtrValue(column)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(ptrValue)
	if val.Kind() != reflect.Ptr {
		return errors.New("some: check must be a pointer")
	}
	val.Elem().Set(reflect.ValueOf(toInsert))

	return nil
}

func getPrtValue(entity Entity, column string) (interface{}, error) {
	ptrValue, err := entity.PtrValue(column)
	if err != nil {
		return nil, err
	}

	val := reflect.ValueOf(ptrValue)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("some: check must be a pointer")
	}

	return val.Elem().Interface(), nil
}
