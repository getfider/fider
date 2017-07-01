package web

import (
	"net/http"
	"strings"

	"fmt"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

//HandlerFunc represents an HTTP handler
type HandlerFunc func(Context) error

//MiddlewareFunc represents an HTTP middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc

//Engine is our web engine wrapper
type Engine struct {
	router *echo.Echo
	Logger *log.Logger
}

//New creates a new Engine
func New(settings *models.AppSettings) *Engine {
	router := echo.New()
	router.Use(middleware.Gzip())

	logger := NewLogger()
	router.Logger = logger

	router.Renderer = NewHTMLRenderer(settings, router.Logger)
	router.HTTPErrorHandler = errorHandler
	return &Engine{router: router, Logger: logger}
}

//Start an HTTP server.
func (e *Engine) Start(address string) {
	cert := env.GetEnvOrDefault("SSL_CERT", "")
	key := env.GetEnvOrDefault("SSL_CERT_KEY", "")

	var err error
	if cert == "" && key == "" {
		err = e.router.Start(address)
	} else {
		err = e.router.StartTLS(address, cert, key)
	}

	if err != nil {
		e.router.Logger.Fatal(err)
	}
}

//Use middleware on root router
func (e *Engine) Use(middleware MiddlewareFunc) {
	e.router.Use(wrapMiddleware(middleware))
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

//Group creates a new router group with prefix
func (e *Engine) Group(preffix string) *Group {
	return &Group{group: e.router.Group(preffix)}
}

//HandleError redirect error to router
func (e *Engine) HandleError(err error, ctx Context) {
	e.router.HTTPErrorHandler(err, ctx)
}

//Use add middleware to sub-routes within the Group
func (g *Group) Use(middleware MiddlewareFunc) {
	g.group.Use(wrapMiddleware(middleware))
}

//Group creates asub-group with prefix
func (g *Group) Group(preffix string) *Group {
	return &Group{group: g.group.Group(preffix)}
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

// NewLogger creates a new logger
func NewLogger() *log.Logger {
	logger := log.New("")
	logger.EnableColor()
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
