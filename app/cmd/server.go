package cmd

import (
	"fmt"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
)

//RunServer starts the Fider Server
//Returns an exitcode, 0 for OK and 1 for ERROR
func RunServer(settings *models.SystemSettings) int {

	fmt.Printf("Application is starting...\n")
	fmt.Printf("GO_ENV: %s\n", env.Current())

	e := getMainEngine(settings)
	e.Start(":" + env.GetEnvOrDefault("PORT", "3000"))

	return 0
}
