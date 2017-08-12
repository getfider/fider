package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/jmoiron/jsonq"
)

// Server is a HTTP server wrapper for testing purpose
type Server struct {
	engine     *web.Engine
	Context    web.Context
	recorder   *httptest.ResponseRecorder
	middleware web.MiddlewareFunc
}

// NewSingleTenantServer creates a new multitenant test server
func NewSingleTenantServer() *Server {
	server := NewServer()
	os.Setenv("HOST_MODE", "single")
	return server
}

// NewServer creates a new test server
func NewServer() *Server {
	os.Setenv("HOST_MODE", "multi")
	settings := &models.AppSettings{}
	engine := web.New(settings)

	request, _ := http.NewRequest("GET", "/", nil)
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

// OnTenant set current context tenant
func (s *Server) OnTenant(tenant *models.Tenant) *Server {
	s.Context.SetTenant(tenant)
	return s
}

// AsUser set current context user
func (s *Server) AsUser(user *models.User) *Server {
	s.Context.SetUser(user)
	return s
}

// WithParam set current context params
func (s *Server) WithParam(name string, value interface{}) *Server {
	s.Context.SetParams(web.Map{name: value})
	return s
}

// WithParams set current context params
func (s *Server) WithParams(params web.Map) *Server {
	s.Context.SetParams(params)
	return s
}

// Execute given handler and return response
func (s *Server) Execute(handler web.HandlerFunc) (int, *httptest.ResponseRecorder) {

	if err := s.middleware(handler)(s.Context); err != nil {
		s.engine.HandleError(err, s.Context)
	}

	return s.recorder.Code, s.recorder
}

// ExecuteAsJSON given handler and return json response
func (s *Server) ExecuteAsJSON(handler web.HandlerFunc) (int, *jsonq.JsonQuery) {

	if err := s.middleware(handler)(s.Context); err != nil {
		s.engine.HandleError(err, s.Context)
	}

	return parseJSONBody(s)
}

// ExecutePost executes given handler as POST and return response
func (s *Server) ExecutePost(handler web.HandlerFunc, body string) (int, *httptest.ResponseRecorder) {

	req, _ := http.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	s.Context.SetRequest(req)

	if err := s.middleware(handler)(s.Context); err != nil {
		s.engine.HandleError(err, s.Context)
	}

	return s.recorder.Code, s.recorder

}

func parseJSONBody(s *Server) (int, *jsonq.JsonQuery) {

	if s.recorder.Code == 200 && hasJSON(s.recorder) {
		var data interface{}
		decoder := json.NewDecoder(s.recorder.Body)
		decoder.Decode(&data)
		query := jsonq.NewQuery(data)
		return s.recorder.Code, query
	}

	return s.recorder.Code, nil
}

func hasJSON(r *httptest.ResponseRecorder) bool {
	isJSONContentType := strings.Contains(r.Result().Header.Get("Content-Type"), "application/json")

	if r.Body.Len() > 0 && isJSONContentType {
		return true
	}

	return false
}
