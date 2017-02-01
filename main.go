package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/env"
	"github.com/WeCanHearYou/wchy/router"
	"github.com/WeCanHearYou/wchy/services"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")
var ctx context.WchyContext
var settings context.WchySettings
var db *sql.DB
var buildtime string

func init() {
	log.Info("Application is starting...")
	log.Infof("GO_ENV: %s", env.Current())

	db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	ctx = context.WchyContext{
		Health: &services.PostgresHealthCheckService{DB: db},
		Tenant: &services.PostgresTenantService{DB: db},
		Settings: context.WchySettings{
			BuildTime: buildtime,
		},
	}
}

func main() {
	router.GetMainEngine(ctx).Run(":" + env.GetEnvOrDefault("PORT", "3000"))
}
