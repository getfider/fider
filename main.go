package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"fmt"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/env"
	"github.com/WeCanHearYou/wchy/router"
	"github.com/WeCanHearYou/wchy/service"
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
