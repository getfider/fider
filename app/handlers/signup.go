package handlers

import (
	"net/http"

	"strings"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/validate"
	"github.com/labstack/echo"
)

// CheckAvailability checks if given domain is available to be used
func CheckAvailability() web.HandlerFunc {
	return func(c web.Context) error {
		subdomain := c.Param("subdomain")
		ok, messages := validate.Subdomain(subdomain)
		if !ok {
			return c.BadRequest(echo.Map{
				"message": strings.Join(messages, ","),
			})
		}

		available, err := c.Services().Tenants.IsSubdomainAvailable(subdomain)
		if err != nil {
			return c.Failure(err)
		}

		if !available {
			return c.BadRequest(echo.Map{
				"message": "This subodmain is not available anymore",
			})
		}

		return c.Ok(echo.Map{})
	}
}

//SignUp is the entry point for installation / signup
func SignUp() web.HandlerFunc {
	return func(c web.Context) error {
		if env.IsSingleHostMode() {
			tenant, err := c.Services().Tenants.First()
			if err != nil {
				return c.Failure(err)
			}

			if tenant != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}
		}
		return c.Page(echo.Map{})
	}
}
