package handlers

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq" //

	"github.com/gin-gonic/gin"
)

var db *sql.DB

//GetMainEngine returns main HTTP engine
func GetMainEngine() *gin.Engine {
	db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/status", Status(db))
	return router
}
