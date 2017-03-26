package main

import (
	"database/sql"

	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/infra"
	"github.com/WeCanHearYou/wechy/app/postgres"
	"github.com/WeCanHearYou/wechy/app/toolbox/env"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/driver/postgres"
	mig "github.com/mattes/migrate/migrate"

	"fmt"
)

var buildtime string
var version = "0.0.0"

func migrate() {
	if env.IsTest() {
		return
	}

	fmt.Printf("Running migrations... \n")
	errors, ok := mig.UpSync(env.MustGet("DATABASE_URL"), "./migrations")
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
	db, err := sql.Open("postgres", env.MustGet("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	ctx := &WechyServices{
		OAuth:  &identity.HTTPOAuthService{},
		Health: &postgres.HealthCheckService{DB: db},
		Idea:   &postgres.IdeaService{DB: db},
		User:   &postgres.UserService{DB: db},
		Tenant: &postgres.TenantService{DB: db},
		Settings: &infra.WechySettings{
			BuildTime:    buildtime,
			Version:      version,
			AuthEndpoint: env.MustGet("AUTH_ENDPOINT"),
		},
	}

	e := GetMainEngine(ctx)
	e.Logger.Fatal(e.Start(":" + env.GetEnvOrDefault("PORT", "3000")))
}
