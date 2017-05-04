package mock2

import (
	"net/http"
	"net/http/httptest"

	"github.com/WeCanHearYou/wechy/app/middlewares"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/labstack/echo"
)

// Server is a HTTP server wrapper for testing purpose
type Server struct {
	engine     *web.Engine
	handler    echo.HandlerFunc
	Context    web.Context
	recorder   *httptest.ResponseRecorder
	middleware web.MiddlewareFunc
}

// NewServer creates a new test server
func NewServer() *Server {
	settings := &models.AppSettings{}
	engine := web.New(settings)

	request, _ := http.NewRequest(echo.GET, "/", nil)
	recorder := httptest.NewRecorder()
	context := engine.NewContext(request, recorder)

	return &Server{
		engine:     engine,
		recorder:   recorder,
		Context:    web.Context{Context: context},
		middleware: middlewares.Noop(),
	}
}

// Use adds a new middleware to pipeline
func (s *Server) Use(middleware web.MiddlewareFunc) {
	s.middleware = middleware
}

// Execute given handler and return response
func (s *Server) Execute(handler web.HandlerFunc) (int, *httptest.ResponseRecorder) {
	if err := s.middleware(handler)(s.Context); err != nil {
		s.engine.HandleError(err, s.Context)
	}

	return s.recorder.Code, s.recorder
}
