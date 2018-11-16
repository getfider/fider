package mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/jsonq"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
)

// Server is a HTTP server wrapper for testing purpose
type Server struct {
	engine     *web.Engine
	context    web.Context
	recorder   *httptest.ResponseRecorder
	middleware []web.MiddlewareFunc
}

func createServer(services *app.Services) *Server {
	settings := &models.SystemSettings{}
	engine := web.New(settings)

	request, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	context := engine.NewContext(recorder, request, make(web.StringMap, 0))
	context.SetServices(services)

	return &Server{
		engine:     engine,
		recorder:   recorder,
		context:    context,
		middleware: []web.MiddlewareFunc{},
	}
}

// Engine returns current engine from mocked server
func (s *Server) Engine() *web.Engine {
	return s.engine
}

// Use adds a new middleware to pipeline
func (s *Server) Use(middleware web.MiddlewareFunc) *Server {
	s.middleware = append(s.middleware, middleware)
	return s
}

// OnTenant set current context tenant
func (s *Server) OnTenant(tenant *models.Tenant) *Server {
	s.context.SetTenant(tenant)
	s.context.Services().SetCurrentTenant(tenant)
	return s
}

// AsUser set current context user
func (s *Server) AsUser(user *models.User) *Server {
	s.context.SetUser(user)
	s.context.Services().SetCurrentUser(user)
	return s
}

// AddParam to current context route parameters
func (s *Server) AddParam(name string, value interface{}) *Server {
	s.context.AddParam(name, fmt.Sprintf("%v", value))
	return s
}

// AddHeader add key-value to current context headers
func (s *Server) AddHeader(name string, value string) *Server {
	s.context.Request.SetHeader(name, value)
	return s
}

// AddCookie add key-value to current context cookies
func (s *Server) AddCookie(name string, value string) *Server {
	s.context.Request.AddCookie(&http.Cookie{Name: name, Value: value})
	return s
}

// WithClientIP set current ClientIP address
func (s *Server) WithClientIP(clientIP string) *Server {
	s.context.Request.ClientIP = clientIP
	return s
}

// WithClaims set current context user claims
func (s *Server) WithClaims(claims *jwt.FiderClaims) *Server {
	s.context.SetClaims(claims)
	return s
}

// WithURL set current context Request URL
func (s *Server) WithURL(fullURL string) *Server {
	u, _ := url.Parse(fullURL)
	s.context.Request.URL = u
	return s
}

// Execute given handler and return response
func (s *Server) Execute(handler web.HandlerFunc) (int, *httptest.ResponseRecorder) {
	next := handler
	for i := len(s.middleware) - 1; i >= 0; i-- {
		next = s.middleware[i](next)
	}

	if err := next(s.context); err != nil {
		s.context.Failure(err)
	}

	return s.recorder.Code, s.recorder
}

// ExecuteAsJSON given handler and return json response
func (s *Server) ExecuteAsJSON(handler web.HandlerFunc) (int, *jsonq.Query) {
	code, response := s.Execute(handler)
	return code, toJSONQuery(response)
}

// ExecutePost executes given handler as POST and return response
func (s *Server) ExecutePost(handler web.HandlerFunc, body string) (int, *httptest.ResponseRecorder) {
	s.context.Request.Method = "POST"
	s.context.Request.Body = body
	s.context.Request.ContentLength = int64(len(body))
	s.context.Request.SetHeader("Content-Type", web.UTF8JSONContentType)

	return s.Execute(handler)
}

// ExecutePostAsJSON executes given handler as POST and return json response
func (s *Server) ExecutePostAsJSON(handler web.HandlerFunc, body string) (int, *jsonq.Query) {
	code, response := s.ExecutePost(handler, body)
	return code, toJSONQuery(response)
}

func toJSONQuery(response *httptest.ResponseRecorder) *jsonq.Query {
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return jsonq.New(string(b))
}
