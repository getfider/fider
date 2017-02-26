package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/WeCanHearYou/wchy/router"
	"github.com/jmoiron/jsonq"
	"github.com/labstack/echo"
)

// Server is a HTTP server wrapper for testing purpose
type Server struct {
	engine   *echo.Echo
	handler  echo.HandlerFunc
	Context  echo.Context
	recorder *httptest.ResponseRecorder
}

// NewServer creates a new test server
func NewServer() *Server {
	engine := echo.New()
	engine.Renderer = router.NewHTMLRenderer(engine.Logger)

	request, _ := http.NewRequest(echo.GET, "/", nil)
	recorder := httptest.NewRecorder()
	context := engine.NewContext(request, recorder)

	return &Server{
		engine:   engine,
		recorder: recorder,
		Context:  context,
	}
}

// Execute a given handler using current server setup
func (s *Server) Execute(handler echo.HandlerFunc) (int, *jsonq.JsonQuery) {
	handler(s.Context)

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
