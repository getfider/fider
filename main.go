package main

import (
	"os"

	"github.com/WeCanHearYou/wchy-api/handlers"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

func init() {
	log.Info("Application is starting...")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	handlers.GetMainEngine().Run(":" + port)
}
