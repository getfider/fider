package main

import (
	"database/sql"

	"github.com/WeCanHearYou/wchy/app"
	"github.com/WeCanHearYou/wchy/identity"
	"github.com/WeCanHearYou/wchy/postgres"
	"github.com/WeCanHearYou/wchy/toolbox/env"
	_ "github.com/lib/pq"

	"fmt"
)

var buildtime string

func main() {
	fmt.Printf("Application is starting...\n")
	fmt.Printf("GO_ENV: %s\n", env.Current())

	db, err := sql.Open("postgres", env.MustGet("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	ctx := &app.WchyServices{
		OAuth:  &identity.HTTPOAuthService{},
		Health: &postgres.HealthCheckService{DB: db},
		Idea:   &postgres.IdeaService{DB: db},
		User:   &postgres.UserService{DB: db},
		Tenant: &postgres.TenantService{DB: db},
		Settings: &app.WchySettings{
			BuildTime:    buildtime,
			AuthEndpoint: env.MustGet("AUTH_ENDPOINT"),
		},
	}

	e := app.GetMainEngine(ctx)
	e.Logger.Fatal(e.Start(":" + env.GetEnvOrDefault("PORT", "3000")))
}
