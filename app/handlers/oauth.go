package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	webutil "github.com/getfider/fider/app/pkg/web/util"
)

// OAuthEcho exchanges OAuth Code for a user profile and return directly to the UI, without storing it
func OAuthEcho() web.HandlerFunc {
	return func(c *web.Context) error {
		provider := c.Param("provider")

		code := c.QueryParam("code")
		if code == "" {
			return c.RedirectTo("/")
		}

		identifier := c.QueryParam("identifier")
		if identifier == "" || identifier != c.SessionID() {
			log.Warn(c, "OAuth identifier doesn't match with user session ID. Aborting sign in process.")
			return c.RedirectTo("/")
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
		redirectURL, _ := url.ParseRequestURI(c.QueryParam("redirect"))
		redirectURL.ResolveReference(c.Request.URL)

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(redirectURL.String())
		}

		identifier := c.QueryParam("identifier")
		if identifier == "" || identifier != c.SessionID() {
			log.Warn(c, "OAuth identifier doesn't match with user session ID. Aborting sign in process.")
			return c.Redirect(redirectURL.String())
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
			return c.RedirectTo("/access-denied")
		}
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				isTrusted := customConfig != nil && customConfig.IsTrusted
				if c.Tenant().IsPrivate && !isTrusted {
					return c.RedirectTo("/not-invited")
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

		return c.Redirect(redirectURL.String())
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

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(redirectURL.String())
		}

		//Test OAuth
		if redirectURL.Path == fmt.Sprintf("/oauth/%s/echo", provider) {
			var query = redirectURL.Query()
			query.Set("code", code)
			query.Set("identifier", claims.Identifier)
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
		var query = redirectURL.Query()
		query.Set("code", code)
		query.Set("redirect", redirectURL.RequestURI())
		query.Set("identifier", claims.Identifier)
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
		redirect := c.QueryParam("redirect")

		if redirect == "" {
			redirect = c.BaseURL()
		} else if redirect != c.BaseURL() && !strings.HasPrefix(redirect, c.BaseURL()+"/") {
			return c.Forbidden()
		}

		redirectURL, _ := url.ParseRequestURI(redirect)
		redirectURL.ResolveReference(c.Request.URL)

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

