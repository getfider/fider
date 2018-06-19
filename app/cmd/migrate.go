package cmd

import (
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/log"
)

// RunMigrate run all pending migrations on current DATABASE_URL
// Returns an exitcode, 0 for OK and 1 for ERROR
func RunMigrate() int {
	logger := log.NewConsoleLogger("MIGRATE")
	db := dbx.NewWithLogger(logger)
	err := db.Migrate("/migrations")
	if err != nil {
		logger.Error(err)
		return 1
	}
	return 0
}
