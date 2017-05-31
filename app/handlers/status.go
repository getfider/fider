package handlers

import (
	"net/http"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/labstack/echo"
)

//Status returns some useful information
func Status(settings *models.AppSettings) web.HandlerFunc {
	return func(c web.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"build":    settings.BuildTime,
			"version":  settings.Version,
			"env":      settings.Environment,
			"compiler": settings.Compiler,
			"now":      time.Now().Format("2006.01.02.150405"),
		})
	}
}
