package main

import (
	"runtime"

	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/pkg/env"
	"github.com/WeCanHearYou/wechy/app/pkg/oauth"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"

	"fmt"
)

var buildtime string
var version = "0.2.0"

func main() {
	fmt.Printf("Application is starting...\n")
	fmt.Printf("GO_ENV: %s\n", env.Current())

	db, err := dbx.New()
	if err != nil {
		panic(err)
	}
	db.Migrate()

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
