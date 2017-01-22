package handlers

import (
	"runtime"
	"time"

	"github.com/WeCanHearYou/wchy-api/context"
	"github.com/gin-gonic/gin"
)

var buildtime string

type statusHandler struct {
	ctx context.WchyContext
}

// Status creates a new Status HTTP handler
func Status(ctx context.WchyContext) gin.HandlerFunc {
	return statusHandler{ctx: ctx}.get()
}

func (h statusHandler) get() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"healthy": gin.H{
				"database": h.ctx.Health.IsDatabaseOnline(),
			},
			"build":   buildtime,
			"version": runtime.Version(),
			"now":     time.Now().Format("2006.01.02.150405"),
		})
	}
}
