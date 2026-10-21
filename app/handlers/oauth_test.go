package handlers_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"

	"github.com/getfider/fider/app/services/oauth"

	"github.com/getfider/fider/app"

	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func initOAuthWithMocks() {
	bus.Init(&oauth.Service{})
	// Mock tenant provider status to return enabled by default
	bus.AddHandler(func(ctx context.Context, q *query.GetTenantProviderStatus) error {
		q.Result = &entity.TenantProvider{
			Provider:  q.Provider,
			IsEnabled: true,
		}
		return nil
	})
	mockTenantByDomain()
}

// mockTenantByDomain resolves the hostnames the mocked tenants are served on, so that the
// OAuth handlers can check a redirect origin against a real tenant record.
func mockTenantByDomain() {
	bus.AddHandler(func(ctx context.Context, q *query.GetTenantByDomain) error {
		switch q.Domain {
		case "demo.test.fider.io":
			q.Result = mock.DemoTenant
			return nil
		case "avengers.test.fider.io", "feedback.theavengers.com":
			q.Result = mock.AvengersTenant
			return nil
		}
		return app.ErrNotFound
	})
}

// newHandoff mints the token OAuthCallback would hand over to the tenant address.
func newHandoff(provider, origin, sessionID string) string {
	sum := sha256.Sum256([]byte(sessionID))
	token, err := jwt.Encode(jwt.OAuthHandoffClaims{
		Provider:      provider,
		Origin:        origin,
		SessionIDHash: hex.EncodeToString(sum[:]),
		Metadata: jwt.Metadata{
			ExpiresAt: jwt.Time(time.Now().Add(2 * time.Minute)),
		},
	})
	Expect(err).IsNil()
	return token
}

func TestSignOutHandler(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/signout").
		AddCookie(web.CookieAuthName, "some-value").
		Execute(handlers.SignOut())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/")
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring(web.CookieAuthName + "=; Path=/; Expires=")
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring("Max-Age=0; HttpOnly")
}

func TestSignInByOAuthHandler_RootRedirect(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, _ := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://avengers.test.fider.io/oauth/facebook?redirect=http://avengers.test.fider.io").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusTemporaryRedirect)
}

func TestSignInByOAuthHandler_PathRedirect(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, _ := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://avengers.test.fider.io/oauth/facebook?redirect=http://avengers.test.fider.io/something").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusTemporaryRedirect)
}

// A tenant with a custom domain is still served on its subdomain, so both are valid targets.
func TestSignInByOAuthHandler_CustomDomainRedirect(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, _ := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://feedback.theavengers.com/oauth/facebook?redirect=http://feedback.theavengers.com/something").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusTemporaryRedirect)
}

func TestSignInByOAuthHandler_EvilRedirect(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, _ := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://avengers.test.fider.io/oauth/facebook?redirect=http://evil.com").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusForbidden)
}

func TestSignInByOAuthHandler_EvilRedirect2(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, _ := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://avengers.test.fider.io/oauth/facebook?redirect=http://avengers.test.fider.io.evil.com").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusForbidden)
}

// The redirect has to be built from the tenant record, never from the request host, which
// any caller can set with Host or X-Forwarded-Host. A poisoned host falls back to the
// tenant's canonical origin rather than signing the attacker's origin into the state.
func TestSignInByOAuthHandler_PoisonedHost(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, response := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://evil.com/oauth/facebook").
		AddHeader("X-Forwarded-Host", "evil.com").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	ExpectOAuthState(response, "http://feedback.theavengers.com", "MY_SESSION_ID")
}

// The poisoned host can't be laundered through the redirect parameter either.
func TestSignInByOAuthHandler_PoisonedHost_AsRedirect(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, _ := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://evil.com/oauth/facebook?redirect=http://evil.com").
		AddHeader("X-Forwarded-Host", "evil.com").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusForbidden)
}

// Same attack, but on a host that resolves to no tenant at all — which is reachable
// because this route sits in front of the RequireTenant middleware.
func TestSignInByOAuthHandler_PoisonedHost_WithoutTenant(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, _ := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://evil.com/oauth/facebook").
		AddHeader("X-Forwarded-Host", "evil.com").
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusForbidden)
}

