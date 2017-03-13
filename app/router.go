package app

import (
	"net/http"

	"strings"

	"github.com/WeCanHearYou/wechy/feedback"
	"github.com/WeCanHearYou/wechy/identity"
	"github.com/WeCanHearYou/wechy/toolbox/env"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// WechySettings is an application-wide settings
type WechySettings struct {
	BuildTime    string
	AuthEndpoint string
}

// WechyServices holds reference to all Wechy services
type WechyServices struct {
	OAuth    identity.OAuthService
	User     identity.UserService
	Tenant   identity.TenantService
	Idea     feedback.IdeaService
	Health   HealthCheckService
	Settings *WechySettings
}

func errorHandler(e error, c echo.Context) {
	if strings.Contains(e.Error(), "code=404") {
		c.Render(http.StatusNotFound, "404.html", echo.Map{})
	} else {
		c.Logger().Error(e)
		c.Render(http.StatusInternalServerError, "500.html", echo.Map{})
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
func GetMainEngine(ctx *WechyServices) *echo.Echo {
	router := echo.New()

	router.Logger = createLogger()
	router.Renderer = NewHTMLRenderer(router.Logger)
	router.HTTPErrorHandler = errorHandler

	router.Static("/favicon.ico", "public/imgs/favicon.ico")
	router.Static("/public", "public")
	router.Static("/assets", "node_modules") //TODO: Don't expose node_modules

	authGroup := router.Group("", identity.HostChecker(env.MustGet("AUTH_ENDPOINT")))
	{
		authGroup.GET("/oauth/facebook", identity.OAuth(ctx.OAuth, ctx.User).Login(identity.OAuthFacebookProvider))
		authGroup.GET("/oauth/facebook/callback", identity.OAuth(ctx.OAuth, ctx.User).Callback(identity.OAuthFacebookProvider))
		authGroup.GET("/oauth/google", identity.OAuth(ctx.OAuth, ctx.User).Login(identity.OAuthGoogleProvider))
		authGroup.GET("/oauth/google/callback", identity.OAuth(ctx.OAuth, ctx.User).Callback(identity.OAuthGoogleProvider))
	}

	appGroup := router.Group("")
	{
		appGroup.Use(identity.JwtGetter())
		appGroup.Use(identity.JwtSetter())
		appGroup.Use(identity.MultiTenant(ctx.Tenant))

		appGroup.GET("/", feedback.Index(ctx.Idea))
		appGroup.GET("/logout", identity.OAuth(ctx.OAuth, ctx.User).Logout())
		appGroup.GET("/status", Status(ctx.Health, ctx.Settings))
	}

	return router
}
