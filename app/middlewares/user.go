package middlewares

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app/pkg/validate"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
	webutil "github.com/getfider/fider/app/pkg/web/util"
)

// User gets JWT Auth token from cookie and insert into context
func User() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			var (
				token            string
				user             *entity.User
				fromSignUpCookie bool
				fromAPIKey       bool
			)

			cookie, err := c.Request.Cookie(web.CookieAuthName)
			if err == nil {
				token = cookie.Value
			} else {
				// The signup-transfer cookie is domain-wide, so it reaches every tenant
				// subdomain. We do NOT promote it to a durable host-only auth cookie here:
				// that only happens later, and only once we have confirmed the token's user
				// actually belongs to the tenant selected by the current host.
				token = webutil.GetSignUpAuthCookie(c)
				fromSignUpCookie = token != ""
			}

			if token != "" {
				claims, err := jwt.DecodeFiderClaims(token)
				if err != nil {
					c.RemoveCookie(web.CookieAuthName)
					return next(c)
				}

				// Scope the lookup to the tenant selected by the request host. A globally
				// valid user ID that belongs to a different tenant must resolve to "not
				// found" so a token cannot be replayed across tenant boundaries.
				tenantID := 0
				if c.Tenant() != nil {
					tenantID = c.Tenant().ID
				}
				userByClaimsID := &query.GetUserByID{UserID: claims.UserID, TenantID: tenantID}
				err = bus.Dispatch(c, userByClaimsID)
				user = userByClaimsID.Result
				if err != nil {
					if errors.Cause(err) == app.ErrNotFound {
						c.RemoveCookie(web.CookieAuthName)
						return next(c)
					}
					return err
				}

				// Security stamp check: if the JWT contains a stamp (new tokens only),
				// validate it against the current DB stamp. A mismatch means the user's
				// security-relevant data has changed (e.g. role changed, account blocked,
				// or OAuth allowed-roles updated) and they must re-authenticate so that
				// access controls are re-evaluated.
				if claims.SecurityStamp != "" && user != nil && claims.SecurityStamp != user.SecurityStamp {
					c.RemoveCookie(web.CookieAuthName)
					if c.IsAjax() {
						return c.JSON(401, web.Map{})
					}
					redirectTarget := c.Request.URL.RequestURI()
					if redirectTarget != "" &&
						redirectTarget != "/" &&
						!strings.HasPrefix(redirectTarget, "/signin") &&
						!strings.HasPrefix(redirectTarget, "/signout") {
						return c.Redirect("/signin?redirect=" + url.QueryEscape(redirectTarget))
					}
					return c.Redirect("/signin")
				}
			} else if c.Request.IsAPI() {
				if apiKey, ok := web.ParseBearerToken(c.Request.GetHeader("Authorization")); ok {
					getUserByAPIKey := &query.GetUserByAPIKey{APIKey: apiKey}
					err = bus.Dispatch(c, getUserByAPIKey)
					if err != nil {
						if errors.Cause(err) == app.ErrNotFound {
							return c.HandleValidation(validate.Failed("API Key is invalid"))
						}
						return err
					}
					user = getUserByAPIKey.Result
					fromAPIKey = true

					if !user.IsCollaborator() {
						return c.HandleValidation(validate.Failed("API Key is invalid"))
					}

					if impersonateUserIDStr := c.Request.GetHeader("X-Fider-UserID"); impersonateUserIDStr != "" {
						if !user.IsAdministrator() {
							return c.HandleValidation(validate.Failed("Only Administrators are allowed to impersonate another user"))
						}
						impersonateUserID, err := strconv.Atoi(impersonateUserIDStr)
						if err != nil {
							return c.HandleValidation(validate.Failed(fmt.Sprintf("User not found for given impersonate UserID '%s'", impersonateUserIDStr)))
						}
						userByImpersonateID := &query.GetUserByID{UserID: impersonateUserID, TenantID: user.Tenant.ID}
						err = bus.Dispatch(c, userByImpersonateID)
						user = userByImpersonateID.Result
						if err != nil {
							if errors.Cause(err) == app.ErrNotFound {
								return c.HandleValidation(validate.Failed(fmt.Sprintf("User not found for given impersonate UserID '%s'", impersonateUserIDStr)))
							}
							return err
						}
					}
				}
			}

			if user != nil && c.Tenant() != nil && user.Tenant.ID == c.Tenant().ID {
				// blocked users are unable to sign in
				if user.Status == enum.UserBlocked {
					c.RemoveCookie(web.CookieAuthName)
					return c.Unauthorized()
				}

				// Only now that the user is confirmed to belong to the current tenant do we
				// promote a signup-transfer cookie into a durable host-only auth cookie.
				// This is deliberately skipped on any other tenant, so the domain-wide
				// signup token cannot become a normal session on an unrelated subdomain.
				if fromSignUpCookie {
					webutil.AddAuthTokenCookie(c, token)
				}

				c.SetUser(user)
				if fromAPIKey {
					c.SetAuthenticatedByAPIKey()
				}
			}

			return next(c)
		}
	}
}
