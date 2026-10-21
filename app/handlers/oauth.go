package handlers

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	webutil "github.com/getfider/fider/app/pkg/web/util"
)

// handoffLifetime is how long the token that carries an OAuth sign in from the callback
// address over to the tenant address stays valid. It only has to survive one redirect.
const handoffLifetime = 2 * time.Minute

// tenantOrigins returns every origin a given tenant is served on.
// A tenant with a custom domain is still reachable on its subdomain, so both are valid.
func tenantOrigins(ctx context.Context, tenant *entity.Tenant) []string {
	origins := []string{web.TenantBaseURL(ctx, tenant)}

	if tenant.CNAME != "" && !env.IsSingleHostMode() {
		origins = append(origins, web.TenantBaseURL(ctx, &entity.Tenant{Subdomain: tenant.Subdomain}))
	}

	return origins
}

// oauthAllowedOrigins returns the origins an OAuth flow is allowed to continue on for the
// current request. These come from the resolved tenant record and from configuration, never
// from the request's Host/X-Forwarded-Host header, which any caller can set to anything.
func oauthAllowedOrigins(c *web.Context) []string {
	if c.Tenant() == nil {
		// Sign up runs on the host-wide OAuth address, before a tenant exists.
		return []string{web.OAuthBaseURL(c)}
	}

	return tenantOrigins(c, c.Tenant())
}

// isAllowedOAuthRedirect reports whether redirect points at one of the given origins.
func isAllowedOAuthRedirect(origins []string, redirect string) bool {
	for _, origin := range origins {
		if redirect == origin || strings.HasPrefix(redirect, origin+"/") {
			return true
		}
	}
	return false
}

// isKnownOAuthOrigin reports whether a URL points at an origin Fider actually serves.
// OAuth state is signed with a process-wide secret, so a valid signature on its own says
// nothing about whether the redirect it carries was ever legitimate — it has to be checked
// against the tenant records again.
func isKnownOAuthOrigin(c *web.Context, u *url.URL) bool {
	origin := u.Scheme + "://" + u.Host

	if origin == web.OAuthBaseURL(c) {
		return true
	}

	byDomain := &query.GetTenantByDomain{Domain: u.Hostname()}
	if err := bus.Dispatch(c, byDomain); err != nil || byDomain.Result == nil || byDomain.Result.IsDisabled() {
		return false
	}

	return slices.Contains(tenantOrigins(c, byDomain.Result), origin)
}

// newOAuthHandoff mints the token that carries an OAuth sign in from the callback address
// over to the tenant address, bound to that origin, provider and browser session.
func newOAuthHandoff(provider, origin, sessionID string) (string, error) {
	return jwt.Encode(jwt.OAuthHandoffClaims{
		Provider:      provider,
		Origin:        origin,
		SessionIDHash: hashSessionID(sessionID),
		Metadata: jwt.Metadata{
			ExpiresAt: jwt.Time(time.Now().Add(handoffLifetime)),
		},
	})
}

// isValidOAuthHandoff checks that a handoff token belongs to this exact request: same
// provider, same origin, same browser session, and not yet expired.
func isValidOAuthHandoff(c *web.Context, provider, handoff string) bool {
	if handoff == "" || c.SessionID() == "" {
		return false
	}

	claims, err := jwt.DecodeOAuthHandoffClaims(handoff)
	if err != nil {
		return false
	}

	if claims.Provider != provider || !isAllowedOAuthRedirect(oauthAllowedOrigins(c), claims.Origin) {
		return false
	}

	expected := hashSessionID(c.SessionID())
	return subtle.ConstantTimeCompare([]byte(claims.SessionIDHash), []byte(expected)) == 1
}

// hashSessionID hashes a session ID so the handoff token can be bound to a browser session
// without the session ID itself travelling through URLs, logs and referrers.
func hashSessionID(sessionID string) string {
	sum := sha256.Sum256([]byte(sessionID))
	return hex.EncodeToString(sum[:])
}

// sanitiseOAuthRedirect reduces a caller supplied redirect to a path on the current tenant.
// Anything pointing elsewhere falls back to the root, so this endpoint can never be used
// as an open redirect.
func sanitiseOAuthRedirect(c *web.Context, redirect string) string {
	u, err := url.Parse(redirect)
	if err != nil || redirect == "" {
		return "/"
	}

	// Protocol relative URLs ("//evil.com") parse with an empty scheme, so check the host
	// rather than relying on IsAbs.
	if u.Host != "" && !isAllowedOAuthRedirect(oauthAllowedOrigins(c), u.Scheme+"://"+u.Host) {
		return "/"
	}

	// Go parses "///evil.com" as a path with no host, but browsers skip any run of leading
	// slashes and treat what follows as the host, so a path must not start with "//".
	path := u.EscapedPath()
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return "/"
	}

	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	return path
}

