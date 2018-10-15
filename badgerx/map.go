package badgerx

type DBMap map[string]*DB // path -> db
type TxMap map[string]*Tx // path -> tx

func NewDBMap(paths []string) (DBMap, error) {
	rootMap := make(DBMap)

	for _, path := range paths {
		db, err := NewDB(path)
		if err != nil {
			return nil, err
		}

		rootMap[path] = &db
	}

	return rootMap, nil
}

func (m DBMap) TxMap() TxMap {
	txMap := make(TxMap)

	for path, db := range m {
		tx := db.Tx()
		txMap[path] = &tx
	}

	return txMap
}

func (m DBMap) Close() error {
	for _, db := range m {
		if err := db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m TxMap) Close(commit bool) error {
	for _, tx := range m {
		if err := tx.Close(commit); err != nil {
			return err
		}
	}
	return nil
}
