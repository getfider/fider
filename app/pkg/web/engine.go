package web

import (
	"context"
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/log/database"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/worker"
	"github.com/julienschmidt/httprouter"
)

var (
	cspBase    = "base-uri 'self'"
	cspDefault = "default-src 'self'"
	cspStyle   = "style-src 'self' 'nonce-%[1]s' https://fonts.googleapis.com %[2]s"
	cspScript  = "script-src 'self' 'nonce-%[1]s' https://cdn.polyfill.io https://www.google-analytics.com %[2]s"
	cspFont    = "font-src 'self' https://fonts.gstatic.com data: %[2]s"
	cspImage   = "img-src 'self' https: data: %[2]s"
	cspObject  = "object-src 'none'"
	cspMedia   = "media-src 'none'"
	cspConnect = "connect-src 'self' https://www.google-analytics.com %[2]s"

	//CspPolicyTemplate is the template used to generate the policy
	CspPolicyTemplate = fmt.Sprintf("%s; %s; %s; %s; %s; %s; %s; %s; %s", cspBase, cspDefault, cspStyle, cspScript, cspImage, cspFont, cspObject, cspMedia, cspConnect)
)

type notFoundHandler struct {
	engine  *Engine
	handler HandlerFunc
}

func (h *notFoundHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := h.engine.NewContext(res, req, nil)
	h.handler(ctx)
}

//HandlerFunc represents an HTTP handler
type HandlerFunc func(Context) error

//MiddlewareFunc represents an HTTP middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc

//Engine is our web engine wrapper
type Engine struct {
	mux         *httprouter.Router
	logger      log.Logger
	renderer    *Renderer
	db          *dbx.Database
	binder      *DefaultBinder
	middlewares []MiddlewareFunc
	worker      worker.Worker
	server      *http.Server
}

//New creates a new Engine
func New(settings *models.SystemSettings) *Engine {
	db := dbx.New()
	logger := database.NewLogger("WEB", db)
	logger.SetProperty(log.PropertyKeyContextID, rand.String(32))

	bgLogger := database.NewLogger("BGW", db)
	bgLogger.SetProperty(log.PropertyKeyContextID, rand.String(32))

	router := &Engine{
		mux:         httprouter.New(),
		db:          db,
		logger:      logger,
		renderer:    NewRenderer(settings, logger),
		binder:      NewDefaultBinder(),
		middlewares: make([]MiddlewareFunc, 0),
		worker:      worker.New(db, bgLogger),
	}

	return router
}

//Start the server.
func (e *Engine) Start(address string) {
	e.logger.Info("Application is starting")
	e.logger.Infof("GO_ENV: @{Env}", log.Props{
		"Env": env.Current(),
	})

	var (
		autoSSL      = env.GetEnvOrDefault("SSL_AUTO", "")
		certFilePath = ""
		keyFilePath  = ""
	)

	if env.IsDefined("SSL_CERT") {
		certFilePath = env.Etc(env.GetEnvOrDefault("SSL_CERT", ""))
		keyFilePath = env.Etc(env.GetEnvOrDefault("SSL_CERT_KEY", ""))
	}

	e.server = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         address,
		Handler:      e.mux,
		ErrorLog:     stdLog.New(e.logger, "", 0),
		TLSConfig:    getDefaultTLSConfig(),
	}

	for i := 0; i < runtime.NumCPU()*2; i++ {
		go e.Worker().Run(strconv.Itoa(i))
	}

	var (
		err         error
		certManager *CertificateManager
	)
	if autoSSL == "true" {
		certManager, err = NewCertificateManager(certFilePath, keyFilePath, e.db.Connection())
		if err != nil {
			panic(errors.Wrap(err, "failed to initialize CertificateManager"))
		}

		e.server.TLSConfig.GetCertificate = certManager.GetCertificate
		e.logger.Infof("https (auto ssl) server started on @{Address}", log.Props{"Address": address})
		go certManager.StartHTTPServer()
		err = e.server.ListenAndServeTLS("", "")
	} else if certFilePath == "" && keyFilePath == "" {
		e.logger.Infof("http server started on @{Address}", log.Props{"Address": address})
		err = e.server.ListenAndServe()
	} else {
		e.logger.Infof("https server started on @{Address}", log.Props{"Address": address})
		err = e.server.ListenAndServeTLS(certFilePath, keyFilePath)
	}

	if err != nil && err != http.ErrServerClosed {
		panic(errors.Wrap(err, "failed to start server"))
	}
}

