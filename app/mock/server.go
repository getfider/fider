package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"strings"

	"fmt"

	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/jmoiron/jsonq"
	"github.com/labstack/echo"
)

// Server is a HTTP server wrapper for testing purpose
type Server struct {
	engine   *echo.Echo
	handler  echo.HandlerFunc
	Context  web.Context
	recorder *httptest.ResponseRecorder
}

// NewServer creates a new test server
func NewServer() *Server {
	engine := echo.New()
	engine.Renderer = web.NewHTMLRenderer(engine.Logger)
	engine.HTTPErrorHandler = func(e error, c echo.Context) {
		fmt.Println(e)
		c.NoContent(e.(*echo.HTTPError).Code)
	}

	request, _ := http.NewRequest(echo.GET, "/", nil)
	recorder := httptest.NewRecorder()
	context := engine.NewContext(request, recorder)

	return &Server{
		engine:   engine,
		recorder: recorder,
		Context:  web.Context{Context: context},
	}
}

// NewContext creates a new HTTP context
func (s *Server) NewContext(req *http.Request, w http.ResponseWriter) web.Context {
	c := s.engine.NewContext(req, w)
	return web.Context{Context: c}
}

// Execute given handler and return response as JSON
func (s *Server) Execute(handler web.HandlerFunc) (int, *jsonq.JsonQuery) {
	ctx := web.Context{Context: s.Context}
	if err := handler(ctx); err != nil {
		s.engine.HTTPErrorHandler(err, ctx)
	}

	return parseJSONBody(s)
}

// ExecutePost executes given handler with posted JSON
func (s *Server) ExecutePost(handler web.HandlerFunc, body string) (int, *jsonq.JsonQuery) {
	ctx := web.Context{Context: s.Context}

	req, _ := http.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx.SetRequest(req)

	if err := handler(ctx); err != nil {
		s.engine.HTTPErrorHandler(err, ctx)
	}
	return parseJSONBody(s)
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

// ExecuteRaw executes given handler and return raw response
func (s *Server) ExecuteRaw(handler web.HandlerFunc) (int, *http.Response) {
	ctx := web.Context{Context: s.Context}
	if err := handler(ctx); err != nil {
		s.engine.HTTPErrorHandler(err, ctx)
	}

	return s.recorder.Code, s.recorder.Result()
}

func hasJSON(r *httptest.ResponseRecorder) bool {
	isJSONContentType := strings.Contains(r.Result().Header.Get("Content-Type"), "application/json")

	if r.Body.Len() > 0 && isJSONContentType {
		return true
	}

	return false
}
