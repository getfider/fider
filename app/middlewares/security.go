package middlewares

import (
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

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
			csp := fmt.Sprintf(web.CspPolicyTemplate, c.ContextID(), cdnHost)

			c.Response.Header().Set("Content-Security-Policy", strings.TrimSpace(csp))
			c.Response.Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response.Header().Set("X-Content-Type-Options", "nosniff")
			c.Response.Header().Set("Referrer-Policy", "no-referrer-when-downgrade")
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
//   - an "Authorization: Bearer ..." header (API key clients).
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
				if !hasJSONContentType(c) && !hasBearerAuthorization(c) {
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

// hasJSONContentType returns true only if the Content-Type media type is exactly
// application/json. Substring matching is unsafe, as "text/plain; x=application/json"
// is a CORS-safelisted Content-Type.
func hasJSONContentType(c *web.Context) bool {
	contentType := c.Request.GetHeader("Content-Type")
	if contentType == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return mediaType == web.JSONContentType
}

// hasBearerAuthorization returns true if the request has an Authorization header
// using the Bearer scheme. Browsers never attach Bearer credentials automatically
// (unlike cached Basic credentials), and cross-origin pages cannot set this header
// without a CORS preflight.
func hasBearerAuthorization(c *web.Context) bool {
	auth := strings.TrimSpace(c.Request.GetHeader("Authorization"))
	return len(auth) > len("Bearer ") && strings.EqualFold(auth[:len("Bearer ")], "Bearer ")
}
