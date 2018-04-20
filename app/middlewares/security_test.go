package middlewares_test

import (
	"fmt"
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

	var id string
	status, response := server.Execute(func(c web.Context) error {
		id = c.ContextID()
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Content-Security-Policy-Report-Only")).To(Equal(fmt.Sprintf(web.CspPolicyTemplate, id)))
	Expect(response.Header().Get("X-XSS-Protection")).To(Equal("1; mode=block"))
	Expect(response.Header().Get("X-Content-Type-Options")).To(Equal("nosniff"))
	Expect(response.Header().Get("Referrer-Policy")).To(Equal("no-referrer-when-downgrade"))
}
