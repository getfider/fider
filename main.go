package main

import (
	"os"
	"time"

	"database/sql"

	"runtime"

	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var buildtime string

func statusHandler(c *gin.Context) {
	isHealthy := true
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		isHealthy = false
	} else {
		_, err := db.Query("SELECT now()")
		if err != nil {
			fmt.Println(err)
			isHealthy = false
		}
	}
	db.Close()

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
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/status", statusHandler)

	router.Run(":" + port)
}
