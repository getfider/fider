package web_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/pkg/web"

	. "github.com/getfider/fider/app/pkg/assert"
)

var e *web.Engine

func StartServer() {
	e = web.New()

	group := e.Group()
	{
		group.Get("/api/ping", func(c *web.Context) error {
			if c.Value("the-name") != nil {
				panic("key: the-name should not be set")
			}
			return c.String(http.StatusOK, "pong")
		})

		group.Use(func(next web.HandlerFunc) web.HandlerFunc {
			return func(c *web.Context) error {
				c.Set("the-name", c.QueryParam("name"))
				return next(c)
			}
		})

		group.Get("/api/echo", func(c *web.Context) error {
			return c.String(http.StatusOK, c.Value("the-name").(string))
		})
	}

	e.Get("/hello", func(c *web.Context) error {
		if c.Value("name") != nil {
			panic("ERROR!")
		}
		return c.String(http.StatusOK, "")
	})

	go e.Start(":8080")

	//Wait until web server is ready
	Expect(func() error {
		_, err := http.Get("http://127.0.0.1:8080/hello")
		return err
	}).EventuallyEquals(nil)

	//Wait until metrics server is ready
	Expect(func() error {
		_, err := http.Get("http://127.0.0.1:4000/metrics")
		return err
	}).EventuallyEquals(nil)
}

func StopServer() {
	err := e.Stop()
	Expect(err).IsNil()
}

func TestEngine_StartRequestStop(t *testing.T) {
	RegisterT(t)
	StartServer()

	resp, err := http.Get("http://127.0.0.1:8080/hello")
	Expect(err).IsNil()
	Expect(resp.StatusCode).Equals(http.StatusOK)
	resp.Body.Close()

	resp, err = http.Get("http://127.0.0.1:8080/world")
	Expect(err).IsNil()
	Expect(resp.StatusCode).Equals(http.StatusNotFound)
	resp.Body.Close()

	StopServer()

	resp, err = http.Get("http://127.0.0.1:8080/hello")
	Expect(err).IsNotNil()
	Expect(resp).IsNil()
}

func TestEngine_MiddlewareAfterHandler(t *testing.T) {
	RegisterT(t)
	StartServer()

	resp, err := http.Get("http://127.0.0.1:8080/api/ping")
	Expect(err).IsNil()
	Expect(resp.StatusCode).Equals(http.StatusOK)
	resp.Body.Close()

	resp, err = http.Get("http://127.0.0.1:8080/api/echo?name=John")
	Expect(err).IsNil()
	Expect(resp.StatusCode).Equals(http.StatusOK)
	content, err := ioutil.ReadAll(resp.Body)
	Expect(err).IsNil()
	Expect(string(content)).Equals("John")
	resp.Body.Close()

	StopServer()
}
