package infra

import (
	"net/http"
	"net/url"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
)

// IsAuthenticated blocks non-authenticated requests
func IsAuthenticated() app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(c app.Context) error {
			if c.User() == nil {
				return c.NoContent(http.StatusForbidden)
			}
			return next(c)
		}
	}
}

// IsAuthorized blocks non-authorized requests
func IsAuthorized(roles ...models.Role) app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(c app.Context) error {
			user := c.User()
			for _, role := range roles {
				if user.Role == role {
					return next(c)
				}
			}
			return c.NoContent(http.StatusForbidden)
		}
	}
}

// HostChecker checks for a specific host
func HostChecker(baseURL string) app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(c app.Context) error {
			u, err := url.Parse(baseURL)
			if err != nil {
				return c.Failure(err)
			}

			if c.Request().Host != u.Host {
				c.Logger().Errorf("%s is not valid for this operation. Only %s is allowed.", c.Request().Host, u.Host)
				return c.NoContent(http.StatusBadRequest)
			}

			return next(c)
		}
	}
}
