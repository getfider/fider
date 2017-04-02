package dbx

import (
	"database/sql"
	"io/ioutil"

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

	db := &Database{conn}
	if env.IsTest() {
		db.load("/app/dbx/setup.sql")
	}
	return db, nil
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

// Exists returns true if at least one record is found
func (db Database) Exists(command string, args ...interface{}) (bool, error) {
	rows, err := db.conn.Query(command, args...)
	if rows != nil {
		rows.Close()
		return true, nil
	}
	return false, err
}

// Count returns number of rows
func (db Database) Count(command string, args ...interface{}) (int, error) {
	rows, err := db.conn.Query(command, args...)
	defer rows.Close()
	count := 0
	for rows != nil && rows.Next() {
		count++
	}
	return count, err
}

// Begin returns a new SQL transaction
func (db Database) Begin() (*sql.Tx, error) {
	return db.conn.Begin()
}

// Close connection to database
func (db Database) Close() {
	if db.conn != nil {
		throw(db.conn.Close())
	}
}

func (db Database) load(path string) {
	content, err := ioutil.ReadFile(env.Path(path))
	if err != nil {
		panic(err)
	}

	err = db.Execute(string(content))
	if err != nil {
		panic(err)
	}
}

func throw(err error) error {
	if err != nil && env.IsTest() {
		panic(err)
	}
	return err
}
