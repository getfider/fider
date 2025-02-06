package cmd

import (
	"context"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/dbx"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/log"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/rand"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/log/console"
)

// RunMigrate run all pending migrations on current DATABASE_URL
// Returns an exitcode, 0 for OK and 1 for ERROR
func RunMigrate() int {
	bus.Init()

	ctx := log.WithProperties(context.Background(), dto.Props{
		log.PropertyKeyTag:       "MIGRATE",
		log.PropertyKeyContextID: rand.String(32),
	})

	err := dbx.Migrate(ctx, "/migrations")
	if err != nil {
		log.Error(ctx, err)
		return 1
	}
	return 0
}