// Without a tenant, sign up on the host-wide OAuth address is the only flow that can complete.
func TestSignInByOAuthHandler_WithoutTenant_SignUp(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, response := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://login.test.fider.io/oauth/facebook?redirect=http://login.test.fider.io/signup").
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	ExpectOAuthState(response, "http://login.test.fider.io/signup", "MY_SESSION_ID")
}

func TestSignInByOAuthHandler_WithoutTenant_NotSignUp(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, _ := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://login.test.fider.io/oauth/facebook?redirect=http://login.test.fider.io/something").
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusForbidden)
}

func TestSignInByOAuthHandler_InvalidURL(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, response := server.
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://avengers.test.fider.io/oauth/facebook").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	// No redirect given, so it defaults to the tenant's canonical origin — its custom domain.
	ExpectOAuthState(response, "http://feedback.theavengers.com", "MY_SESSION_ID")
}

func TestSignInByOAuthHandler_AuthenticatedUser(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, response := server.
		AsUser(mock.JonSnow).
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://avengers.test.fider.io/oauth/facebook?redirect=http://avengers.test.fider.io").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io")
}

func TestSignInByOAuthHandler_AuthenticatedUser_UsingEcho(t *testing.T) {
	RegisterT(t)
	initOAuthWithMocks()

	server := mock.NewServer()
	code, response := server.
		AsUser(mock.JonSnow).
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		WithURL("http://avengers.test.fider.io/oauth/facebook?redirect=http://avengers.test.fider.io/oauth/facebook/echo").
		OnTenant(mock.AvengersTenant).
		Use(middlewares.Session()).
		Execute(handlers.SignInByOAuth())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	ExpectOAuthState(response, "http://avengers.test.fider.io/oauth/facebook/echo", "MY_SESSION_ID")
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterT(t)
	mockTenantByDomain()

	server := mock.NewServer()
	code, _ := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=abc").
		AddParam("provider", app.FacebookProvider).
		Execute(handlers.OAuthCallback())

	Expect(code).Equals(http.StatusForbidden)
}

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterT(t)
	mockTenantByDomain()

	server := mock.NewServer()
	state, _ := jwt.Encode(jwt.OAuthStateClaims{
		Redirect:   "http://avengers.test.fider.io",
		Identifier: "",
	})

	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state="+state).
		AddParam("provider", app.FacebookProvider).
		Execute(handlers.OAuthCallback())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io")
}

// State is signed with a process-wide secret, so a valid signature alone doesn't make the
// origin it carries legitimate — it has to resolve to a tenant Fider actually serves.
func TestCallbackHandler_UnknownOrigin(t *testing.T) {
	RegisterT(t)
	mockTenantByDomain()

	state, _ := jwt.Encode(jwt.OAuthStateClaims{
		Redirect:   "http://evil.com",
		Identifier: "888",
	})

	server := mock.NewServer()
	code, _ := server.
		WithURL("http://login.test.fider.io/oauth/callback?state="+state+"&code=123").
		AddParam("provider", app.FacebookProvider).
		Execute(handlers.OAuthCallback())

	Expect(code).Equals(http.StatusForbidden)
}

func TestCallbackHandler_SignIn(t *testing.T) {
	RegisterT(t)
	mockTenantByDomain()

	state, _ := jwt.Encode(jwt.OAuthStateClaims{
		Redirect:   "http://avengers.test.fider.io",
		Identifier: "888",
	})

	server := mock.NewServer()
	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state="+state+"&code=123").
		AddParam("provider", app.FacebookProvider).
		Execute(handlers.OAuthCallback())

	Expect(code).Equals(http.StatusTemporaryRedirect)

	location, _ := url.Parse(response.Header().Get("Location"))
	Expect(location.Host).Equals("avengers.test.fider.io")
	Expect(location.Path).Equals("/oauth/facebook/token")
	Expect(location.Query().Get("code")).Equals("123")
	Expect(location.Query().Get("redirect")).Equals("/")
	// The session ID itself must not travel through the URL any more.
	Expect(location.Query().Get("identifier")).IsEmpty()
	ExpectOAuthHandoff(location.Query().Get("handoff"), "facebook", "http://avengers.test.fider.io", "888")
}