// OAuthEcho exchanges OAuth Code for a user profile and return directly to the UI, without storing it
func OAuthEcho() web.HandlerFunc {
	return func(c *web.Context) error {
		provider := c.Param("provider")

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect("/")
		}

		if !isValidOAuthHandoff(c, provider, c.QueryParam("handoff")) {
			log.Warn(c, "OAuth handoff token is not valid for this request. Aborting sign in process.")
			return c.Redirect("/")
		}

		rawProfile := &query.GetOAuthRawProfile{Provider: provider, Code: code}
		err := bus.Dispatch(c, rawProfile)
		if err != nil {
			return c.Page(http.StatusOK, web.Props{
				Page:  "OAuthEcho/OAuthEcho.page",
				Title: "OAuth Test Page",
				Data: web.Map{
					"err": errors.Cause(err).Error(),
				},
			})
		}

		parseRawProfile := &cmd.ParseOAuthRawProfile{Provider: provider, Body: rawProfile.Result}
		_ = bus.Dispatch(c, parseRawProfile)

		// Fetch provider config to show configured allowedRoles on the test page.
		// Errors are intentionally ignored here — this is a non-critical diagnostic fetch.
		var configuredAllowedRoles, configuredRolesPath string
		if providerConfig, err := getCustomOAuthConfig(c, provider); err == nil && providerConfig != nil {
			configuredAllowedRoles = providerConfig.AllowedRoles
			configuredRolesPath = providerConfig.JSONUserRolesPath
		}

		return c.Page(http.StatusOK, web.Props{
			Page:  "OAuthEcho/OAuthEcho.page",
			Title: "OAuth Test Page",
			Data: web.Map{
				"body":                 rawProfile.Result,
				"profile":              parseRawProfile.Result,
				"configuredRolesPath":  configuredRolesPath,
				"configuredAllowedRoles": configuredAllowedRoles,
			},
		})
	}
}

// OAuthToken exchanges OAuth Code for a user profile
// The user profile is then used to either get an existing user on Fider or creating a new one
// Once Fider user is retrieved/created, an authentication cookie is store in user's browser
func OAuthToken() web.HandlerFunc {
	return func(c *web.Context) error {
		provider := c.Param("provider")
		redirect := sanitiseOAuthRedirect(c, c.QueryParam("redirect"))

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(redirect)
		}

		if !isValidOAuthHandoff(c, provider, c.QueryParam("handoff")) {
			log.Warn(c, "OAuth handoff token is not valid for this request. Aborting sign in process.")
			return c.Redirect(redirect)
		}

		oauthUser := &query.GetOAuthProfile{Provider: provider, Code: code}
		if err := bus.Dispatch(c, oauthUser); err != nil {
			return c.Failure(err)
		}

		// Fetch custom provider config once — used for role checking and trust check below.
		// Returns nil for built-in providers (Google, Facebook, …).
		// Returns an error if the DB is unavailable — we fail hard rather than silently
		// bypassing access controls.
		customConfig, err := getCustomOAuthConfig(c, provider)
		if err != nil {
			return c.Failure(err)
		}

		// Look up the existing Fider user first (by provider UID, then by email).
		// We need this before the role check so that administrators and collaborators
		// can always sign in regardless of OAuth role changes.
		var user *entity.User

		userByProvider := &query.GetUserByProvider{Provider: provider, UID: oauthUser.Result.ID}
		err = bus.Dispatch(c, userByProvider)
		user = userByProvider.Result

		if errors.Cause(err) == app.ErrNotFound && oauthUser.Result.Email != "" {
			userByEmail := &query.GetUserByEmail{Email: oauthUser.Result.Email}
			err = bus.Dispatch(c, userByEmail)
			user = userByEmail.Result
		}

		// Check if user has the required roles for this provider.
		// Both AllowedRoles and JSONUserRolesPath must be set on the provider for the check to run.
		// Administrators and collaborators already trusted in Fider are always allowed through,
		// regardless of their current OAuth roles.
		var providerRolesPath, providerAllowedRoles string
		if customConfig != nil {
			providerRolesPath = customConfig.JSONUserRolesPath
			providerAllowedRoles = customConfig.AllowedRoles
		}
		isFiderPrivileged := user != nil && (user.Role == enum.RoleAdministrator || user.Role == enum.RoleCollaborator)
		if !isFiderPrivileged && !hasAllowedRole(oauthUser.Result.Roles, providerRolesPath, providerAllowedRoles) {
			log.Warnf(c, "User @{UserID} attempted OAuth login but does not have required role. User roles: @{UserRoles}, Allowed roles: @{AllowedRoles}",
				dto.Props{
					"UserID":       oauthUser.Result.ID,
					"UserRoles":    oauthUser.Result.Roles,
					"AllowedRoles": providerAllowedRoles,
				})
			return c.Redirect("/access-denied")
		}
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				isTrusted := customConfig != nil && customConfig.IsTrusted
				if c.Tenant().IsPrivate && !isTrusted {
					return c.Redirect("/not-invited")
				}

				user = &entity.User{
					Name:   oauthUser.Result.Name,
					Tenant: c.Tenant(),
					Email:  oauthUser.Result.Email,
					Role:   enum.RoleVisitor,
					Providers: []*entity.UserProvider{
						{
							UID:  oauthUser.Result.ID,
							Name: provider,
						},
					},
				}

				if err = bus.Dispatch(c, &cmd.RegisterUser{User: user}); err != nil {
					return c.Failure(err)
				}
			} else {
				return c.Failure(err)
			}
		} else if !user.HasProvider(provider) {
			if err = bus.Dispatch(c, &cmd.RegisterUserProvider{
				UserID:       user.ID,
				ProviderName: provider,
				ProviderUID:  oauthUser.Result.ID,
			}); err != nil {
				return c.Failure(err)
			}
		}

		webutil.AddAuthUserCookie(c, user)

		return c.Redirect(redirect)
	}
}

