package dbx

import (
	"context"
	"database/sql"
	"io/ioutil"
	"reflect"
	"time"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"

	"github.com/lib/pq"
)

var conn *sql.DB
var rowMapper *RowMapper

func init() {
	var err error
	conn, err = sql.Open("postgres", env.Config.Database.URL)
	if err != nil {
		panic(wrap(err, "failed to open connection to the database"))
	}

	conn.SetMaxIdleConns(env.Config.Database.MaxIdleConns)
	conn.SetMaxOpenConns(env.Config.Database.MaxOpenConns)
	rowMapper = NewRowMapper()
}

func Connection() *sql.DB {
	return conn
}

// Ping checks if current database connection is healthy
func Ping() error {
	_, err := conn.Exec("SELECT 1")
	return err
}

// BeginTx returns a new SQL transaction
func BeginTx(ctx context.Context) (*Trx, error) {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, wrap(err, "failed to start new transaction")
	}
	return &Trx{tx: tx, ctx: ctx}, nil
}

func load(path string) {
	content, err := ioutil.ReadFile(env.Path(path))
	if err != nil {
		panic(wrap(err, "failed to read file %s", path))
	}

	_, err = conn.Exec(string(content))
	if err != nil {
		panic(wrap(err, "failed to execute %s", path))
	}
}

// Seed clean and insert new seed data for testing
func Seed() {
	if env.IsTest() {
		load("/app/pkg/dbx/setup.sql")
	}
}

//Trx represents a Database transaction
type Trx struct {
	tx  *sql.Tx
	ctx context.Context
}

var formatter = strings.NewReplacer("\t", "", "\n", " ")

// Execute given SQL command
func (trx *Trx) Execute(command string, args ...interface{}) (int64, error) {
	if log.IsEnabled(log.DEBUG) {
		command = formatter.Replace(command)
		start := time.Now()
		defer func() {
			log.Debugf(trx.ctx, "@{Command:yellow} @{Args:blue} executed in @{ElapsedMs:magenta}ms", dto.Props{
				"Command":   command,
				"Args":      args,
				"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
			})
		}()
	}

	result, err := trx.tx.ExecContext(trx.ctx, command, args...)
	if err != nil {
		return 0, wrap(err, "failed to execute trx.Execute")
	}

	rows, _ := result.RowsAffected()
	return rows, nil
}

// Scalar returns first row and first column
func (trx *Trx) Scalar(data interface{}, command string, args ...interface{}) error {
	if log.IsEnabled(log.DEBUG) {
		command = formatter.Replace(command)
		start := time.Now()
		defer func() {
			log.Debugf(trx.ctx, "@{Command:yellow} @{Args:blue} executed in @{ElapsedMs:magenta}ms", dto.Props{
				"Command":   command,
				"Args":      args,
				"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
			})
		}()
	}

	row := trx.tx.QueryRowContext(trx.ctx, command, args...)
	err := row.Scan(data)
	if err != nil {
		if err == sql.ErrNoRows {
			return app.ErrNotFound
		}
		return wrap(err, "failed to execute trx.Scalar")
	}
	return nil
}

// Get first row and bind to given data
func (trx *Trx) Get(data interface{}, command string, args ...interface{}) error {
	if log.IsEnabled(log.DEBUG) {
		command = formatter.Replace(command)
		start := time.Now()
		defer func() {
			log.Debugf(trx.ctx, "@{Command:yellow} @{Args:blue} executed in @{ElapsedMs:magenta}ms", dto.Props{
				"Command":   command,
				"Args":      args,
				"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
			})
		}()
	}

	rows, err := trx.tx.QueryContext(trx.ctx, command, args...)
	if err != nil {
		return wrap(err, "failed to execute trx.Get")
	}

	defer rows.Close()
	if rows.Next() {
		columns, _ := rows.Columns()
		err := rowMapper.Map(data, columns, rows.Scan)
		if err != nil {
			return wrap(err, "failed to map result to model")
		}
		return nil
	}

	return app.ErrNotFound
}

