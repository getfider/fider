package main

import (
	"net/http"

	"strings"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/feedback"
	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/infra"
	"github.com/WeCanHearYou/wechy/app/toolbox/env"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

// WechyServices holds reference to all Wechy services
type WechyServices struct {
	OAuth    identity.OAuthService
	User     identity.UserService
	Tenant   identity.TenantService
	Idea     feedback.IdeaService
	Health   infra.HealthCheckService
	Settings *infra.WechySettings
}

func errorHandler(e error, c echo.Context) {
	if strings.Contains(e.Error(), "code=404") {
		c.Logger().Debug(e)
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

func wrapFunc(handler app.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := app.Context{Context: c}
		return handler(ctx)
	}
}

func wrapMiddleware(mw app.MiddlewareFunc) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return wrapFunc(mw(func(c app.Context) error {
			return h(c)
		}))
	}
}

func get(group *echo.Group, path string, handler app.HandlerFunc) {
	group.GET(path, wrapFunc(handler))
}

func post(group *echo.Group, path string, handler app.HandlerFunc) {
	group.POST(path, wrapFunc(handler))
}

func use(group *echo.Group, mw app.MiddlewareFunc) {
	group.Use(wrapMiddleware(mw))
}

func group(router *echo.Echo, name string, middlewares ...app.MiddlewareFunc) *echo.Group {
	var mw []echo.MiddlewareFunc
	for _, m := range middlewares {
		mw = append(mw, wrapMiddleware(m))
	}
	return router.Group(name, mw...)
}

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *WechyServices) *echo.Echo {
	router := echo.New()

	router.Logger = createLogger()
	router.Renderer = app.NewHTMLRenderer(router.Logger)
	router.HTTPErrorHandler = errorHandler

	router.Use(middleware.Gzip())
	router.Static("/favicon.ico", "favicon.ico")
	router.Static("/assets", "dist")

	hostChecker := identity.HostChecker(env.MustGet("AUTH_ENDPOINT"))
	authGroup := group(router, "", hostChecker)
	{
		get(authGroup, "/oauth/facebook", identity.OAuth(ctx.OAuth, ctx.User).Login(identity.OAuthFacebookProvider))
		get(authGroup, "/oauth/facebook/callback", identity.OAuth(ctx.OAuth, ctx.User).Callback(identity.OAuthFacebookProvider))
		get(authGroup, "/oauth/google", identity.OAuth(ctx.OAuth, ctx.User).Login(identity.OAuthGoogleProvider))
		get(authGroup, "/oauth/google/callback", identity.OAuth(ctx.OAuth, ctx.User).Callback(identity.OAuthGoogleProvider))
	}

	appGroup := group(router, "")
	{
		use(appGroup, identity.JwtGetter())
		use(appGroup, identity.JwtSetter())
		use(appGroup, identity.MultiTenant(ctx.Tenant))

		get(appGroup, "/", feedback.Index(ctx.Idea).List())
		get(appGroup, "/ideas/:id", feedback.Index(ctx.Idea).Details())
		get(appGroup, "/logout", identity.OAuth(ctx.OAuth, ctx.User).Logout())

		get(appGroup, "/api/status", infra.Status(ctx.Health, ctx.Settings))
		post(appGroup, "/api/ideas", feedback.Index(ctx.Idea).Post())
		post(appGroup, "/api/ideas/:id/comments", feedback.Index(ctx.Idea).PostComment())
	}

	return router
}
