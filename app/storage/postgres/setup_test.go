package postgres_test

import (
	"os"
	"testing"

	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	_ "github.com/mattes/migrate/driver/postgres"
)

func TestMain(m *testing.M) {
	db, _ := dbx.New()
	db.Migrate()
	db.Close()

	code := m.Run()
	os.Exit(code)
}
