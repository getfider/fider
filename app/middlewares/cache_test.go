package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/mock"
	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

func TestCache(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Use(middlewares.OneYearCache())
	handler := func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	}

	status, response := server.Execute(handler)

	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Cache-Control")).To(Equal("public, max-age=30672000"))
}
