package handler

import (
	"runtime"
	"time"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/gin-gonic/gin"
)

type statusHandler struct {
	ctx *context.WchyContext
}

// Status creates a new Status HTTP handler
func Status(ctx *context.WchyContext) gin.HandlerFunc {
	return statusHandler{ctx: ctx}.get()
}

func (h statusHandler) get() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"build": h.ctx.Settings.BuildTime,
			"healthy": gin.H{
				"database": h.ctx.Health.IsDatabaseOnline(),
			},
			"version": runtime.Version(),
			"now":     time.Now().Format("2006.01.02.150405"),
		})
	}
}
