package dbx

import (
	"database/sql"

	"github.com/WeCanHearYou/wechy/app/toolbox/env"

	//required
	_ "github.com/lib/pq"
)

// New creates a new Database instance
func New() (*Database, error) {
	conn, err := sql.Open("postgres", env.MustGet("DATABASE_URL"))
	if err != nil {
		return nil, throw(err)
	}
	return &Database{conn}, nil
}

// Database represents a connection to a SQL database
type Database struct {
	conn *sql.DB
}

// Execute given SQL command
func (db Database) Execute(command string, args ...interface{}) error {
	_, err := db.conn.Exec(command, args...)
	if err != nil {
		return throw(err)
	}
	return nil
}

// Query the database with given SQL command
func (db Database) Query(command string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(command, args...)
}

// QueryRow the database with given SQL command and returns only one row
func (db Database) QueryRow(command string, args ...interface{}) *sql.Row {
	return db.conn.QueryRow(command, args...)
}

// Begin returns a new SQL transaction
func (db Database) Begin() (*sql.Tx, error) {
	return db.conn.Begin()
}

// Close connection to database
func (db Database) Close() {
	if env.IsTest() {
		db.Execute("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
		db.Execute("TRUNCATE TABLE tenants RESTART IDENTITY CASCADE")
	}

	if db.conn != nil {
		db.conn.Close()
	}
}

func throw(err error) error {
	if env.IsTest() {
		panic(err)
	}
	return err
}