func TestCallbackHandler_SignIn_WithPath(t *testing.T) {
	RegisterT(t)
	mockTenantByDomain()
	server := mock.NewServer()

	state, _ := jwt.Encode(jwt.OAuthStateClaims{
		Redirect:   "http://avengers.test.fider.io/some-page",
		Identifier: "888",
	})

	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state="+state+"&code=123").
		AddParam("provider", app.FacebookProvider).
		Execute(handlers.OAuthCallback())

	Expect(code).Equals(http.StatusTemporaryRedirect)

	location, _ := url.Parse(response.Header().Get("Location"))
	Expect(location.Host).Equals("avengers.test.fider.io")
	Expect(location.Path).Equals("/oauth/facebook/token")
	Expect(location.Query().Get("code")).Equals("123")
	Expect(location.Query().Get("redirect")).Equals("/some-page")
	ExpectOAuthHandoff(location.Query().Get("handoff"), "facebook", "http://avengers.test.fider.io", "888")
}

func TestCallbackHandler_SignUp(t *testing.T) {
	RegisterT(t)
	mockTenantByDomain()

	oauthUser := &dto.OAuthUserProfile{
		ID:    "FB123",
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetOAuthProfile) error {
		if q.Provider == app.FacebookProvider && q.Code == "123" {
			q.Result = oauthUser
			return nil
		}
		return app.ErrNotFound
	})

	state, _ := jwt.Encode(jwt.OAuthStateClaims{
		Redirect:   "http://demo.test.fider.io/signup",
		Identifier: "",
	})

	server := mock.NewServer()
	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state="+state+"&code=123").
		AddParam("provider", app.FacebookProvider).
		Execute(handlers.OAuthCallback())
	Expect(code).Equals(http.StatusTemporaryRedirect)

	location, _ := url.Parse(response.Header().Get("Location"))
	Expect(location.Host).Equals("demo.test.fider.io")
	Expect(location.Scheme).Equals("http")
	Expect(location.Path).Equals("/signup")
	ExpectOAuthToken(location.Query().Get("token"), &jwt.OAuthClaims{
		OAuthProvider: "facebook",
		OAuthID:       oauthUser.ID,
		OAuthName:     oauthUser.Name,
		OAuthEmail:    oauthUser.Email,
	})
}

func TestOAuthTokenHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterT(t)

	oauthUser := &dto.OAuthUserProfile{
		ID:    "FB123",
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetOAuthProfile) error {
		if q.Provider == app.FacebookProvider && q.Code == "123" {
			q.Result = oauthUser
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		if q.Provider == app.FacebookProvider && q.UID == oauthUser.ID {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=123&redirect=/hello&handoff="+newHandoff(app.FacebookProvider, "http://demo.test.fider.io", "MY_SESSION_ID")).
		OnTenant(mock.DemoTenant).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		AddParam("provider", app.FacebookProvider).
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/hello")
	ExpectFiderAuthCookie(response, mock.JonSnow)
}

func TestOAuthTokenHandler_NewUser(t *testing.T) {
	RegisterT(t)

	var registeredUser *entity.User
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		registeredUser = c.User
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		Expect(q.Provider).Equals(app.FacebookProvider)
		Expect(q.UID).Equals("FB456")
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		Expect(q.Email).Equals("some.guy@facebook.com")
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		Expect(q.Provider).Equals(app.FacebookProvider)
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetOAuthProfile) error {
		Expect(q.Provider).Equals(app.FacebookProvider)
		Expect(q.Code).Equals("456")

		q.Result = &dto.OAuthUserProfile{
			ID:    "FB456",
			Name:  "Some Facebook Guy",
			Email: "some.guy@facebook.com",
		}
		return nil
	})

	server := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=456&redirect=/hello&handoff="+newHandoff(app.FacebookProvider, "http://demo.test.fider.io", "MY_SESSION_ID")).
		OnTenant(mock.DemoTenant).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		AddParam("provider", app.FacebookProvider).
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/hello")

	Expect(registeredUser.Name).Equals("Some Facebook Guy")

	ExpectFiderAuthCookie(response, registeredUser)
}

