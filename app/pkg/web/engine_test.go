package web_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app/pkg/web"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestEngine_StartRequestStop(t *testing.T) {
	RegisterT(t)
	w := web.New(&models.SystemSettings{})
	group := w.Group()
	{
		group.Get("/hello", func(c web.Context) error {
			return c.Ok(web.Map{})
		})
	}

	go w.Start(":8080")
	time.Sleep(time.Second)
	resp, err := http.Get("http://127.0.0.1:8080/hello")
	Expect(err).IsNil()
	Expect(resp.StatusCode).Equals(http.StatusOK)
	resp.Body.Close()

	resp, err = http.Get("http://127.0.0.1:8080/world")
	Expect(err).IsNil()
	Expect(resp.StatusCode).Equals(http.StatusNotFound)
	resp.Body.Close()

	err = w.Stop()
	Expect(err).IsNil()

	resp, err = http.Get("http://127.0.0.1:8080/hello")
	Expect(err).IsNotNil()
	Expect(resp).IsNil()
}
