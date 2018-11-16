package pqx

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type ConnectionInfo interface {
	Database() string
	Host() string
	Port() string
	Username() string
}

type ConnectionConfig struct {
	database string
	host     string
	port     string
	username string
}

func (c ConnectionConfig) Database() string {
	return c.database
}
func (c ConnectionConfig) Host() string {
	return c.host
}
func (c ConnectionConfig) Port() string {
	return c.port
}
func (c ConnectionConfig) Username() string {
	return c.username
}

func NewDefaultConnectionInfo(db string) ConnectionInfo {
	return ConnectionConfig{
		database: db,
		port:     "5432",
		host:     "localhost",
		username: "postgres",
	}
}

var singleDB *sql.DB = nil

type DB struct {
	db *sql.DB
}

func OpenDatabaseConnection(info ConnectionInfo) (DB, error) {
	db, err := sql.Open("postgres",
		"user="+info.Username()+
			" dbname="+info.Database()+
			" sslmode=disable"+
			" host="+info.Host()+
			" port="+info.Port())

	if err != nil {
		return DB{}, err
	}

	err = db.Ping() // check db connection
	return DB{db: db}, err
}

func OpenSingleDatabaseConnection(info ConnectionInfo) error {
	// Close old db connection before open a new one
	if singleDB != nil {
		err := singleDB.Close()
		if err != nil {
			return err
		}
	}

	db, err := OpenDatabaseConnection(info)
	singleDB = db.db

	return err
}

func Query(sql SQL) RowsResult {
	return query(singleDB, sql)
}
