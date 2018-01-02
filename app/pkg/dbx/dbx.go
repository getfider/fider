package dbx

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/lib/pq"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/mattes/migrate"

	//required
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

// New creates a new Database instance
func New() (*Database, error) {
	return NewWithLogger(log.NewConsoleLogger())
}

// NewWithLogger creates a new Database instance with logging
func NewWithLogger(logger log.Logger) (*Database, error) {
	conn, err := sql.Open("postgres", env.MustGet("DATABASE_URL"))
	conn.SetMaxIdleConns(20)
	conn.SetMaxOpenConns(50)
	if err != nil {
		return nil, err
	}

	db := &Database{conn, logger, NewRowMapper()}
	return db, nil
}

// Database represents a connection to a SQL database
type Database struct {
	conn   *sql.DB
	logger log.Logger
	mapper RowMapper
}

// Begin returns a new SQL transaction
func (db Database) Begin() (*Trx, error) {
	tx, err := db.conn.Begin()
	return &Trx{tx: tx, logger: db.logger, mapper: db.mapper}, err
}

// Close connection to database
func (db Database) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func scan(prefix string, data interface{}) map[string]interface{} {
	fields := make(map[string]interface{})

	s := reflect.ValueOf(data).Elem()

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		typeField := s.Type().Field(i)
		tag := typeField.Tag.Get("db")

		if tag != "" {
			if typeField.Type.Kind() == reflect.Slice {
				obj := reflect.New(reflect.MakeSlice(typeField.Type, 0, 0).Type()).Elem()
				field.Set(obj)
				fields[prefix+tag] = pq.Array(field.Addr().Interface())
			} else if typeField.Type.Kind() != reflect.Ptr {
				fields[prefix+tag] = field.Addr().Interface()
			} else if field.Type().Elem().Kind() != reflect.Struct || field.Type().Elem().String() == "time.Time" {
				obj := reflect.New(field.Type().Elem()).Elem()
				field.Set(obj.Addr())
				fields[prefix+tag] = field.Interface()
			} else {
				obj := reflect.New(field.Type().Elem()).Elem()
				field.Set(obj.Addr())
				for name, address := range scan(prefix+tag+"_", field.Interface()) {
					fields[name] = address
				}
			}
		}
	}

	return fields
}

func fill(rows *sql.Rows, data interface{}) error {
	columns, _ := rows.Columns()
	fields := scan("", data)

	pointers := make([]interface{}, len(columns))
	for i, column := range columns {
		if pointer, ok := fields[column]; ok {
			pointers[i] = pointer
		} else {
			panic(fmt.Sprintf("No target for column %s", column))
		}
	}

	return rows.Scan(pointers...)
}

func (db Database) load(path string) {
	content, err := ioutil.ReadFile(env.Path(path))
	if err != nil {
		panic(err)
	}

	_, err = db.conn.Exec(string(content))
	if err != nil {
		panic(err)
	}
}

// Seed clean and insert new seed data for testing
func (db Database) Seed() {
	if env.IsTest() {
		db.load("/app/pkg/dbx/setup.sql")
	}
}

// Migrate the database to latest verion
func (db Database) Migrate() {

	db.logger.Infof("Running migrations...")
	m, err := migrate.New(
		"file://"+env.Path("migrations"),
		env.MustGet("DATABASE_URL"),
	)

	if err == nil {
		err = m.Up()
	}

	if err != nil && err != migrate.ErrNoChange {
		db.logger.Infof("Error: %s.", err)

		panic("Migrations failed.")
	} else {
		db.logger.Infof("Migrations finished with success.")
	}
}

//Trx represents a Database transaction
type Trx struct {
	tx     *sql.Tx
	logger log.Logger
	mapper RowMapper
}

// QueryRow the database with given SQL command and returns only one row
func (trx Trx) QueryRow(command string, args ...interface{}) *sql.Row {
	command = formatCommand(command)
	trx.logger.Debugf("%s %v", log.Yellow(command), log.Blue(args))
	return trx.tx.QueryRow(command, args...)
}

// Query the database with given SQL command
func (trx Trx) Query(command string, args ...interface{}) (*sql.Rows, error) {
	command = formatCommand(command)
	trx.logger.Debugf("%s %v", log.Yellow(command), log.Blue(args))
	return trx.tx.Query(command, args...)
}

// Execute given SQL command
func (trx Trx) Execute(command string, args ...interface{}) error {
	command = formatCommand(command)
	trx.logger.Debugf("%s %v", log.Yellow(command), log.Blue(args))
	_, err := trx.tx.Exec(command, args...)
	return err
}

func formatCommand(cmd string) string {
	cmd = strings.Replace(cmd, "\t", "", -1)
	cmd = strings.Replace(cmd, "\n", " ", -1)
	return cmd
}

// Scalar returns first row and first column
func (trx Trx) Scalar(data interface{}, command string, args ...interface{}) error {
	row := trx.QueryRow(command, args...)
	err := row.Scan(data)
	if err == sql.ErrNoRows {
		return app.ErrNotFound
	}
	return err
}

// Get first row and bind to given data
func (trx Trx) Get(data interface{}, command string, args ...interface{}) error {
	rows, err := trx.Query(command, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		columns, _ := rows.Columns()
		return trx.mapper.Map(data, columns, rows.Scan)
	}

	return app.ErrNotFound
}

// QueryIntArray executes given SQL command and return first column as int
func (trx Trx) QueryIntArray(command string, args ...interface{}) ([]int, error) {
	values := make([]int, 0)
	var value int

	rows, err := trx.Query(command, args...)
	if err != nil {
		return nil, err
	}

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&value); err != nil {
				return nil, err
			}
			values = append(values, value)
		}
	}

	return values, nil
}

// Exists returns true if at least one record is found
func (trx Trx) Exists(command string, args ...interface{}) (bool, error) {
	rows, err := trx.Query(command, args...)
	if rows != nil {
		defer rows.Close()
		return rows.Next(), nil
	}
	return false, err
}

// Count returns number of rows
func (trx Trx) Count(command string, args ...interface{}) (int, error) {
	rows, err := trx.Query(command, args...)
	defer rows.Close()
	count := 0
	for rows != nil && rows.Next() {
		count++
	}
	return count, err
}

//Select all matched rows bind to given data
func (trx Trx) Select(data interface{}, command string, args ...interface{}) error {
	rows, err := trx.Query(command, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	sliceType := reflect.TypeOf(data).Elem()
	items := reflect.New(sliceType).Elem()
	itemType := sliceType.Elem().Elem()
	var columns []string
	for rows.Next() {
		if columns == nil {
			columns, _ = rows.Columns()
		}
		item := reflect.New(itemType)
		if err = trx.mapper.Map(item.Interface(), columns, rows.Scan); err != nil {
			return err
		}
		items = reflect.Append(items, item)
	}

	reflect.Indirect(reflect.ValueOf(data)).Set(items)
	return nil
}

// Commit current transaction
func (trx Trx) Commit() error {
	return trx.tx.Commit()
}

// Rollback current transaction
func (trx Trx) Rollback() error {
	return trx.tx.Rollback()
}
