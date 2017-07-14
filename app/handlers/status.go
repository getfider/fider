package handlers

import (
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

//Status returns some useful information
func Status(settings *models.AppSettings) web.HandlerFunc {
	return func(c web.Context) error {
		return c.Ok(web.Map{
			"build":    settings.BuildTime,
			"version":  settings.Version,
			"env":      settings.Environment,
			"compiler": settings.Compiler,
			"now":      time.Now().Format("2006.01.02.150405"),
		})
	}
}
