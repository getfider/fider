package cmd

import (
	"context"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/rand"
	_ "github.com/getfider/fider/app/services/log/console"
)

// RunMigrate run all pending migrations on current DATABASE_URL
// Returns an exitcode, 0 for OK and 1 for ERROR
func RunMigrate() int {
	bus.Init()

	ctx := context.Background()
	ctx = log.SetProperty(ctx, log.PropertyKeyTag, "MIGRATE")
	ctx = log.SetProperty(ctx, log.PropertyKeyContextID, rand.String(32))

	db := dbx.New()
	defer db.Close()

	err := db.Migrate(ctx, "/migrations")
	if err != nil {
		log.Error(ctx, err)
		return 1
	}
	return 0
}
