package main

import (
	"runtime"

	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/pkg/env"
	"github.com/WeCanHearYou/wechy/app/pkg/oauth"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"
	_ "github.com/mattes/migrate/driver/postgres"
	mig "github.com/mattes/migrate/migrate"

	"fmt"
)

var buildtime string
var version = "0.1.0"

func migrate() {
	fmt.Printf("Running migrations... \n")
	errors, ok := mig.UpSync(env.MustGet("DATABASE_URL"), env.Path("/migrations"))
	if !ok {
		for i, err := range errors {
			fmt.Printf("Error #%d: %s.\n", i, err)
		}

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

	ctx := &WeCHYServices{
		OAuth:  &oauth.HTTPService{},
		Idea:   &postgres.IdeaStorage{DB: db},
		User:   &postgres.UserStorage{DB: db},
		Tenant: &postgres.TenantStorage{DB: db},
		Settings: &models.WeCHYSettings{
			BuildTime:    buildtime,
			Version:      version,
			AuthEndpoint: env.MustGet("AUTH_ENDPOINT"),
			Compiler:     runtime.Version(),
			Environment:  env.Current(),
		},
	}

	e := GetMainEngine(ctx)
	e.Start(":" + env.GetEnvOrDefault("PORT", "3000"))
}
