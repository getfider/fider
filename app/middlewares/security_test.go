package middlewares_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSecure(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.Secure())

	var id string
	status, response := server.Execute(func(c web.Context) error {
		id = c.ContextID()
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Header().Get("Content-Security-Policy")).Equals(fmt.Sprintf(web.CspPolicyTemplate, id))
	Expect(response.Header().Get("X-XSS-Protection")).Equals("1; mode=block")
	Expect(response.Header().Get("X-Content-Type-Options")).Equals("nosniff")
	Expect(response.Header().Get("Referrer-Policy")).Equals("no-referrer-when-downgrade")
}
