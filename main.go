package main

import (
	"runtime"

	"fmt"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

var buildtime string
var version = "0.12.0-beta"

func main() {
	fmt.Printf("Application is starting...\n")
	fmt.Printf("GO_ENV: %s\n", env.Current())

	settings := &models.SystemSettings{
		BuildTime:       buildtime,
		Version:         version,
		Compiler:        runtime.Version(),
		Environment:     env.Current(),
		GoogleAnalytics: env.GetEnvOrDefault("GOOGLE_ANALYTICS", ""),
		Mode:            env.Mode(),
		Domain:          env.MultiTenantDomain(),
	}

	e := GetMainEngine(settings)
	e.Start(":" + env.GetEnvOrDefault("PORT", "3000"))
}
