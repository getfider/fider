package router

import (
	"net/http"
	"strings"

	"fmt"

	"github.com/WeCanHearYou/wchy/auth"
	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/util"
	"github.com/labstack/echo"
)

// MultiTenant extract tenant information from hostname and inject it into current context
func MultiTenant(ctx *context.WchyContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hostname := util.StripPort(c.Request().Host)
			tenant, err := ctx.Tenant.GetByDomain(hostname)
			if err == nil {
				c.Set("Tenant", tenant)
				return next(c)
			}

			fmt.Printf("Tenant not found for %s.\n", hostname)
			return c.NoContent(http.StatusNotFound)
		}
	}
}

// HostChecker checks for a specific host
func HostChecker(baseURL string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			scheme := "http"
			if c.Request().TLS != nil {
				scheme = "https"
			}
			url := scheme + "://" + c.Request().Host

			if strings.Index(url, baseURL) < 0 {
				c.Logger().Errorf("%s is not valid for this operation", url)
				return c.NoContent(http.StatusBadRequest)
			}
			return next(c)
		}
	}
}

// JwtGetter gets JWT token from cookie and add into context
func JwtGetter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if cookie, err := c.Cookie("auth"); err == nil {
				if claims, err := auth.Decode(cookie.Value); err == nil {
					c.Set("Claims", claims)
				} else {
					c.Logger().Error(err)
				}
			}

			return next(c)
		}
	}
}

// JwtSetter sets JWT token into cookie
func JwtSetter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			query := c.Request().URL.Query()

			jwt := query.Get("jwt")
			if jwt != "" {
				c.SetCookie(&http.Cookie{
					Name:     "auth",
					Value:    jwt,
					HttpOnly: true,
				})

				scheme := "http"
				if c.Request().TLS != nil {
					scheme = "https"
				}

				query.Del("jwt")

				url := scheme + "://" + c.Request().Host + c.Request().URL.Path + "?" + query.Encode()
				return c.Redirect(http.StatusTemporaryRedirect, url)
			}

			return next(c)
		}
	}
}
