package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/jmoiron/jsonq"
)

// Server is a HTTP server wrapper for testing purpose
type Server struct {
	engine     *web.Engine
	context    web.Context
	recorder   *httptest.ResponseRecorder
	middleware web.MiddlewareFunc
}

func createServer(services *app.Services) *Server {
	settings := &models.AppSettings{}
	engine := web.New(settings)

	request, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	context := engine.NewContext(request, recorder)
	context.SetServices(services)

	return &Server{
		engine:     engine,
		recorder:   recorder,
		context:    context,
		middleware: middlewares.Noop(),
	}
}

// Use adds a new middleware to pipeline
func (s *Server) Use(middleware web.MiddlewareFunc) {
	s.middleware = middleware
}

// OnTenant set current context tenant
func (s *Server) OnTenant(tenant *models.Tenant) *Server {
	s.context.SetTenant(tenant)
	return s
}

// AsUser set current context user
func (s *Server) AsUser(user *models.User) *Server {
	s.context.SetUser(user)
	return s
}

// WithParam set current context params
func (s *Server) WithParam(name string, value interface{}) *Server {
	s.context.SetParams(web.Map{name: value})
	return s
}

// AddHeader add key-value to current context headers
func (s *Server) AddHeader(name string, value string) *Server {
	s.context.Request().Header.Add(name, value)
	return s
}

// AddCookie add key-value to current context cookies
func (s *Server) AddCookie(name string, value string) *Server {
	s.context.Request().AddCookie(&http.Cookie{Name: name, Value: value})
	return s
}

// WithParams set current context params
func (s *Server) WithParams(params web.Map) *Server {
	s.context.SetParams(params)
	return s
}

// WithURL set current context Request URL
func (s *Server) WithURL(fullURL string) *Server {
	s.context.Request().URL, _ = url.Parse(fullURL)
	s.context.Request().Host = s.context.Request().URL.Host
	return s
}

// Execute given handler and return response
func (s *Server) Execute(handler web.HandlerFunc) (int, *httptest.ResponseRecorder) {

	if err := s.middleware(handler)(s.context); err != nil {
		s.engine.HandleError(err, s.context)
	}

	return s.recorder.Code, s.recorder
}

// ExecuteAsJSON given handler and return json response
func (s *Server) ExecuteAsJSON(handler web.HandlerFunc) (int, *jsonq.JsonQuery) {

	if err := s.middleware(handler)(s.context); err != nil {
		s.engine.HandleError(err, s.context)
	}

	return parseJSONBody(s)
}

// ExecutePost executes given handler as POST and return response
func (s *Server) ExecutePost(handler web.HandlerFunc, body string) (int, *httptest.ResponseRecorder) {

	req, _ := http.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	s.context.SetRequest(req)

	if err := s.middleware(handler)(s.context); err != nil {
		s.engine.HandleError(err, s.context)
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
