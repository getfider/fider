package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"fmt"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/env"
	"github.com/WeCanHearYou/wchy/router"
	"github.com/WeCanHearYou/wchy/service"
)

var ctx *context.WchyContext
var db *sql.DB
var buildtime string

func init() {
	fmt.Printf("Application is starting...\n")
	fmt.Printf("GO_ENV: %s\n", env.Current())

	db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	ctx = &context.WchyContext{
		Health: &service.PostgresHealthCheckService{DB: db},
		Idea:   &service.PostgresIdeaService{DB: db},
		Tenant: &service.PostgresTenantService{DB: db},
		Settings: context.WchySettings{
			BuildTime: buildtime,
		},
	}
}

func main() {
	router.GetMainEngine(ctx).Run(":" + env.GetEnvOrDefault("PORT", "3000"))
}