// getCustomOAuthConfig fetches the custom OAuth provider config for the given provider.
// Built-in providers (Google, Facebook, GitHub, …) are identified by the absence of a
// leading "_" and never have a custom config row, so the bus dispatch is skipped for them.
// Returns (nil, nil) for built-in providers.
// Returns (nil, err) if the DB lookup fails — callers must treat this as a hard error so
// that a transient DB outage cannot silently bypass access controls.
func getCustomOAuthConfig(ctx context.Context, provider string) (*entity.OAuthConfig, error) {
	if len(provider) == 0 || provider[0] != '_' {
		return nil, nil
	}
	q := &query.GetCustomOAuthConfigByProvider{Provider: provider}
	if err := bus.Dispatch(ctx, q); err != nil {
		return nil, err
	}
	return q.Result, nil
}

// OAuthCallback handles the redirect back from the OAuth provider
// This callback can run on either Tenant or Login address
// If the request is for a sign in, we redirect the user to the tenant address
// If the request is for a sign up, we exchange the OAuth code and get the user profile
func OAuthCallback() web.HandlerFunc {
	return func(c *web.Context) error {
		c.Response.Header().Add("X-Robots-Tag", "noindex")

		provider := c.Param("provider")
		state := c.QueryParam("state")
		claims, err := jwt.DecodeOAuthStateClaims(state)
		if err != nil {
			return c.Forbidden()
		}

		if claims.Redirect == "" {
			log.Warnf(c, "Missing redirect URL in OAuth callback state for provider @{Provider}.", dto.Props{"Provider": provider})
			return c.NotFound()
		}

		redirectURL, err := url.ParseRequestURI(claims.Redirect)
		if err != nil {
			return c.Failure(err)
		}

		if !isKnownOAuthOrigin(c, redirectURL) {
			log.Warnf(c, "OAuth callback state for provider @{Provider} points at an unknown origin. Aborting sign in process.", dto.Props{"Provider": provider})
			return c.Forbidden()
		}

		origin := redirectURL.Scheme + "://" + redirectURL.Host

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(redirectURL.String())
		}

		//Test OAuth
		if redirectURL.Path == fmt.Sprintf("/oauth/%s/echo", provider) {
			handoff, err := newOAuthHandoff(provider, origin, claims.Identifier)
			if err != nil {
				return c.Failure(err)
			}

			var query = redirectURL.Query()
			query.Set("code", code)
			query.Set("handoff", handoff)
			redirectURL.RawQuery = query.Encode()
			return c.Redirect(redirectURL.String())
		}

		//Sign up process
		if redirectURL.Path == "/signup" {
			oauthUser := &query.GetOAuthProfile{Provider: provider, Code: code}
			if err := bus.Dispatch(c, oauthUser); err != nil {
				return c.Failure(err)
			}

			claims := jwt.OAuthClaims{
				OAuthID:       oauthUser.Result.ID,
				OAuthProvider: provider,
				OAuthName:     oauthUser.Result.Name,
				OAuthEmail:    oauthUser.Result.Email,
				Metadata: jwt.Metadata{
					ExpiresAt: jwt.Time(time.Now().Add(10 * time.Minute)),
				},
			}

			token, err := jwt.Encode(claims)
			if err != nil {
				return c.Failure(err)
			}

			var query = redirectURL.Query()
			query.Set("token", token)
			redirectURL.RawQuery = query.Encode()
			return c.Redirect(redirectURL.String())
		}

		//Sign in process
		handoff, err := newOAuthHandoff(provider, origin, claims.Identifier)
		if err != nil {
			return c.Failure(err)
		}

		var query = redirectURL.Query()
		query.Set("code", code)
		query.Set("redirect", redirectURL.RequestURI())
		query.Set("handoff", handoff)
		redirectURL.RawQuery = query.Encode()
		redirectURL.Path = fmt.Sprintf("/oauth/%s/token", provider)
		return c.Redirect(redirectURL.String())
	}
}

