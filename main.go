package main

import (
	"runtime"

	"fmt"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
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

	settings := &models.AppSettings{
		BuildTime:   buildtime,
		Version:     version,
		Compiler:    runtime.Version(),
		Environment: env.Current(),
	}

	e := GetMainEngine(settings, db)
	e.Start(":" + env.GetEnvOrDefault("PORT", "3000"))
}
