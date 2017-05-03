package handlers

import (
	"net/http"
	"time"

	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/labstack/echo"
)

type statusHandler struct {
	settings *models.AppSettings
}

// Status creates a new Status HTTP handler
func Status(settings *models.AppSettings) web.HandlerFunc {
	return statusHandler{settings}.get()
}

func (h statusHandler) get() web.HandlerFunc {
	return func(c web.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"build":    h.settings.BuildTime,
			"version":  h.settings.Version,
			"env":      h.settings.Environment,
			"compiler": h.settings.Compiler,
			"now":      time.Now().Format("2006.01.02.150405"),
		})
	}
}
