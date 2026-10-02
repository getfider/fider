package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// permissionsPolicy disables browser features that Fider never uses
const permissionsPolicy = "camera=(), microphone=(), geolocation=(), payment=(), usb=()"

// Secure middleware is responsible for
// 1. Setting the HTTP Security Headers
// 2. Protecting from Host attacks
func Secure() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			cdnHost := env.Config.CDN.Host
			if cdnHost != "" {
				if !env.IsSingleHostMode() {
					cdnHost = "*." + cdnHost
				}
				cdnHost = " " + cdnHost
			}
			analyticsScript, analyticsConnect := "", ""
			if env.Config.GoogleAnalytics != "" {
				analyticsScript, analyticsConnect = web.CspGoogleAnalyticsScript, web.CspGoogleAnalyticsConnect
			}
			csp := fmt.Sprintf(web.CspPolicyTemplate, c.ContextID(), cdnHost, analyticsScript, analyticsConnect)

			c.Response.Header().Set("Content-Security-Policy", strings.TrimSpace(csp))
			c.Response.Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response.Header().Set("X-Content-Type-Options", "nosniff")
			c.Response.Header().Set("Referrer-Policy", "no-referrer-when-downgrade")
			c.Response.Header().Set("Permissions-Policy", permissionsPolicy)

			// Only when Fider terminates TLS itself. Behind a reverse proxy, HSTS is the operator's decision.
			if env.Config.TLS.Automatic || env.Config.TLS.Certificate != "" {
				c.Response.Header().Set("Strict-Transport-Security", "max-age=15552000")
			}
			return next(c)
		}
	}
}

// CSRF middleware is responsible for blocking Cross-Site Request Forgery attacks.
//
// A state-changing request (anything other than GET/HEAD/OPTIONS) is only
// allowed when it carries a signal that a cross-origin web page cannot produce
// without first passing a CORS preflight (which Fider never grants):
//
//   - a Content-Type whose media type is exactly "application/json"; or
//   - an "Authorization: Bearer ..." header on an /api/ path, i.e. exactly the
//     requests that the User middleware authenticates with an API key.
//
// The Accept header is deliberately NOT considered: it is CORS-safelisted, so a
// cross-origin page can send "Accept: application/json" in a no-cors request.
//
// As defence in depth, requests that the browser labels as
// "Sec-Fetch-Site: cross-site" are rejected. "same-site" is allowed so that
// multi-tenant subdomains keep working, and requests without the header
// (non-browser clients, older browsers) are unaffected.
func CSRF() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			if isWriteMethod(c.Request.Method) {
				if strings.EqualFold(c.Request.GetHeader("Sec-Fetch-Site"), "cross-site") {
					return c.Forbidden()
				}
				if !web.IsJSONContentType(c.Request.GetHeader("Content-Type")) && !isAPIKeyRequest(c) {
					return c.Forbidden()
				}
			}
			return next(c)
		}
	}
}

func isWriteMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	default:
		return true
	}
}

// isAPIKeyRequest returns true if the request would be authenticated with an API key
// by the User middleware: an /api/ path with a Bearer token. It is scoped this way so
// the CSRF exemption matches exactly what the API key authentication path accepts.
func isAPIKeyRequest(c *web.Context) bool {
	if !c.Request.IsAPI() {
		return false
	}
	_, ok := web.ParseBearerToken(c.Request.GetHeader("Authorization"))
	return ok
}
