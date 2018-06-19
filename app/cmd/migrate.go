package cmd

import (
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/log/database"
	"github.com/getfider/fider/app/pkg/uuid"
)

// RunMigrate run all pending migrations on current DATABASE_URL
// Returns an exitcode, 0 for OK and 1 for ERROR
func RunMigrate() int {
	db := dbx.New()
	logger := database.NewLogger("MIGRATE", db)
	logger.SetProperty(log.PropertyKeyContextID, uuid.NewV4().String())
	db.SetLogger(logger)

	defer db.Close()
	err := db.Migrate("/migrations")
	if err != nil {
		logger.Error(err)
		return 1
	}
	return 0
}