func TestOAuthTokenHandler_NewUserWithoutEmail(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	var newUser *entity.User
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		c.User.ID = 1
		newUser = c.User
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetOAuthProfile) error {
		if q.Provider == app.FacebookProvider && q.Code == "798" {
			q.Result = &dto.OAuthUserProfile{
				ID:    "FB798",
				Name:  "Mark",
				Email: "",
			}
			return nil
		}
		return app.ErrNotFound
	})

	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=798&redirect=/&handoff="+newHandoff(app.FacebookProvider, "http://demo.test.fider.io", "MY_SESSION_ID")).
		OnTenant(mock.DemoTenant).
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())

	Expect(newUser.ID).Equals(1)
	Expect(newUser.Name).Equals("Mark")
	Expect(newUser.Providers).HasLen(1)

	Expect(code).Equals(http.StatusTemporaryRedirect)

	Expect(response.Header().Get("Location")).Equals("/")
	ExpectFiderAuthCookie(response, &entity.User{
		ID:   1,
		Name: "Mark",
	})
}

func TestOAuthTokenHandler_ExistingUser_WithoutEmail(t *testing.T) {
	RegisterT(t)

	user := &entity.User{
		ID:     3,
		Name:   "Some Facebook Guy",
		Email:  "",
		Tenant: mock.DemoTenant,
		Providers: []*entity.UserProvider{
			{UID: "FB456", Name: app.FacebookProvider},
		},
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetOAuthProfile) error {
		if q.Provider == app.FacebookProvider && q.Code == "456" {
			q.Result = &dto.OAuthUserProfile{
				ID:    "FB456",
				Name:  "Some Facebook Guy",
				Email: "some.guy@facebook.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		if q.Provider == "facebook" && q.UID == "FB456" {
			q.Result = user
			return nil
		}
		return app.ErrNotFound
	})

	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=456&redirect=/&handoff="+newHandoff(app.FacebookProvider, "http://demo.test.fider.io", "MY_SESSION_ID")).
		OnTenant(mock.DemoTenant).
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())

	Expect(code).Equals(http.StatusTemporaryRedirect)

	Expect(response.Header().Get("Location")).Equals("/")
	ExpectFiderAuthCookie(response, &entity.User{
		ID:    3,
		Name:  "Some Facebook Guy",
		Email: "",
	})
}

func TestOAuthTokenHandler_ExistingUser_NewProvider(t *testing.T) {
	RegisterT(t)

	var newProvider *entity.UserProvider
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUserProvider) error {
		newProvider = &entity.UserProvider{
			Name: c.ProviderName,
			UID:  c.ProviderUID,
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetOAuthProfile) error {
		if q.Provider == app.GoogleProvider && q.Code == "123" {
			q.Result = &dto.OAuthUserProfile{
				ID:    "GO123",
				Name:  "Jon Snow",
				Email: "jon.snow@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		if q.Provider == app.GoogleProvider && q.UID == "GO123" {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/google/token?code=123&redirect=/&handoff="+newHandoff(app.GoogleProvider, "http://demo.test.fider.io", "MY_SESSION_ID")).
		OnTenant(mock.DemoTenant).
		AddParam("provider", app.GoogleProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())

	Expect(code).Equals(http.StatusTemporaryRedirect)

	Expect(newProvider.Name).Equals("google")
	Expect(newProvider.UID).Equals("GO123")

	Expect(response.Header().Get("Location")).Equals("/")
	ExpectFiderAuthCookie(response, mock.JonSnow)
}

func TestOAuthTokenHandler_NewUser_PrivateSite(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.AvengersTenant.IsPrivate = true

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		Expect(q.Provider).Equals(app.FacebookProvider)
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetOAuthProfile) error {
		Expect(q.Provider).Equals(app.FacebookProvider)
		Expect(q.Code).Equals("456")
		q.Result = &dto.OAuthUserProfile{
			ID:    "FB456",
			Name:  "Some Facebook Guy",
			Email: "some.guy@facebook.com",
		}
		return nil
	})

	code, response := server.
		WithURL("http://feedback.theavengers.com/oauth/facebook/token?code=456&redirect=/&handoff="+newHandoff(app.FacebookProvider, "http://feedback.theavengers.com", "MY_SESSION_ID")).
		OnTenant(mock.AvengersTenant).
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/not-invited")
	ExpectFiderAuthCookie(response, nil)
}

