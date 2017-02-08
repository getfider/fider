package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/WeCanHearYou/wchy/router"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/jsonq"
)

// Server is a HTTP server wrapper for testing purpose
type Server struct {
	engine *gin.Engine
}

// NewServer creates a new test server
func NewServer() *Server {
	engine := gin.New()
	engine.HTMLRender = router.CreateTemplateRender()

	return &Server{
		engine: engine,
	}
}

// Use to register new middlewares
func (s *Server) Use(f gin.HandlerFunc) {
	s.engine.Use(f)
}

// Param can be used to set request parameters
func (s *Server) Param(key, value string) {
	s.Use(func(c *gin.Context) {
		c.Params = append(c.Params, gin.Param{Key: key, Value: value})
	})
}

// Set can be used to set context values
func (s *Server) Set(key string, value interface{}) {
	s.Use(func(c *gin.Context) {
		c.Set(key, value)
	})
}

// Register given handler
func (s *Server) Register(handler gin.HandlerFunc) {
	s.engine.Handle("GET", "/", handler)
}

// Request registered handler
func (s *Server) Request() (int, *jsonq.JsonQuery) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	s.engine.ServeHTTP(response, request)

	status := response.Result().StatusCode

	if status == 200 && hasJSON(response) {
		var data interface{}
		decoder := json.NewDecoder(response.Body)
		decoder.Decode(&data)
		query := jsonq.NewQuery(data)
		return status, query
	}

	return status, nil
}

func hasJSON(r *httptest.ResponseRecorder) bool {
	isJSONContentType := strings.Contains(r.Result().Header.Get("Content-Type"), gin.MIMEJSON)

	if r.Body.Len() > 0 && isJSONContentType {
		return true
	}

	return false
}
