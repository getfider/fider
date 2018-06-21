package database_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/log/database"
)

type logEntry struct {
	Tag   string `db:"tag"`
	Level string `db:"level"`
	Text  string `db:"text"`
}

func TestLog(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	db.Seed()
	defer db.Close()
	trx, _ := db.Begin()
	defer trx.Rollback()

	logger := database.NewLogger("TEST", db)
	logger.SetLevel(log.INFO)
	db.SetLogger(logger)

	logger.Infof("2 + 2 is @{Result}", log.Props{"Result": 4})
	logger.Debugf("2 + 2 is @{Result}", log.Props{"Result": 4})

	Expect(func() int {
		count, err := trx.Count("SELECT id FROM logs")
		Expect(err).IsNil()
		return count
	}).EventuallyEquals(1)

	entry := logEntry{}
	err := trx.Get(&entry, "SELECT tag, level, text FROM logs LIMIT 1")
	Expect(err).IsNil()
	Expect(entry.Tag).Equals("TEST")
	Expect(entry.Level).Equals("INFO")
	Expect(entry.Text).Equals("2 + 2 is 4")
}
