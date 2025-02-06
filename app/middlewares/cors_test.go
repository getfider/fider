package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/Spicy-Bush/fider-tarkov-community/app/middlewares"
	. "github.com/Spicy-Bush/fider-tarkov-community/app/pkg/assert"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/mock"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
)

func TestCORS(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.CORS())
	handler := func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	}

	status, response := server.Execute(handler)

	Expect(status).Equals(http.StatusOK)
	Expect(response.Header().Get("Access-Control-Allow-Origin")).Equals("*")
	Expect(response.Header().Get("Access-Control-Allow-Methods")).Equals("GET")
}
