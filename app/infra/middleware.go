package infra

import (
	"net/http"
	"net/url"

	"github.com/WeCanHearYou/wechy/app"
)

// IsAuthenticated blocks non-authenticated requests
func IsAuthenticated() app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(c app.Context) error {
			claims := c.Claims()
			if claims == nil {
				return c.NoContent(http.StatusForbidden)
			}
			return next(c)
		}
	}
}

// JwtGetter gets JWT token from cookie and add into context
func JwtGetter() app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(c app.Context) error {

			if cookie, err := c.Cookie("auth"); err == nil {
				if claims, err := Decode(cookie.Value); err == nil {
					if claims.TenantID == c.Tenant().ID {
						c.SetClaims(claims)
					}
				} else {
					c.Logger().Error(err)
				}
			}

			return next(c)
		}
	}
}

// JwtSetter sets JWT token into cookie
func JwtSetter() app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(c app.Context) error {

			query := c.Request().URL.Query()

			jwt := query.Get("jwt")
			if jwt != "" {
				c.SetCookie(&http.Cookie{
					Name:     "auth",
					Value:    jwt,
					HttpOnly: true,
					Path:     "/",
				})

				scheme := "http"
				if c.Request().TLS != nil {
					scheme = "https"
				}

				query.Del("jwt")

				url := scheme + "://" + c.Request().Host + c.Request().URL.Path
				querystring := query.Encode()
				if querystring != "" {
					url += "?" + querystring
				}

				return c.Redirect(http.StatusTemporaryRedirect, url)
			}

			return next(c)
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
