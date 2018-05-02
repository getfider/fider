package web_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/pkg/web"

	. "github.com/onsi/gomega"
)

func TestEngine_StartRequestStop(t *testing.T) {
	RegisterTestingT(t)
	w := web.New(nil)
	group := w.Group()
	{
		group.Get("/hello", func(c web.Context) error {
			return c.Ok(web.Map{})
		})
	}
	go w.Start(":8080")
	resp, err := http.Get("http://localhost:8080/hello")
	Expect(err).To(BeNil())
	resp.Body.Close()
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	err = w.Stop()
	Expect(err).To(BeNil())

	resp, err = http.Get("http://localhost:8080/hello")
	Expect(err).NotTo(BeNil())
	Expect(resp).To(BeNil())
}
