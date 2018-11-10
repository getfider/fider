package middlewares

import (
	"time"

	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/web"
)

// Session starts a new Session if an Session ID is not yet set
func Session() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			cookie, err := c.Request.Cookie(web.CookieSessionName)
			if err != nil {
				nextYear := time.Now().Add(365 * 24 * time.Hour)
				cookie = c.AddCookie(web.CookieSessionName, rand.String(48), nextYear)
			}
			c.SetSessionID(cookie.Value)
			return next(c)
		}
	}
}