func TestOAuthTokenHandler_NewUser_SignUpDisabled(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.AvengersTenant.IsSignUpDisabled = true

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		Expect(q.Provider).Equals(app.FacebookProvider)
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetOAuthProfile) error {
		Expect(q.Provider).Equals(app.FacebookProvider)
		Expect(q.Code).Equals("456")
		q.Result = &dto.OAuthUserProfile{
			ID:    "FB456",
			Name:  "Some Facebook Guy",
			Email: "some.guy@facebook.com",
		}
		return nil
	})

	code, response := server.
		WithURL("http://feedback.theavengers.com/oauth/facebook/token?code=456&redirect=/&handoff="+newHandoff(app.FacebookProvider, "http://feedback.theavengers.com", "MY_SESSION_ID")).
		OnTenant(mock.AvengersTenant).
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/not-invited")
	ExpectFiderAuthCookie(response, nil)
}

func TestOAuthTokenHandler_NewUser_PrivateSite_UsingTrustedProvider(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.AvengersTenant.IsPrivate = true

	providerCode := "_jd72hfjv"

	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		Expect(c.User.Name).Equals("Mark Doe")
		Expect(c.User.Email).Equals("mark.doe@microsoft.com")
		Expect(c.User.Providers).HasLen(1)
		Expect(c.User.Providers[0].UID).Equals("1234-5678")
		Expect(c.User.Providers[0].Name).Equals(providerCode)
		Expect(c.User.Role).Equals(enum.RoleVisitor)

		c.User.ID = 999
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		Expect(q.Provider).Equals(providerCode)
		q.Result = &entity.OAuthConfig{
			Provider:    providerCode,
			DisplayName: "Microsoft AD",
			IsTrusted:   true,
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetOAuthProfile) error {
		Expect(q.Provider).Equals(providerCode)
		Expect(q.Code).Equals("000111")
		q.Result = &dto.OAuthUserProfile{
			ID:    "1234-5678",
			Name:  "Mark Doe",
			Email: "mark.doe@microsoft.com",
		}
		return nil
	})

	code, response := server.
		WithURL("http://feedback.theavengers.com/oauth/"+providerCode+"/token?code=000111&redirect=/&handoff="+newHandoff(providerCode, "http://feedback.theavengers.com", "MY_SESSION_ID")).
		OnTenant(mock.AvengersTenant).
		AddParam("provider", providerCode).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())

	Expect(code).Equals(http.StatusTemporaryRedirect)

	Expect(response.Header().Get("Location")).Equals("/")
	ExpectFiderAuthCookie(response, &entity.User{
		ID:    999,
		Name:  "Mark Doe",
		Email: "mark.doe@microsoft.com",
	})
}

// executeOAuthTokenWithHandoff runs OAuthToken on the Avengers tenant with the given
// handoff token and returns the response, so the rejection cases can share a setup.
func executeOAuthTokenWithHandoff(handoff string) (int, *httptest.ResponseRecorder) {
	return mock.NewServer().
		WithURL("http://feedback.theavengers.com/oauth/facebook/token?code=456&redirect=/&handoff="+handoff).
		OnTenant(mock.AvengersTenant).
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())
}

func TestOAuthTokenHandler_MissingHandoff(t *testing.T) {
	RegisterT(t)

	code, response := executeOAuthTokenWithHandoff("")

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/")
	ExpectFiderAuthCookie(response, nil)
}

// The browser that started the flow is the only one allowed to finish it. A handoff must
// not be usable by a caller who simply sets their own user_session_id cookie.
func TestOAuthTokenHandler_HandoffForAnotherSession(t *testing.T) {
	RegisterT(t)

	handoff := newHandoff(app.FacebookProvider, "http://feedback.theavengers.com", "SOME_OTHER_ID")
	code, response := executeOAuthTokenWithHandoff(handoff)

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/")
	ExpectFiderAuthCookie(response, nil)
}

// A handoff issued for one tenant must not be redeemable on another.
func TestOAuthTokenHandler_HandoffForAnotherTenant(t *testing.T) {
	RegisterT(t)

	handoff := newHandoff(app.FacebookProvider, "http://demo.test.fider.io", "MY_SESSION_ID")
	code, response := executeOAuthTokenWithHandoff(handoff)

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/")
	ExpectFiderAuthCookie(response, nil)
}

func TestOAuthTokenHandler_HandoffForAnotherProvider(t *testing.T) {
	RegisterT(t)

	handoff := newHandoff(app.GoogleProvider, "http://feedback.theavengers.com", "MY_SESSION_ID")
	code, response := executeOAuthTokenWithHandoff(handoff)

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/")
	ExpectFiderAuthCookie(response, nil)
}

