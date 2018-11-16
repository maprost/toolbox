package pqx

// FillEntity fill the entity (pk has to be set)
// SELECT column1, column2,... FROM [Table] WHERE [PK] = [PK_value]
func FillEntity(entity Entity) (bool, Result) {
	return fillEntityFunc(singleDB, entity)
}

// FillEntity fill the entity (pk has to be set)
// SELECT column1, column2,... FROM [Table] WHERE [PK] = [PK_value]
func (db *DB) FillEntity(entity Entity) (bool, Result) {
	return fillEntityFunc(db.db, entity)
}

// FillEntity fill the entity (pk has to be set)
// SELECT column1, column2,... FROM [Table] WHERE [PK] = [PK_value]
func (tx *Transaction) FillEntity(entity Entity) (bool, Result) {
	return fillEntityFunc(tx.tx, entity)
}

// SELECT column1, column2,... FROM [Table] WHERE [PK] = [PK_value]
func fillEntityFunc(exe originExecutor, entity Entity) (bool, Result) {
	config := entity.EntityConfig()

	pkValue, err := getPrtValue(entity, config.PKName)
	if err != nil {
		return false, Result{Err: err}
	}

	// preparation of the statement
	sql := NewSQL()
	sql.Writef("SELECT %s FROM %s WHERE %s=%a;", config.ColumnList(), config.Table, config.PKName, pkValue)

	// execute statement
	rowsResult := query(exe, sql)
	defer rowsResult.Close()

	found, err := rowsResult.ScanEntity(entity)
	if err != nil {
		return false, resultError(rowsResult.Result, err)
	}

	return found, rowsResult.Result
}
