package web

import (
	"context"
	"crypto/tls"
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/uuid"
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
	cspConnect = "connect-src 'self' https://www.google-analytics.com"

	//CspPolicyTemplate is the template used to generate the policy
	CspPolicyTemplate = fmt.Sprintf("%s; %s; %s; %s; %s; %s; %s; %s; %s", cspBase, cspDefault, cspStyle, cspScript, cspImage, cspFont, cspObject, cspMedia, cspConnect)
)

type notFoundHandler struct {
	router *Engine
}

func (h *notFoundHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := h.router.NewContext(res, req, nil)
	ctx.NotFound()
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
	binder      *DefaultBinder
	middlewares []MiddlewareFunc
	worker      worker.Worker
	server      *http.Server
}

//New creates a new Engine
func New(settings *models.SystemSettings) *Engine {
	logger := log.NewConsoleLogger("WEB")
	router := &Engine{
		mux:         httprouter.New(),
		logger:      logger,
		renderer:    NewRenderer(settings, logger),
		binder:      NewDefaultBinder(),
		middlewares: make([]MiddlewareFunc, 0),
		worker:      worker.New(),
	}

	router.mux.NotFound = &notFoundHandler{router}
	return router
}

//Start the server.
func (e *Engine) Start(address string) {
	certFile := env.GetEnvOrDefault("SSL_CERT", "")
	keyFile := env.GetEnvOrDefault("SSL_CERT_KEY", "")
	autoSSL := env.GetEnvOrDefault("SSL_AUTO", "")

	certFilePath := env.Etc(certFile)
	keyFilePath := env.Etc(keyFile)

	e.server = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         address,
		Handler:      e.mux,
		ErrorLog:     stdLog.New(e.logger, "", 0),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			MaxVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			},
		},
	}

	for i := 0; i < runtime.NumCPU()*2; i++ {
		go e.Worker().Run(strconv.Itoa(i))
	}

	var (
		err         error
		certManager *CertificateManager
	)
	if autoSSL == "true" {
		certManager, err = NewCertificateManager(certFilePath, keyFilePath, env.Etc("certs"))
		if err != nil {
			panic(errors.Wrap(err, "failed to initialize CertificateManager"))
		}

		e.server.TLSConfig.GetCertificate = certManager.GetCertificate
		e.logger.Infof("https (auto ssl) server started on %s", address)
		go certManager.StartHTTPServer()
		err = e.server.ListenAndServeTLS("", "")
	} else if certFile == "" && keyFile == "" {
		e.logger.Infof("http server started on %s", address)
		err = e.server.ListenAndServe()
	} else {
		e.logger.Infof("https server started on %s", address)
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

	e.logger.Infof("server is shutting down")
	if err := e.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}
	e.logger.Infof("server has shutdown")

	e.logger.Infof("worker is shutting down")
	if err := e.worker.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown worker")
	}
	e.logger.Infof("worker has shutdown")

	return nil
}

//NewContext creates and return a new context
func (e *Engine) NewContext(res http.ResponseWriter, req *http.Request, params StringMap) Context {
	return Context{
		id:       strings.Replace(uuid.NewV4().String(), "-", "", 4),
		Response: res,
		Request:  WrapRequest(req),
		engine:   e,
		logger:   e.logger,
		params:   params,
		worker:   e.worker,
	}
}

//Logger returns current logger
func (e *Engine) Logger() log.Logger {
	return e.logger
}

//Worker returns current worker referenc
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

//Delete handles HTTP DELETE requests
func (g *Group) Delete(path string, handler HandlerFunc) {
	g.engine.mux.Handle("DELETE", path, g.handler(handler))
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
		next(ctx)
	}
	return h
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
	g.engine.mux.Handle("GET", prefix, g.handler(h))
}

// ParseCookies return a list of cookie parsed from raw Set-Cookie
func ParseCookies(s string) []*http.Cookie {
	return (&http.Response{Header: http.Header{"Set-Cookie": {s}}}).Cookies()
}
