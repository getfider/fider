package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/jsonq"
)

type testServer struct {
	engine *gin.Engine
}

func NewTestServer() *testServer {
	engine := gin.New()
	return &testServer{
		engine: engine,
	}
}

func (s *testServer) do(f gin.HandlerFunc) {
	s.engine.Use(f)
}

func (s *testServer) param(key, value string) {
	s.do(func(c *gin.Context) {
		c.Params = append(c.Params, gin.Param{Key: key, Value: value})
	})
}

func (s *testServer) register(handler gin.HandlerFunc) {
	s.engine.Handle("GET", "/", handler)
}

func (s *testServer) request() (int, *jsonq.JsonQuery) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	s.engine.ServeHTTP(response, request)

	status := response.Result().StatusCode

	if status == 200 {
		var data interface{}
		decoder := json.NewDecoder(response.Body)
		decoder.Decode(&data)
		query := jsonq.NewQuery(data)
		return status, query
	}

	return status, nil
}
