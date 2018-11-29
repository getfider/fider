package cmd

import (
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/log/console"
	"github.com/getfider/fider/app/pkg/rand"
)

// RunMigrate run all pending migrations on current DATABASE_URL
// Returns an exitcode, 0 for OK and 1 for ERROR
func RunMigrate() int {
	logger := console.NewLogger("MIGRATE")
	logger.SetProperty(log.PropertyKeyContextID, rand.String(32))
	db := dbx.NewWithLogger(logger)
	defer db.Close()

	err := db.Migrate("/migrations")
	if err != nil {
		logger.Error(err)
		return 1
	}
	return 0
}
