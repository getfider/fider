package router

import (
	"net/http"

	"strings"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/env"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func errorHandler(e error, c echo.Context) {
	if strings.Contains(e.Error(), "code=404") {
		c.Logger().Debugf("%s not found.", c.Request().URL)
		c.NoContent(http.StatusNotFound)
	} else {
		c.Logger().Error(e)
		c.NoContent(http.StatusInternalServerError)
	}
}

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *context.WchyContext) *echo.Echo {
	router := echo.New()

	logger := log.New("")
	logger.SetHeader(`${level} [${time_rfc3339}] ${short_file}:${line}`)

	if env.IsProduction() {
		logger.SetLevel(log.INFO)
	} else {
		logger.SetLevel(log.DEBUG)
	}

	router.Logger = logger
	router.Renderer = NewHTMLRenderer()
	router.HTTPErrorHandler = errorHandler

	router.Static("/favicon.ico", "public/imgs/favicon.ico")
	router.Static("/public", "public")
	router.Static("/assets", "node_modules") //TODO: Don't expose node_modules

	router.GET("/oauth", handler.OAuth(ctx))

	app := router.Group("", MultiTenant(ctx))
	{
		app.GET("/", handler.Index(ctx))

		api := app.Group("/api")
		{
			api.GET("/status", handler.Status(ctx))
			api.GET("/tenants/:domain", handler.TenantByDomain(ctx))
		}
	}

	return router
}
