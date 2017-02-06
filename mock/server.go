package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"fmt"

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
	fmt.Println(status)

	if status == 200 && response.Body.Len() > 0 {
		var data interface{}
		decoder := json.NewDecoder(response.Body)
		decoder.Decode(&data)
		query := jsonq.NewQuery(data)
		return status, query
	}

	return status, nil
}
