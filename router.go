package main

import (
	"net/http"

	"strings"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/feedback"
	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/infra"
	"github.com/WeCanHearYou/wechy/app/middlewares"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/oauth"
	"github.com/WeCanHearYou/wechy/app/pkg/env"
	"github.com/WeCanHearYou/wechy/app/storage"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

// WechyServices holds reference to all Wechy services
type WechyServices struct {
	OAuth    oauth.Service
	User     storage.User
	Tenant   storage.Tenant
	Idea     storage.Idea
	Settings *models.WechySettings
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

func group(router *echo.Echo, name string) *echo.Group {
	return router.Group(name)
}

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *WechyServices) *echo.Echo {
	router := echo.New()

	router.Logger = createLogger()
	router.Renderer = app.NewHTMLRenderer(router.Logger)
	router.HTTPErrorHandler = errorHandler

	router.Use(middleware.Gzip())
	router.Static("/favicon.ico", "favicon.ico")
	assetsGroup := group(router, "/assets")
	{
		use(assetsGroup, app.OneYearCache())
		assetsGroup.Static("/", "dist")
	}

	oauthHandlers := identity.OAuth(ctx.Tenant, ctx.OAuth, ctx.User)
	authGroup := group(router, "/oauth")
	{
		use(authGroup, middlewares.HostChecker(env.MustGet("AUTH_ENDPOINT")))

		get(authGroup, "/facebook", oauthHandlers.Login(oauth.FacebookProvider))
		get(authGroup, "/facebook/callback", oauthHandlers.Callback(oauth.FacebookProvider))
		get(authGroup, "/google", oauthHandlers.Login(oauth.GoogleProvider))
		get(authGroup, "/google/callback", oauthHandlers.Callback(oauth.GoogleProvider))
	}

	appGroup := group(router, "")
	{
		use(appGroup, middlewares.MultiTenant(ctx.Tenant))
		use(appGroup, middlewares.JwtGetter(ctx.User))
		use(appGroup, middlewares.JwtSetter())

		get(appGroup, "/", feedback.Handlers(ctx.Idea).List())
		get(appGroup, "/ideas/:number", feedback.Handlers(ctx.Idea).Details())
		get(appGroup, "/logout", oauthHandlers.Logout())
		get(appGroup, "/api/status", infra.Status(ctx.Settings))
	}

	securedGroup := group(router, "/api")
	{
		use(securedGroup, middlewares.MultiTenant(ctx.Tenant))
		use(securedGroup, middlewares.JwtGetter(ctx.User))
		use(securedGroup, middlewares.JwtSetter())
		use(securedGroup, middlewares.IsAuthenticated())

		post(securedGroup, "/ideas", feedback.Handlers(ctx.Idea).PostIdea())
		post(securedGroup, "/ideas/:id/comments", feedback.Handlers(ctx.Idea).PostComment())
	}

	adminGroup := group(router, "/admin")
	{
		use(adminGroup, middlewares.MultiTenant(ctx.Tenant))
		use(adminGroup, middlewares.JwtGetter(ctx.User))
		use(adminGroup, middlewares.JwtSetter())
		use(adminGroup, middlewares.IsAuthenticated())
		use(adminGroup, middlewares.IsAuthorized(models.RoleMember, models.RoleAdministrator))

		get(adminGroup, "", func(ctx app.Context) error {
			return ctx.HTML(http.StatusOK, "Welcome to Admin Page :)")
		})
	}

	return router
}
