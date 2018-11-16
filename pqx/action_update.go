package pqx

import (
	"time"
)

// Update an entity
// UPDATE [Table] SET column1 = value1, column2 = value2, ... WHERE [PK] = [PK_value]
func Update(entity Entity) Result {
	return updateFunc(singleDB, entity)
}

// Update an entity
// UPDATE [Table] SET column1 = value1, column2 = value2, ... WHERE [PK] = [PK_value]
func (db *DB) Update(entity Entity) Result {
	return updateFunc(db.db, entity)
}

// Update an entity
// UPDATE [Table] SET column1 = value1, column2 = value2, ... WHERE [PK] = [PK_value]
func (tx *Transaction) Update(entity Entity) Result {
	return updateFunc(tx.tx, entity)
}

// UPDATE [Table] SET column1 = value1, column2 = value2, ... WHERE [PK] = [PK_value]
func updateFunc(exe originExecutor, entity Entity) Result {
	config := entity.EntityConfig()

	// refresh time
	if config.ChangedName != "" {
		err := setPrtValue(entity, config.ChangedName, time.Now())
		if err != nil {
			return Result{Err: err}
		}
	}

	sql := NewSQL()
	sql.Writef("UPDATE %s SET", config.Table)

	for _, column := range config.Columns {
		if column == config.CreatedName || column == config.PKName {
			// don't change the create date or pk
			continue
		}

		value, err := getPrtValue(entity, column)
		if err != nil {
			return Result{Err: err}
		}

		sql.Listf("%s=%a", column, value)
	}

	pkValue, err := getPrtValue(entity, config.PKName)
	if err != nil {
		return Result{Err: err}
	}
	sql.Writef(" WHERE %s=%a", config.PKName, pkValue)

	// execute statement
	rowsResult := query(exe, sql)
	rowsResult.Close()
	return rowsResult.Result
}
