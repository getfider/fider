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
	defer db.Close()
	db.Connection().Exec("DELETE FROM logs WHERE properties->>'ContextID' = 'MyContextID'")

	trx, _ := db.Begin()
	defer trx.Commit()

	logger := database.NewLogger("TEST", db)
	logger.SetLevel(log.INFO)
	logger.SetProperty(log.PropertyKeyContextID, "MyContextID")
	db.SetLogger(logger)

	logger.Infof("2 + 2 is @{Result}", log.Props{"Result": 4})
	logger.Debugf("2 + 2 is @{Result}", log.Props{"Result": 4})

	Expect(func() int {
		count, err := trx.Count("SELECT * FROM logs WHERE properties->>'ContextID' = 'MyContextID'")
		Expect(err).IsNil()
		return count
	}).EventuallyEquals(1)

	var entries []*logEntry
	err := trx.Select(&entries, "SELECT tag, level, text FROM logs WHERE properties->>'ContextID' = 'MyContextID'")
	Expect(err).IsNil()
	Expect(entries).HasLen(1)
	Expect(entries[0].Tag).Equals("TEST")
	Expect(entries[0].Level).Equals("INFO")
	Expect(entries[0].Text).Equals("2 + 2 is 4")
}
