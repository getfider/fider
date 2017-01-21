package main

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/WeCanHearYou/wchy-api/handlers"
	logging "github.com/op/go-logging"
)

var db *sql.DB
var log = logging.MustGetLogger("main")

func init() {
	log.Info("Application is starting...")
}

func main() {
	db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/status", handlers.Status(db))

	router.Run(":" + port)
}
