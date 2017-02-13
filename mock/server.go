package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/jmoiron/jsonq"
	"github.com/labstack/echo"
)

// Server is a HTTP server wrapper for testing purpose
type Server struct {
	engine  *echo.Echo
	handler echo.HandlerFunc
}

// NewServer creates a new test server
func NewServer() *Server {
	engine := echo.New()
	//TODO: engine.HTMLRender = router.CreateTemplateRender()

	return &Server{
		engine: engine,
	}
}

// Use to register new middlewares
func (s *Server) Use(mw echo.MiddlewareFunc) {
	s.engine.Use(mw)
}

// UseFunc to register new middlewares based on HandlerFunc
func (s *Server) UseFunc(f echo.HandlerFunc) {
	s.engine.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return f
	})
}

// Param can be used to set request parameters
func (s *Server) Param(key, value string) {
	s.UseFunc(func(c echo.Context) error {
		c.SetParamNames(key)
		c.SetParamValues(value)
		return nil
	})
}

// Set can be used to set context values
func (s *Server) Set(key string, value interface{}) {
	s.UseFunc(func(c echo.Context) error {
		c.Set(key, value)
		return nil
	})
}

// Register given handler
func (s *Server) Register(handler echo.HandlerFunc) {
	s.handler = handler
}

// Request registered handler
func (s *Server) Request() (int, *jsonq.JsonQuery) {
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	c := s.engine.NewContext(req, rec)
	s.handler(c)

	status := rec.Code

	if status == 200 && hasJSON(rec) {
		var data interface{}
		decoder := json.NewDecoder(rec.Body)
		decoder.Decode(&data)
		query := jsonq.NewQuery(data)
		return status, query
	}

	return status, nil
}

func hasJSON(r *httptest.ResponseRecorder) bool {
	isJSONContentType := strings.Contains(r.Result().Header.Get("Content-Type"), "application/json")

	if r.Body.Len() > 0 && isJSONContentType {
		return true
	}

	return false
}