// Exists returns true if at least one record is found
func (trx *Trx) Exists(command string, args ...interface{}) (bool, error) {
	if log.IsEnabled(log.DEBUG) {
		command = formatter.Replace(command)
		start := time.Now()
		defer func() {
			log.Debugf(trx.ctx, "@{Command:yellow} @{Args:blue} executed in @{ElapsedMs:magenta}ms", dto.Props{
				"Command":   command,
				"Args":      args,
				"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
			})
		}()
	}

	rows, err := trx.tx.QueryContext(trx.ctx, command, args...)
	if err != nil {
		return false, wrap(err, "failed to execute trx.Exists")
	}

	defer rows.Close()
	return rows.Next(), nil
}

// Count returns number of rows
func (trx *Trx) Count(command string, args ...interface{}) (int, error) {
	if log.IsEnabled(log.DEBUG) {
		command = formatter.Replace(command)
		start := time.Now()
		defer func() {
			log.Debugf(trx.ctx, "@{Command:yellow} @{Args:blue} executed in @{ElapsedMs:magenta}ms", dto.Props{
				"Command":   command,
				"Args":      args,
				"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
			})
		}()
	}

	rows, err := trx.tx.QueryContext(trx.ctx, command, args...)
	if err != nil {
		return 0, wrap(err, "failed to execute trx.Count")
	}

	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
	}
	return count, nil
}

//Select all matched rows bind to given data
func (trx *Trx) Select(data interface{}, command string, args ...interface{}) error {
	if log.IsEnabled(log.DEBUG) {
		command = formatter.Replace(command)
		start := time.Now()
		defer func() {
			log.Debugf(trx.ctx, "@{Command:yellow} @{Args:blue} executed in @{ElapsedMs:magenta}ms", dto.Props{
				"Command":   command,
				"Args":      args,
				"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
			})
		}()
	}

	rows, err := trx.tx.QueryContext(trx.ctx, command, args...)
	if err != nil {
		return wrap(err, "failed to execute trx.Select")
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
		if err = rowMapper.Map(item.Interface(), columns, rows.Scan); err != nil {
			return wrap(err, "failed to map result to model")
		}
		items = reflect.Append(items, item)
	}

	if items.Len() > 0 {
		reflect.Indirect(reflect.ValueOf(data)).Set(items)
	}
	return nil
}

//Query all matched rows and return raw sql.Rows
func (trx *Trx) Query(command string, args ...interface{}) (*sql.Rows, error) {
	if log.IsEnabled(log.DEBUG) {
		command = formatter.Replace(command)
		start := time.Now()
		defer func() {
			log.Debugf(trx.ctx, "@{Command:yellow} @{Args:blue} executed in @{ElapsedMs:magenta}ms", dto.Props{
				"Command":   command,
				"Args":      args,
				"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
			})
		}()
	}

	rows, err := trx.tx.QueryContext(trx.ctx, command, args...)
	if err != nil {
		return nil, wrap(err, "failed to execute trx.Select")
	}
	return rows, nil
}

// Commit current transaction
func (trx *Trx) Commit() error {
	err := trx.tx.Commit()
	if err != nil && err != sql.ErrTxDone {
		return wrap(err, "failed to commit transaction")
	}
	return nil
}

// MustCommit current transaction
func (trx *Trx) MustCommit() {
	err := trx.Commit()
	if err != nil {
		panic(err)
	}
}

// Rollback current transaction
func (trx *Trx) Rollback() error {
	if trx.tx == nil {
		return nil
	}

	err := trx.tx.Rollback()
	if err != nil && err != sql.ErrTxDone {
		return wrap(err, "failed to rollback transaction")
	}

	return nil
}

// MustRollback current transaction
func (trx *Trx) MustRollback() {
	err := trx.Rollback()
	if err != nil {
		panic(err)
	}
}

func wrap(err error, format string, a ...interface{}) error {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "57014" { //query canceled
			return errors.Wrap(context.Canceled, format, a...)
		}
	}
	return errors.Wrap(err, format, a...)
}
