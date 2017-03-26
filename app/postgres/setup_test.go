package postgres_test

import (
	"database/sql"

	"github.com/WeCanHearYou/wechy/app/toolbox/env"
	_ "github.com/lib/pq"
)

func setup() *sql.DB {
	db, err := sql.Open("postgres", env.MustGet("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	return db
}

func teardown(db *sql.DB) {
	execute(db, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	execute(db, "TRUNCATE TABLE tenants RESTART IDENTITY CASCADE")

	if db != nil {
		db.Close()
	}
}

func execute(db *sql.DB, command string, args ...interface{}) {
	_, err := db.Exec(command, args...)
	if err != nil {
		panic(err)
	}
}