//Stop the server.
func (e *Engine) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	e.logger.Info("server is shutting down")
	if err := e.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}
	e.logger.Info("server has shutdown")

	e.logger.Info("worker is shutting down")
	if err := e.worker.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown worker")
	}
	e.logger.Info("worker has shutdown")

	return nil
}

//NewContext creates and return a new context
func (e *Engine) NewContext(res http.ResponseWriter, req *http.Request, params StringMap) Context {
	contextID := rand.String(32)
	request := WrapRequest(req)
	ctxLogger := e.logger.New()
	ctxLogger.SetProperty(log.PropertyKeyContextID, contextID)
	ctxLogger.SetProperty("UserAgent", req.Header.Get("User-Agent"))

	return Context{
		id:       contextID,
		Response: res,
		Request:  request,
		engine:   e,
		logger:   ctxLogger,
		params:   params,
		worker:   e.worker,
	}
}

//Logger returns current logger
func (e *Engine) Logger() log.Logger {
	return e.logger
}

//Database returns current database
func (e *Engine) Database() *dbx.Database {
	return e.db
}

//Worker returns current worker reference
func (e *Engine) Worker() worker.Worker {
	return e.worker
}

//Group creates a new route group
func (e *Engine) Group() *Group {
	g := &Group{
		engine:      e,
		middlewares: e.middlewares,
	}
	return g
}

//Use adds a middleware to the root engine
func (e *Engine) Use(middleware MiddlewareFunc) {
	e.middlewares = append(e.middlewares, middleware)
}

//Get handles HTTP GET requests
func (e *Engine) Get(path string, handler HandlerFunc) {
	e.mux.Handle("GET", path, e.handle(e.middlewares, handler))
}

//Post handles HTTP POST requests
func (e *Engine) Post(path string, handler HandlerFunc) {
	e.mux.Handle("POST", path, e.handle(e.middlewares, handler))
}

//Put handles HTTP PUT requests
func (e *Engine) Put(path string, handler HandlerFunc) {
	e.mux.Handle("PUT", path, e.handle(e.middlewares, handler))
}

//Delete handles HTTP DELETE requests
func (e *Engine) Delete(path string, handler HandlerFunc) {
	e.mux.Handle("DELETE", path, e.handle(e.middlewares, handler))
}

//NotFound register how to handle routes that are not found
func (e *Engine) NotFound(handler HandlerFunc) {
	e.mux.NotFound = &notFoundHandler{
		engine:  e,
		handler: handler,
	}
}

func (e *Engine) handle(middlewares []MiddlewareFunc, handler HandlerFunc) httprouter.Handle {
	next := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	var h = func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		params := make(StringMap, 0)
		for _, p := range ps {
			params[p.Key] = p.Value
		}
		ctx := e.NewContext(res, req, params)
		next(ctx)
	}
	return h
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
	g.engine.mux.Handle("GET", path, g.engine.handle(g.middlewares, handler))
}

//Post handles HTTP POST requests
func (g *Group) Post(path string, handler HandlerFunc) {
	g.engine.mux.Handle("POST", path, g.engine.handle(g.middlewares, handler))
}

//Put handles HTTP PUT requests
func (g *Group) Put(path string, handler HandlerFunc) {
	g.engine.mux.Handle("PUT", path, g.engine.handle(g.middlewares, handler))
}

//Delete handles HTTP DELETE requests
func (g *Group) Delete(path string, handler HandlerFunc) {
	g.engine.mux.Handle("DELETE", path, g.engine.handle(g.middlewares, handler))
}

// Static return files from given folder
func (g *Group) Static(prefix, root string) {
	fi, err := os.Stat(env.Path(root))
	if err != nil {
		panic(fmt.Sprintf("Path '%s' not found", root))
	}

	var h HandlerFunc
	if fi.IsDir() {
		h = func(c Context) error {
			path := root + c.Param("filepath")
			fi, err := os.Stat(path)
			if err == nil && !fi.IsDir() {
				http.ServeFile(c.Response, c.Request.instance, path)
				return nil
			}
			return c.NotFound()
		}
	} else {
		h = func(c Context) error {
			http.ServeFile(c.Response, c.Request.instance, root)
			return nil
		}
	}
	g.engine.mux.Handle("GET", prefix, g.engine.handle(g.middlewares, h))
}

// ParseCookie return a list of cookie parsed from raw Set-Cookie
func ParseCookie(s string) *http.Cookie {
	if s == "" {
		return nil
	}
	return (&http.Response{Header: http.Header{"Set-Cookie": {s}}}).Cookies()[0]
}
