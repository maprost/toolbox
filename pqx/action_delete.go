package pqx

// Delete an entity
// DELETE FROM [Table] WHERE [PK] = [PK_value];
func Delete(entity Entity) Result {
	return deleteFunc(singleDB, entity)
}

// Delete an entity
// DELETE FROM [Table] WHERE [PK] = [PK_value];
func (db *DB) Delete(entity Entity) Result {
	return deleteFunc(db.db, entity)
}

// Delete an entity
// DELETE FROM [Table] WHERE [PK] = [PK_value];
func (tx *Transaction) Delete(entity Entity) Result {
	return deleteFunc(tx.tx, entity)
}

// DELETE FROM [Table] WHERE [PK] = [PK_value];
func deleteFunc(exe originExecutor, entity Entity) Result {
	config := entity.EntityConfig()

	value, err := getPrtValue(entity, config.PKName)
	if err != nil {
		return Result{Err: err}
	}

	sql := NewSQL()
	sql.Writef("DELETE FROM %s WHERE %s=%a;", config.Table, config.PKName, value)
	rowsResult := query(exe, sql)
	rowsResult.Close()

	return rowsResult.Result
}
