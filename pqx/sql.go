package pqx

import (
	"fmt"
	"strconv"
	"strings"
)

type SQL struct {
	sql      string
	counter  int
	args     []interface{}
	listMode bool
}

func NewSQL() SQL {
	return SQL{counter: 1}
}

func (s *SQL) Writef(sql string, args ...interface{}) {
	s.noList()
	s.write(sql, args)
}

func (s *SQL) Listf(sql string, args ...interface{}) {
	sql = s.list(sql)
	s.write(sql, args)
}

func (s *SQL) write(sql string, args []interface{}) {
	sql = checkLastSpace(sql)

	// convert args
	fArgs := make([]interface{}, 0, len(args))
	argsIndex := 0
	for i := 0; i < (len(sql) - 1); i++ {
		if sql[i] == '%' {

			if sql[i+1] == 'a' {
				sqlNumber := s.next(args[argsIndex])
				sql = strings.Replace(sql, "%a", sqlNumber, 1)

			} else {
				fArgs = append(fArgs, args[argsIndex])
			}

			argsIndex++
		}
	}

	// insert args counter
	s.sql += fmt.Sprintf(sql, fArgs...)
}

func (s *SQL) next(arg interface{}) string {
	result := "$" + strconv.Itoa(s.counter)
	s.counter++
	s.args = append(s.args, arg)
	return result
}

func (s *SQL) list(sql string) string {
	if s.listMode {
		sql = "," + sql
	}

	s.listMode = true
	return sql
}

func (s *SQL) noList() {
	s.listMode = false
}

func (s SQL) String() string {
	return fmt.Sprint(s.sql, s.args)
}

func checkLastSpace(sql string) string {
	if !strings.HasSuffix(sql, " ") {
		sql += " "
	}
	return sql
}
