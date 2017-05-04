package web

import (
	"net/http"
	"strings"

	"fmt"

	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/env"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

//Engine is our web engine wrapper
type Engine struct {
	router *echo.Echo
}

//New creates a new Engine
func New(settings *models.AppSettings) *Engine {
	router := echo.New()
	router.Use(middleware.Gzip())
	router.Static("/favicon.ico", "favicon.ico")
	router.Logger = createLogger()
	router.Renderer = NewHTMLRenderer(settings, router.Logger)
	router.HTTPErrorHandler = errorHandler
	return &Engine{router: router}
}

//Start an HTTP server.
func (e *Engine) Start(address string) {
	e.router.Logger.Fatal(e.router.Start(address))
}

//Group is our router group wrapper
type Group struct {
	group *echo.Group
}

//NewContext creates and return a new context
func (e *Engine) NewContext(req *http.Request, w http.ResponseWriter) Context {
	context := e.router.NewContext(req, w)
	return Context{Context: context}
}

//Group creates a new router group with prefix and optional group-level middleware
func (e *Engine) Group(preffix string) *Group {
	return &Group{group: e.router.Group(preffix)}
}

//HandleError redirect error to router
func (e *Engine) HandleError(err error, ctx Context) {
	e.router.HTTPErrorHandler(err, ctx)
}

//Use add middleware to routes within the main Engine
func (e *Engine) Use(middleware MiddlewareFunc) {
	e.router.Use(wrapMiddleware(middleware))
}

//Use add middleware to sub-routes within the Group
func (g *Group) Use(middleware MiddlewareFunc) {
	g.group.Use(wrapMiddleware(middleware))
}

//Static return files from given folder
func (g *Group) Static(prefix, root string) {
	g.group.Static(prefix, root)
}

//Get handles HTTP GET requests
func (g *Group) Get(path string, handler HandlerFunc) {
	g.group.GET(path, wrapFunc(handler))
}

//Post handles HTTP POST requests
func (g *Group) Post(path string, handler HandlerFunc) {
	g.group.POST(path, wrapFunc(handler))
}

func errorHandler(e error, c echo.Context) {
	if strings.Contains(e.Error(), "code=404") {
		c.Logger().Debug(fmt.Sprintf("%s [%s] %s", e, c.Request().Method, c.Request().URL.String()))
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

func wrapMiddleware(mw MiddlewareFunc) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return wrapFunc(mw(func(c Context) error {
			return h(c)
		}))
	}
}

func wrapFunc(handler HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := Context{Context: c}
		return handler(ctx)
	}
}
