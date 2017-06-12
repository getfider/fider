package handlers

import (
	"net/http"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/labstack/echo"
)

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
