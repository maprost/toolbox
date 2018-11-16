package benchmark

import (
	"fmt"
	"testing"
	"time"

	"github.com/maprost/testbox/must"
	"github.com/maprost/toolbox/pqx"
)

const (
	testEntity_table   = "testEntity"
	testEntity_id      = "id"
	testEntity_message = "msg"
	testEntity_changed = "changed"
	testEntity_created = "created"
)

var testEntityConfig = &pqx.EntityConfig{
	PKName:        testEntity_id,
	Autoincrement: true,
	ChangedName:   testEntity_changed,
	CreatedName:   testEntity_created,
	Table:         testEntity_table,
	Columns:       []string{testEntity_id, testEntity_message, testEntity_changed, testEntity_created},
}

type TestEntity struct {
	ID      int64
	Message string
	Changed time.Time
	Created time.Time
}

func (t TestEntity) EntityConfig() *pqx.EntityConfig {
	return testEntityConfig
}

func (t *TestEntity) PtrValue(column string) (interface{}, error) {
	switch column {
	case testEntity_id:
		return &t.ID, nil
	case testEntity_message:
		return &t.Message, nil
	case testEntity_changed:
		return &t.Changed, nil
	case testEntity_created:
		return &t.Created, nil
	}

	return nil, fmt.Errorf("can't find column %s", column)
}

func CreateTestEntityTableIfNotExists(t testing.TB, queryfunc func(sql pqx.SQL) pqx.RowsResult) {
	// create table
	sql := pqx.NewSQL()
	sql.Writef(`
		CREATE TABLE IF NOT EXISTS %s(
			%s bigserial,
			%s text,
			%s timestamp with time zone,
			%s timestamp with time zone
		);`, testEntity_table, testEntity_id, testEntity_message, testEntity_created, testEntity_changed)

	rowResult := queryfunc(sql)
	must.BeNoError(t, rowResult.Err)
}
