package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/labstack/echo"
)

type statusHandler struct {
	ctx *context.WchyContext
}

// Status creates a new Status HTTP handler
func Status(ctx *context.WchyContext) echo.HandlerFunc {
	return statusHandler{ctx: ctx}.get()
}

func (h statusHandler) get() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"build": h.ctx.Settings.BuildTime,
			"healthy": echo.Map{
				"database": h.ctx.Health.IsDatabaseOnline(),
			},
			"version": runtime.Version(),
			"now":     time.Now().Format("2006.01.02.150405"),
		})
	}
}
