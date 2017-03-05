package router

import (
	"net/http"

	"strings"

	"github.com/WeCanHearYou/wchy/app/auth"
	"github.com/WeCanHearYou/wchy/app/context"
	"github.com/WeCanHearYou/wchy/app/env"
	"github.com/WeCanHearYou/wchy/app/handler"
	"github.com/WeCanHearYou/wchy/app/router/middleware"
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

func createLogger() echo.Logger {
	logger := log.New("")
	logger.SetHeader(`${level} [${time_rfc3339}]`)

	if env.IsProduction() {
		logger.SetLevel(log.INFO)
	} else {
		logger.SetLevel(log.DEBUG)
	}

	return logger
}

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *context.WchyContext) *echo.Echo {
	router := echo.New()

	router.Logger = createLogger()
	router.Renderer = NewHTMLRenderer(router.Logger)
	router.HTTPErrorHandler = errorHandler

	router.Static("/favicon.ico", "public/imgs/favicon.ico")
	router.Static("/public", "public")
	router.Static("/assets", "node_modules") //TODO: Don't expose node_modules

	authGroup := router.Group("", middleware.HostChecker(env.MustGet("AUTH_ENDPOINT")))
	{
		authGroup.GET("/oauth/facebook", handler.OAuth(ctx).Login(auth.OAuthFacebookProvider))
		authGroup.GET("/oauth/facebook/callback", handler.OAuth(ctx).Callback(auth.OAuthFacebookProvider))
		authGroup.GET("/oauth/google", handler.OAuth(ctx).Login(auth.OAuthGoogleProvider))
		authGroup.GET("/oauth/google/callback", handler.OAuth(ctx).Callback(auth.OAuthGoogleProvider))
	}

	app := router.Group("")
	{
		app.Use(middleware.JwtGetter())
		app.Use(middleware.JwtSetter())
		app.Use(middleware.MultiTenant(ctx))

		app.GET("/", handler.Index(ctx))
		app.GET("/logout", handler.OAuth(ctx).Logout())

		api := app.Group("/api")
		{
			api.GET("/status", handler.Status(ctx))
			api.GET("/tenants/:domain", handler.TenantByDomain(ctx))
		}
	}

	return router
}
