package pqx

import (
	"errors"
	"time"
)

// INSERT INTO table_name (AI, column1,column2,column3,...)
// VALUES (DEFAULT, value1,value2,value3,...) RETURNING AI;
func Insert(entity Entity) Result {
	return insertFunc(singleDB, entity)
}

// INSERT INTO table_name (AI, column1,column2,column3,...)
// VALUES (DEFAULT, value1,value2,value3,...) RETURNING AI;
func (db *DB) Insert(entity Entity) Result {
	return insertFunc(db.db, entity)
}

// INSERT INTO table_name (AI, column1,column2,column3,...)
// VALUES (DEFAULT, value1,value2,value3,...) RETURNING AI;
func (tx *Transaction) Insert(entity Entity) Result {
	return insertFunc(tx.tx, entity)
}

// INSERT INTO table_name (AI, column1,column2,column3,...)
// VALUES (DEFAULT, value1,value2,value3,...) RETURNING AI
func insertFunc(exe originExecutor, entity Entity) Result {
	config := entity.EntityConfig()

	// refresh time
	if config.ChangedName != "" {
		err := setPrtValue(entity, config.ChangedName, time.Now())
		if err != nil {
			return Result{Err: err}
		}
	}

	if config.CreatedName != "" {
		err := setPrtValue(entity, config.CreatedName, time.Now())
		if err != nil {
			return Result{Err: err}
		}
	}

	// preparation of the statement
	sql := NewSQL()
	sql.Writef("INSERT INTO %s (%s) VALUES (", config.Table, config.ColumnList())

	for _, column := range config.Columns {
		if column == config.PKName && config.Autoincrement {
			sql.Listf("DEFAULT")

		} else {
			value, err := getPrtValue(entity, column)
			if err != nil {
				return Result{Err: err}
			}

			sql.Listf("%a", value)
		}
	}

	sql.Writef(")")
	if config.Autoincrement {
		sql.Writef(" RETURNING %s", config.PKName)
	}
	sql.Writef(";")

	// execute statement
	result := query(exe, sql)
	defer result.Close()
	if result.Err != nil {
		return result.Result
	}

	// update pk with returning value (if needed)
	if config.Autoincrement {
		if result.Next() == false {
			return resultError(result.Result, errors.New("no return element in insert statement"))
		}

		id, err := entity.PtrValue(config.PKName)
		if err != nil {
			return resultError(result.Result, err)
		}

		err = result.Scan(id)
		if err != nil {
			return resultError(result.Result, err)
		}
	}

	return result.Result
}
