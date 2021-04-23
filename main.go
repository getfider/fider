package main

import (
	"os"
	"runtime"

	"github.com/getfider/fider/app/cmd"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	_ "github.com/lib/pq"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

var (
	// replaced during CI build
	buildtime   = ""
	// buildnumber = "local"

	// Use this for non-stable releases
	// version = "x.y.z-" + buildnumber

	// Use this for stable releases
	// version = "x.y.z"

	version = "0.18.1"
)

func main() {
	client := appinsights.NewTelemetryClient(env.Config.InstrumentationKey)
	settings := &models.SystemSettings{
		BuildTime:       buildtime,
		Version:         version,
		Compiler:        runtime.Version(),
		Environment:     env.Config.Environment,
		GoogleAnalytics: env.Config.GoogleAnalytics,
		Mode:            env.Config.HostMode,
		Domain:          env.MultiTenantDomain(),
		HasLegal:        env.HasLegal(),
	}

	client.TrackEvent("Main executed")

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "ping" {
		os.Exit(cmd.RunPing())
	} else if len(args) > 0 && args[0] == "migrate" {
		os.Exit(cmd.RunMigrate())
	} else {
		os.Exit(cmd.RunServer(settings))
	}
}
