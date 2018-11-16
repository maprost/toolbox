package pqx_test

import (
	"github.com/maprost/toolbox/pqx"
	"testing"

	"github.com/maprost/testbox/should"
)

func TestInsertSQL(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := []interface{}{1, "hello", false}
	pkName := "a"
	table := "blob"

	t.Run("check insert statement", func(t *testing.T) {
		// build sql
		sql := pqx.NewSQL()
		sql.Writef("Insert INTO %s (", table)
		for _, k := range keys {
			sql.Listf("%s", k)
		}
		sql.Writef(") VALUES (")
		for i, v := range values {
			if keys[i] == pkName {
				sql.Listf("Default")
			} else {
				sql.Listf("%a", v)
			}
		}
		sql.Writef(") RETURNING %s;", pkName)

		// check content
		should.BeEqual(t, sql.String(), "Insert INTO blob ( a ,b ,c ) VALUES ( Default ,$1 ,$2 ) RETURNING a; [hello false]")
	})

	t.Run("check contains statement", func(t *testing.T) {
		// build sql
		sql := pqx.NewSQL()
		sql.Writef("Select %s", pkName)
		sql.Writef("FROM %s", table)
		sql.Writef("WHERE %s=%a;", pkName, values[0])

		// check content
		should.BeEqual(t, sql.String(), "Select a FROM blob WHERE a=$1; [1]")

	})

}
