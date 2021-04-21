package oauth_test

import (
	"context"
	"crypto/tls"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/services/oauth"
)

func newGetContext(rawurl string) *web.Context {
	u, _ := url.Parse(rawurl)
	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", u.RequestURI(), nil)
	req.Host = u.Host

	if u.Scheme == "https" {
		req.TLS = &tls.ConnectionState{}
	}

	return web.NewContext(e, req, res, nil)
}
func TestGetAuthURL_Facebook(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	ctx := newGetContext("http://login.test.fider.io:3000")
	authURL := &query.GetOAuthAuthorizationURL{
		Provider:   app.FacebookProvider,
		Redirect:   "http://example.org",
		Identifier: "456",
	}

	err := bus.Dispatch(ctx, authURL)
	Expect(err).IsNil()
	Expect(authURL.Result).Equals("https://www.facebook.com/v3.2/dialog/oauth?client_id=FB_CL_ID&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Ffacebook%2Fcallback&response_type=code&scope=public_profile+email&state=http%3A%2F%2Fexample.org%7C456")
}

func TestGetAuthURL_Google(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	ctx := newGetContext("http://login.test.fider.io:3000")

	authURL := &query.GetOAuthAuthorizationURL{
		Provider:   app.GoogleProvider,
		Redirect:   "http://example.org",
		Identifier: "123",
	}

	err := bus.Dispatch(ctx, authURL)
	Expect(err).IsNil()
	Expect(authURL.Result).Equals("https://accounts.google.com/o/oauth2/v2/auth?client_id=GO_CL_ID&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Fgoogle%2Fcallback&response_type=code&scope=profile+email&state=http%3A%2F%2Fexample.org%7C123")
}

func TestGetAuthURL_GitHub(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	ctx := newGetContext("http://login.test.fider.io:3000")

	authURL := &query.GetOAuthAuthorizationURL{
		Provider:   app.GitHubProvider,
		Redirect:   "http://example.org",
		Identifier: "456",
	}

	err := bus.Dispatch(ctx, authURL)
	Expect(err).IsNil()
	Expect(authURL.Result).Equals("https://github.com/login/oauth/authorize?client_id=GH_CL_ID&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Fgithub%2Fcallback&response_type=code&scope=user%3Aemail&state=http%3A%2F%2Fexample.org%7C456")
}

