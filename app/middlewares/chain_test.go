package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestChain_ExecutedCorrectOrder(t *testing.T) {
	RegisterT(t)

	mw1 := func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			c.Response.Header().Add("Key1", "Value1")
			c.Response.Header().Add("Key2", "Value2")
			return next(c)
		}
	}

	mw2 := func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			c.Response.Header().Del("Key2")
			return next(c)
		}
	}

	server := mock.NewServer()
	server.Use(middlewares.Chain(mw1, mw2))
	handler := func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	}

	status, response := server.Execute(handler)

	Expect(status).Equals(http.StatusOK)
	Expect(response.Header().Get("Key1")).Equals("Value1")
	Expect(response.Header().Get("Key2")).Equals("")
}

func TestChain_NilMiddleware(t *testing.T) {
	RegisterT(t)

	mw1 := func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			c.Response.Header().Add("Key1", "Value1")
			c.Response.Header().Add("Key2", "Value2")
			return next(c)
		}
	}

	server := mock.NewServer()
	server.Use(middlewares.Chain(mw1, nil))
	handler := func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	}

	status, response := server.Execute(handler)

	Expect(status).Equals(http.StatusOK)
	Expect(response.Header().Get("Key1")).Equals("Value1")
	Expect(response.Header().Get("Key2")).Equals("Value2")
}
