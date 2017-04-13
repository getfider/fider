package infra

import (
	"net/http"
	"runtime"
	"time"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/env"
	"github.com/labstack/echo"
)

type statusHandler struct {
	settings *models.WechySettings
}

// Status creates a new Status HTTP handler
func Status(settings *models.WechySettings) app.HandlerFunc {
	return statusHandler{settings}.get()
}

func (h statusHandler) get() app.HandlerFunc {
	return func(c app.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"build":    h.settings.BuildTime,
			"version":  h.settings.Version,
			"env":      env.Current(),
			"compiler": runtime.Version(),
			"now":      time.Now().Format("2006.01.02.150405"),
		})
	}
}
