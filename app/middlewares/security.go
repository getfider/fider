package middlewares

import (
	"fmt"
	"net"
	"strings"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
)

// ForceHTTPS, when enabled, redirects HTTP requests to HTTPS
func ForceHTTPS() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			if !env.Config.TLS.Force {
				return next(c)
			}

			url := c.Request.URL
			isHttps := c.Request.URL.Scheme == "https" || c.Request.GetHeader("X-Forwarded-Proto") == "https"
			if isHttps {
				return next(c)
			}

			host, _, _ := net.SplitHostPort(url.Host)
			target := "https://" + host + url.RequestURI()
			return c.Redirect(target)
		}
	}
}

// Secure middleware is responsible for
// 1. Setting the HTTP Security Headers
// 2. Protecting from Host attacks
func Secure() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {

			// Host attack is only a problem if Fider is running on single host mode
			if env.IsSingleHostMode() {
				if c.Request.URL.Hostname() != env.Config.HostDomain {
					log.Errorf(c, "Requested hostname '@{URLHostname}' does not match environment HOST_DOMAIN '@{HostDomain}'.", dto.Props{
						"URLHostname": c.Request.URL.Hostname(),
						"HostDomain":  env.Config.HostDomain,
					})
					return c.NotFound()
				}
			}

			cdnHost := env.Config.CDN.Host
			if cdnHost != "" && !env.IsSingleHostMode() {
				cdnHost = "*." + cdnHost
			}
			csp := fmt.Sprintf(web.CspPolicyTemplate, c.ContextID(), cdnHost)

			c.Response.Header().Set("Content-Security-Policy", strings.TrimSpace(csp))
			c.Response.Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response.Header().Set("X-Content-Type-Options", "nosniff")
			c.Response.Header().Set("Referrer-Policy", "no-referrer-when-downgrade")
			return next(c)
		}
	}
}
