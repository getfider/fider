package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/WeCanHearYou/wechy/app/middlewares"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/jmoiron/jsonq"
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
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
