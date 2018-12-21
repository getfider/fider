package main

import (
	"os"
	"runtime"

	"github.com/getfider/fider/app/cmd"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	_ "github.com/lib/pq"
)

var (
	// replaced during CI build
	buildtime   = ""
	buildnumber = "local"

	// Use this for non-stable releases
	// version = "x.y.z-" + buildnumber

	// Use this for stable releases
	// version = "x.y.z"

	version = "0.17.1"
)

func main() {
	settings := &models.SystemSettings{
		BuildTime:       buildtime,
		Version:         version,
		Compiler:        runtime.Version(),
		Environment:     env.Current(),
		GoogleAnalytics: env.GetEnvOrDefault("GOOGLE_ANALYTICS", ""),
		Mode:            env.Mode(),
		Domain:          env.MultiTenantDomain(),
		HasLegal:        env.HasLegal(),
	}

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "ping" {
		os.Exit(cmd.RunPing())
	} else if len(args) > 0 && args[0] == "migrate" {
		os.Exit(cmd.RunMigrate())
	} else {
		os.Exit(cmd.RunServer(settings))
	}
}
