package handlers

import (
	"database/sql"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	logging "github.com/op/go-logging"
)

var buildtime string
var log = logging.MustGetLogger("handlers/status")

type statusHandler struct {
	db *sql.DB
}

//Status creates a new Status HTTP handler
func Status(db *sql.DB) gin.HandlerFunc {
	return statusHandler{db: db}.get(db)
}

func (h statusHandler) get(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		isHealthy := true
		if h.db != nil {
			_, err := h.db.Query("SELECT now()")
			if err != nil {
				log.Error(err)
				isHealthy = false
			}
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
}
