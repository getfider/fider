package handlers

import (
	"net/http"
	"time"

	"github.com/getfider/fider/app/pkg/web"
)

// Logout remove auth cookies
func Logout() web.HandlerFunc {
	return func(c web.Context) error {
		c.SetCookie(&http.Cookie{
			Name:    "auth",
			MaxAge:  -1,
			Expires: time.Now().Add(-100 * time.Hour),
		})
		return c.Redirect(http.StatusTemporaryRedirect, c.QueryParam("redirect"))
	}
}
