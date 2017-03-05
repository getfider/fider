package main

import (
	"database/sql"

	"github.com/WeCanHearYou/wchy/app/auth"
	"github.com/WeCanHearYou/wchy/app/context"
	"github.com/WeCanHearYou/wchy/app/env"
	"github.com/WeCanHearYou/wchy/app/router"
	"github.com/WeCanHearYou/wchy/app/service"
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

	ctx := &context.WchyContext{
		OAuth:  &auth.HTTPOAuthService{},
		Health: &service.PostgresHealthCheckService{DB: db},
		Idea:   &service.PostgresIdeaService{DB: db},
		User:   &service.PostgresUserService{DB: db},
		Tenant: &service.PostgresTenantService{DB: db},
		Settings: context.WchySettings{
			BuildTime:    buildtime,
			AuthEndpoint: env.MustGet("AUTH_ENDPOINT"),
		},
	}

	e := router.GetMainEngine(ctx)
	e.Logger.Fatal(e.Start(":" + env.GetEnvOrDefault("PORT", "3000")))
}