func TestGetAuthURL_Custom(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_custom" {
			q.Result = &models.OAuthConfig{
				Provider:     q.Provider,
				ClientID:     "CU_CL_ID",
				Scope:        "profile email",
				AuthorizeURL: "https://example.org/oauth/authorize",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	authURL := &query.GetOAuthAuthorizationURL{
		Provider:   "_custom",
		Redirect:   "http://example.org",
		Identifier: "456",
	}

	err := bus.Dispatch(ctx, authURL)
	Expect(err).IsNil()
	Expect(authURL.Result).Equals("https://example.org/oauth/authorize?client_id=CU_CL_ID&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2F_custom%2Fcallback&response_type=code&scope=profile+email&state=http%3A%2F%2Fexample.org%7C456")
}

func TestGetAuthURL_Twitch(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_custom" {
			q.Result = &models.OAuthConfig{
				Provider:     q.Provider,
				ClientID:     "CU_CL_ID",
				Scope:        "openid",
				AuthorizeURL: "https://id.twitch.tv/oauth/authorize",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	authURL := &query.GetOAuthAuthorizationURL{
		Provider:   "_custom",
		Redirect:   "http://example.org",
		Identifier: "456",
	}

	err := bus.Dispatch(ctx, authURL)
	Expect(err).IsNil()
	Expect(authURL.Result).Equals("https://id.twitch.tv/oauth/authorize?claims=%7B%22userinfo%22%3A%7B%22preferred_username%22%3Anull%2C%22email%22%3Anull%2C%22email_verified%22%3Anull%7D%7D&client_id=CU_CL_ID&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2F_custom%2Fcallback&response_type=code&scope=openid&state=http%3A%2F%2Fexample.org%7C456")
}

func TestParseProfileResponse_AllFields(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_test1" {
			q.Result = &models.OAuthConfig{
				Provider:          q.Provider,
				JSONUserIDPath:    "id",
				JSONUserNamePath:  "name",
				JSONUserEmailPath: "email",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	profile := &cmd.ParseOAuthRawProfile{
		Provider: "_test1",
		Body:     `{"name":"Jon Snow","email":"jon\u0040got.com","id":"789654"}`,
	}

	err := bus.Dispatch(ctx, profile)
	Expect(err).IsNil()

	Expect(profile.Result.ID).Equals("789654")
	Expect(profile.Result.Name).Equals("Jon Snow")
	Expect(profile.Result.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_WithoutEmail(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_test1" {
			q.Result = &models.OAuthConfig{
				Provider:          q.Provider,
				JSONUserIDPath:    "id",
				JSONUserNamePath:  "name",
				JSONUserEmailPath: "email",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	profile := &cmd.ParseOAuthRawProfile{
		Provider: "_test1",
		Body:     `{"name":"Jon Snow","id":"1"}`,
	}
	err := bus.Dispatch(ctx, profile)
	Expect(err).IsNil()

	Expect(profile.Result.ID).Equals("1")
	Expect(profile.Result.Name).Equals("Jon Snow")
	Expect(profile.Result.Email).Equals("")
}

func TestParseProfileResponse_NestedData(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_test1" {
			q.Result = &models.OAuthConfig{
				Provider:          q.Provider,
				JSONUserIDPath:    "id",
				JSONUserNamePath:  "profile.name",
				JSONUserEmailPath: "profile.email",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	profile := &cmd.ParseOAuthRawProfile{
		Provider: "_test1",
		Body: `{
			"id": "123",
			"profile": {
				"name": "Jon Snow",
				"email": "Jon@got.com"
			}
		}`,
	}
	err := bus.Dispatch(ctx, profile)
	Expect(err).IsNil()

	Expect(profile.Result.ID).Equals("123")
	Expect(profile.Result.Name).Equals("Jon Snow")
	Expect(profile.Result.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_WithFallback(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_test1" {
			q.Result = &models.OAuthConfig{
				Provider:          q.Provider,
				JSONUserIDPath:    "id",
				JSONUserNamePath:  "profile.name, profile.login",
				JSONUserEmailPath: "profile.email",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	profile := &cmd.ParseOAuthRawProfile{
		Provider: "_test1",
		Body: `{
			"id": 123,
			"profile": {
				"login": "jonny",
				"email": "jon@got.com"
			}
		}`,
	}
	err := bus.Dispatch(ctx, profile)
	Expect(err).IsNil()

	Expect(profile.Result.ID).Equals("123")
	Expect(profile.Result.Name).Equals("jonny")
	Expect(profile.Result.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_UseEmailWhenNameIsEmpty(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_test1" {
			q.Result = &models.OAuthConfig{
				Provider:          q.Provider,
				JSONUserIDPath:    "id",
				JSONUserNamePath:  "profile.name",
				JSONUserEmailPath: "profile.email",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	profile := &cmd.ParseOAuthRawProfile{
		Provider: "_test1",
		Body: `{
			"id": "123",
			"profile": {
				"email": "jon@got.com"
			}
		}`,
	}
	err := bus.Dispatch(ctx, profile)
	Expect(err).IsNil()

	Expect(profile.Result.ID).Equals("123")
	Expect(profile.Result.Name).Equals("jon")
	Expect(profile.Result.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_InvalidEmail(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_test1" {
			q.Result = &models.OAuthConfig{
				Provider:          q.Provider,
				JSONUserIDPath:    "id",
				JSONUserNamePath:  "profile.name",
				JSONUserEmailPath: "profile.email",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	profile := &cmd.ParseOAuthRawProfile{
		Provider: "_test1",
		Body: `{
			"id": "AB123",
			"profile": {
				"name": "Jon Snow",
				"email": "jon"
			}
		}`,
	}
	err := bus.Dispatch(ctx, profile)
	Expect(err).IsNil()

	Expect(profile.Result.ID).Equals("AB123")
	Expect(profile.Result.Name).Equals("Jon Snow")
	Expect(profile.Result.Email).Equals("")
}

func TestParseProfileResponse_EmptyID(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_test1" {
			q.Result = &models.OAuthConfig{
				Provider:          q.Provider,
				JSONUserIDPath:    "id",
				JSONUserNamePath:  "name",
				JSONUserEmailPath: "email",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	profile := &cmd.ParseOAuthRawProfile{
		Provider: "_test1",
		Body:     `{}`,
	}

	err := bus.Dispatch(ctx, profile)
	Expect(errors.Cause(err)).Equals(app.ErrUserIDRequired)
	Expect(profile.Result).IsNil()
}

func TestParseProfileResponse_EmptyName(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_test1" {
			q.Result = &models.OAuthConfig{
				Provider:          q.Provider,
				JSONUserIDPath:    "id",
				JSONUserNamePath:  "name",
				JSONUserEmailPath: "email",
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	profile := &cmd.ParseOAuthRawProfile{
		Provider: "_test1",
		Body:     `{ "id": "A0" }`,
	}

	err := bus.Dispatch(ctx, profile)
	Expect(err).IsNil()
	Expect(profile.Result.ID).Equals("A0")
	Expect(profile.Result.Name).Equals("Anonymous")
}

func TestCustomOAuth_Disabled(t *testing.T) {
	RegisterT(t)
	bus.Init(&oauth.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_test1" {
			q.Result = &models.OAuthConfig{
				Provider:          q.Provider,
				JSONUserIDPath:    "id",
				JSONUserNamePath:  "name",
				JSONUserEmailPath: "email",
				Status:            enum.OAuthConfigDisabled,
			}
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	rawProfile := &cmd.ParseOAuthRawProfile{
		Provider: "_test1",
		Body:     `{ "id": "A0", "name": "John" }`,
	}

	err := bus.Dispatch(ctx, rawProfile)
	Expect(err).IsNil()
	Expect(rawProfile.Result.ID).Equals("A0")
	Expect(rawProfile.Result.Name).Equals("John")

	oauthProfile := &query.GetOAuthProfile{
		Provider: "_test1",
	}

	err = bus.Dispatch(ctx, oauthProfile)
	Expect(err).IsNotNil()
	Expect(oauthProfile.Result).IsNil()
}
