package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/julienschmidt/httprouter"
)

//HandlerFunc represents an HTTP handler
type HandlerFunc func(Context) error

//MiddlewareFunc represents an HTTP middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc

//Engine is our web engine wrapper
type Engine struct {
	mux      *httprouter.Router
	logger   log.Logger
	renderer *Renderer
	binder   *DefaultBinder
}

//New creates a new Engine
func New(settings *models.AppSettings) *Engine {
	logger := NewLogger()
	router := &Engine{
		mux:      httprouter.New(),
		logger:   logger,
		renderer: NewRenderer(settings, logger),
		binder:   NewDefaultBinder(),
	}

	router.mux.NotFound = func(res http.ResponseWriter, req *http.Request) {
		ctx := router.NewContext(res, req, nil)
		ctx.NotFound()
	}
	return router
}

//Start an HTTP server.
func (e *Engine) Start(address string) {
	cert := env.GetEnvOrDefault("SSL_CERT", "")
	key := env.GetEnvOrDefault("SSL_CERT_KEY", "")

	var err error
	if cert == "" && key == "" {
		e.logger.Infof("http server started on %s", address)
		err = http.ListenAndServe(address, e.mux)
	} else {
		e.logger.Infof("https server started on %s", address)
		err = http.ListenAndServeTLS(address, cert, key, e.mux)
	}

	if err != nil {
		e.logger.Error(err)
	}
}

//NewContext creates and return a new context
func (e *Engine) NewContext(res http.ResponseWriter, req *http.Request, params StringMap) Context {
	return Context{
		engine: e,
		res:    res,
		req:    req,
		logger: e.logger,
		params: params,
	}
}

//Logger returns current logger
func (e *Engine) Logger() log.Logger {
	return e.logger
}

//Group creates a new route group
func (e *Engine) Group() *Group {
	g := &Group{
		engine:      e,
		middlewares: make([]MiddlewareFunc, 0),
	}
	return g
}

//Group is our router group wrapper
type Group struct {
	engine      *Engine
	middlewares []MiddlewareFunc
}

//Group creates a new route group
func (g *Group) Group() *Group {
	g2 := &Group{
		engine:      g.engine,
		middlewares: g.middlewares,
	}
	return g2
}

//Use adds a middleware to current route stack
func (g *Group) Use(middleware MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middleware)
}

//Get handles HTTP GET requests
func (g *Group) Get(path string, handler HandlerFunc) {
	g.engine.mux.Handle("GET", path, g.handler(handler))
}

//Post handles HTTP POST requests
func (g *Group) Post(path string, handler HandlerFunc) {
	g.engine.mux.Handle("POST", path, g.handler(handler))
}

func (g *Group) handler(handler HandlerFunc) httprouter.Handle {
	next := handler
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		next = g.middlewares[i](next)
	}
	var h = func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		params := make(StringMap, 0)
		for _, p := range ps {
			params[p.Key] = p.Value
		}
		ctx := g.engine.NewContext(res, req, params)
		if err := next(ctx); err != nil {
			ctx.Failure(err)
		}
	}
	return h
}

// Static return files from given folder
func (g *Group) Static(prefix, root string) {

	fi, err := os.Stat(root)
	if err != nil {
		panic(fmt.Sprintf("Path '%s' not found", root))
	}

	var h HandlerFunc
	if fi.IsDir() {
		h = func(c Context) error {
			path := root + c.Param("filepath")
			fi, err := os.Stat(path)
			if err == nil && !fi.IsDir() {
				http.ServeFile(c.Response(), c.Request(), path)
				return nil
			}
			return c.NotFound()
		}
	} else {
		h = func(c Context) error {
			http.ServeFile(c.Response(), c.Request(), root)
			return nil
		}
	}
	g.engine.mux.Handle("GET", prefix, g.handler(h))
}

// NewLogger creates a new logger
func NewLogger() log.Logger {
	logger := log.NewConsoleLogger()

	if env.IsProduction() {
		logger.SetLevel(log.INFO)
	} else {
		logger.SetLevel(log.DEBUG)
	}

	return logger
}
