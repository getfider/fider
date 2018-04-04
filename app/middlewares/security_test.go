package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

func TestSecure(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.Secure())
	status, response := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Header().Get("X-XSS-Protection")).To(Equal("1; mode=block"))
	Expect(response.Header().Get("X-Content-Type-Options")).To(Equal("nosniff"))
	Expect(response.Header().Get("Referrer-Policy")).To(Equal("no-referrer-when-downgrade"))
}
