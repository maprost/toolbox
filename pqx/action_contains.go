package pqx

// Contains returns true if the entity exists.
// SELECT [PK] FROM [Table] WHERE [PK] = [PK_value]
func Contains(entity Entity) (bool, Result) {
	return containsFunc(singleDB, entity)
}

// Contains returns true if the entity exists.
// SELECT [PK] FROM [Table] WHERE [PK] = [PK_value]
func (db *DB) Contains(entity Entity) (bool, Result) {
	return containsFunc(db.db, entity)
}

// Contains returns true if the entity exists.
// SELECT [PK] FROM [Table] WHERE [PK] = [PK_value]
func (tx *Transaction) Contains(entity Entity) (bool, Result) {
	return containsFunc(tx.tx, entity)
}

// SELECT [PK] FROM [Table] WHERE [PK] = [PK_value]
func containsFunc(exe originExecutor, entity Entity) (bool, Result) {
	config := entity.EntityConfig()

	value, err := getPrtValue(entity, config.PKName)
	if err != nil {
		return false, Result{Err: err}
	}

	sql := NewSQL()
	sql.Writef("Select %s FROM %s WHERE %s=%a;", config.PKName, config.Table, config.PKName, value)
	rowsResult := query(exe, sql)
	defer rowsResult.Close()

	return rowsResult.Next(), rowsResult.Result
}
