package main

import (
	"runtime"

	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/pkg/env"
	"github.com/WeCanHearYou/wechy/app/pkg/oauth"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"
	_ "github.com/lib/pq"
	mig "github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"fmt"
)

var buildtime string
var version = "0.2.0"

func migrate() {
	fmt.Printf("Running migrations... \n")
	m, err := mig.New(
		"file://./"+env.Path("/migrations"),
		env.MustGet("DATABASE_URL"),
	)

	if err == nil {
		err = m.Up()
	}

	if err != nil && err != mig.ErrNoChange {
		fmt.Printf("Error: %s.\n", err)

		panic("Migrations failed.")
	} else {
		fmt.Printf("Migrations finished with success.\n")
	}
}

func init() {
	fmt.Printf("Application is starting...\n")
	fmt.Printf("GO_ENV: %s\n", env.Current())
	migrate()
}

func main() {
	db, err := dbx.New()
	if err != nil {
		panic(err)
	}

	ctx := &AppServices{
		OAuth:  &oauth.HTTPService{},
		Idea:   &postgres.IdeaStorage{DB: db},
		User:   &postgres.UserStorage{DB: db},
		Tenant: &postgres.TenantStorage{DB: db},
		Settings: &models.AppSettings{
			BuildTime:   buildtime,
			Version:     version,
			Compiler:    runtime.Version(),
			Environment: env.Current(),
		},
	}

	e := GetMainEngine(ctx)
	e.Start(":" + env.GetEnvOrDefault("PORT", "3000"))
}