// SignInByOAuth is responsible for redirecting the user to the OAuth authorization URL for given provider
// A cookie is stored in user's browser with a random identifier that is later used to verify the authenticity of the request
func SignInByOAuth() web.HandlerFunc {
	return func(c *web.Context) error {
		c.Response.Header().Add("X-Robots-Tag", "noindex")

		provider := c.Param("provider")

		// The origin the user is sent back to has to come from the resolved tenant, not from
		// the request host: Host and X-Forwarded-Host are caller controlled, so deriving it
		// from them lets anyone have the provider's authorization code delivered to an origin
		// they own and then replay it against a tenant of their choosing.
		origins := oauthAllowedOrigins(c)

		redirect := c.QueryParam("redirect")
		if redirect == "" {
			redirect = origins[0]
		} else if !isAllowedOAuthRedirect(origins, redirect) {
			return c.Forbidden()
		}

		redirectURL, err := url.ParseRequestURI(redirect)
		if err != nil {
			return c.Forbidden()
		}

		// Without a tenant the only flow that can complete is sign up, which runs on the
		// host-wide OAuth address. Everything else needs a tenant to sign in to.
		if c.Tenant() == nil && redirectURL.Path != "/signup" {
			return c.Forbidden()
		}

		if c.IsAuthenticated() && redirectURL.Path != fmt.Sprintf("/oauth/%s/echo", provider) {
			return c.Redirect(redirect)
		}

		authURL := &query.GetOAuthAuthorizationURL{
			Provider:   provider,
			Redirect:   redirect,
			Identifier: c.SessionID(),
		}
		if err := bus.Dispatch(c, authURL); err != nil {
			return c.Failure(err)
		}
		return c.Redirect(authURL.Result)
	}
}

// hasAllowedRole checks if the user has any of the allowed roles configured on the provider.
// If allowedRoles is empty, all users are allowed (returns true).
// If jsonUserRolesPath is empty, the role check is skipped (returns true) — this ensures
// providers without a roles path are never accidentally blocked.
func hasAllowedRole(userRoles []string, jsonUserRolesPath string, allowedRoles string) bool {
	allowedRolesConfig := strings.TrimSpace(allowedRoles)

	// If no roles restriction is configured on this provider, allow all users
	if allowedRolesConfig == "" {
		return true
	}

	// If the provider has no roles path configured, skip the role check for this provider
	if strings.TrimSpace(jsonUserRolesPath) == "" {
		return true
	}

	// Parse allowed roles from config (comma-separated)
	allowedRolesList := strings.Split(allowedRolesConfig, ",")
	allowedRolesMap := make(map[string]bool)
	for _, role := range allowedRolesList {
		role = strings.TrimSpace(role)
		if role != "" {
			allowedRolesMap[role] = true
		}
	}

	// If no valid roles in config, allow all
	if len(allowedRolesMap) == 0 {
		return true
	}

	// Check if user has any of the allowed roles
	for _, userRole := range userRoles {
		userRole = strings.TrimSpace(userRole)
		if allowedRolesMap[userRole] {
			return true
		}
	}

	// User doesn't have any of the required roles
	return false
}

