package main

import (
	"os"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func healthyHandler(c *gin.Context) {
	isHealthy := true
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		isHealthy = false
	} else {
		_, err := db.Query("SELECT now()")
		if err != nil {
			isHealthy = false
		}
	}
	db.Close()

	c.JSON(200, gin.H{"isHealthy": isHealthy})
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/healthy", healthyHandler)

	router.Run(":" + port)
}
