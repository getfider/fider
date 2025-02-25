package main

import (
	"os"

	"github.com/Spicy-Bush/fider-tarkov-community/app/cmd"
	_ "github.com/lib/pq"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "ping" {
		os.Exit(cmd.RunPing())
	} else if len(args) > 0 && args[0] == "migrate" {
		os.Exit(cmd.RunMigrate())
	} else {
		os.Exit(cmd.RunServer())
	}
}
