package main

import (
	"database/sql"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/identity"
	"github.com/WeCanHearYou/wechy/postgres"
	"github.com/WeCanHearYou/wechy/toolbox/env"
	_ "github.com/lib/pq"

	"fmt"
)

var buildtime string

func init() {
	fmt.Printf("Application is starting...\n")
	fmt.Printf("GO_ENV: %s\n", env.Current())
}

func main() {
	db, err := sql.Open("postgres", env.MustGet("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	ctx := &app.WechyServices{
		OAuth:  &identity.HTTPOAuthService{},
		Health: &postgres.HealthCheckService{DB: db},
		Idea:   &postgres.IdeaService{DB: db},
		User:   &postgres.UserService{DB: db},
		Tenant: &postgres.TenantService{DB: db},
		Settings: &app.WechySettings{
			BuildTime:    buildtime,
			AuthEndpoint: env.MustGet("AUTH_ENDPOINT"),
		},
	}

	e := app.GetMainEngine(ctx)
	e.Logger.Fatal(e.Start(":" + env.GetEnvOrDefault("PORT", "3000")))
}
