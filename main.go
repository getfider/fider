package main

import (
	"embed"
	"os"

	"github.com/getfider/fider/app/cmd"
	"github.com/getfider/fider/app/pkg/i18n"
	_ "github.com/lib/pq"
)

//go:embed locale/*
var locales embed.FS

func main() {
	i18n.Locales = locales

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "ping" {
		os.Exit(cmd.RunPing())
	} else if len(args) > 0 && args[0] == "migrate" {
		os.Exit(cmd.RunMigrate())
	} else {
		os.Exit(cmd.RunServer())
	}
}
