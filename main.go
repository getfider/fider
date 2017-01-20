package main

import (
	"os"
	"time"

	"database/sql"

	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	logging "github.com/op/go-logging"
)

var buildtime string
var db *sql.DB
var log = logging.MustGetLogger("main")

func statusHandler(c *gin.Context) {
	isHealthy := true
	_, err := db.Query("SELECT now2()")
	if err != nil {
		log.Error(err)
		isHealthy = false
	}

	c.JSON(200, gin.H{
		"healthy": gin.H{
			"database": isHealthy,
		},
		"build":   buildtime,
		"version": runtime.Version(),
		"now":     time.Now().Format("2006.01.02.150405"),
	})
}

func main() {
	db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/status", statusHandler)

	router.Run(":" + port)
}