func TestOAuthTokenHandler_ExpiredHandoff(t *testing.T) {
	RegisterT(t)

	sum := sha256.Sum256([]byte("MY_SESSION_ID"))
	handoff, _ := jwt.Encode(jwt.OAuthHandoffClaims{
		Provider:      app.FacebookProvider,
		Origin:        "http://feedback.theavengers.com",
		SessionIDHash: hex.EncodeToString(sum[:]),
		Metadata: jwt.Metadata{
			ExpiresAt: jwt.Time(time.Now().Add(-1 * time.Minute)),
		},
	})

	code, response := executeOAuthTokenWithHandoff(handoff)

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/")
	ExpectFiderAuthCookie(response, nil)
}

// This endpoint must never send the browser to another origin, even before sign in.
func TestOAuthTokenHandler_OpenRedirect(t *testing.T) {
	RegisterT(t)

	// Browsers read any run of leading slashes, or a backslash, as the start of a host.
	for redirect, location := range map[string]string{
		"http://evil.com/x": "/",
		"//evil.com/x":      "/",
		"///evil.com/x":     "/",
		"////evil.com/x":    "/",
		"https://feedback.theavengers.com.evil.com/x": "/",
		"/\\evil.com/x":   "/%5Cevil.com/x",
		"/%2F/evil.com/x": "/%2F/evil.com/x",
	} {
		code, response := mock.NewServer().
			WithURL("http://feedback.theavengers.com/oauth/facebook/token?redirect="+url.QueryEscape(redirect)).
			OnTenant(mock.AvengersTenant).
			AddParam("provider", app.FacebookProvider).
			AddCookie(web.CookieSessionName, "MY_SESSION_ID").
			Use(middlewares.Session()).
			Execute(handlers.OAuthToken())

		Expect(code).Equals(http.StatusTemporaryRedirect)
		Expect(response.Header().Get("Location")).Equals(location)
	}
}

// A missing redirect used to be parsed into a nil URL and then dereferenced.
func TestOAuthTokenHandler_MissingRedirect(t *testing.T) {
	RegisterT(t)

	code, response := mock.NewServer().
		WithURL("http://feedback.theavengers.com/oauth/facebook/token").
		OnTenant(mock.AvengersTenant).
		AddParam("provider", app.FacebookProvider).
		AddCookie(web.CookieSessionName, "MY_SESSION_ID").
		Use(middlewares.Session()).
		Execute(handlers.OAuthToken())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/")
}

func ExpectOAuthToken(token string, expected *jwt.OAuthClaims) {
	user, err := jwt.DecodeOAuthClaims(token)
	Expect(err).IsNil()
	Expect(user.OAuthID).Equals(expected.OAuthID)
	Expect(user.OAuthName).Equals(expected.OAuthName)
	Expect(user.OAuthEmail).Equals(expected.OAuthEmail)
	Expect(user.OAuthProvider).Equals(expected.OAuthProvider)
}

// ExpectOAuthState checks the provider authorization URL and the claims signed into its state.
func ExpectOAuthState(response *httptest.ResponseRecorder, expectedRedirect, expectedIdentifier string) {
	location, err := url.Parse(response.Header().Get("Location"))
	Expect(err).IsNil()
	Expect(location.Host).Equals("www.facebook.com")
	Expect(location.Query().Get("redirect_uri")).Equals("http://login.test.fider.io/oauth/facebook/callback")

	claims, err := jwt.DecodeOAuthStateClaims(location.Query().Get("state"))
	Expect(err).IsNil()
	Expect(claims.Redirect).Equals(expectedRedirect)
	Expect(claims.Identifier).Equals(expectedIdentifier)
	Expect(claims.ExpiresAt).IsNotNil()
}

// ExpectOAuthHandoff checks a handoff token is bound to the expected provider, origin and session.
func ExpectOAuthHandoff(handoff, expectedProvider, expectedOrigin, expectedSessionID string) {
	claims, err := jwt.DecodeOAuthHandoffClaims(handoff)
	Expect(err).IsNil()
	Expect(claims.Provider).Equals(expectedProvider)
	Expect(claims.Origin).Equals(expectedOrigin)
	Expect(claims.ExpiresAt).IsNotNil()

	sum := sha256.Sum256([]byte(expectedSessionID))
	Expect(claims.SessionIDHash).Equals(hex.EncodeToString(sum[:]))
}
