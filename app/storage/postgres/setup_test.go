package postgres_test

import (
	"os"
	"testing"

	"github.com/getfider/fider/app/pkg/dbx"
)

func TestMain(m *testing.M) {
	db, _ := dbx.New()
	db.Migrate()
	db.Close()

	code := m.Run()
	os.Exit(code)
}
