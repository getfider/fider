package web_test

import (
	"testing"

	"github.com/getfider/fider/app/models"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
)

func TestGetAuthURL_Facebook(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	authURL, err := svc.GetAuthURL(oauth.FacebookProvider, "http://example.org", "456")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://www.facebook.com/dialog/oauth?client_id=&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Ffacebook%2Fcallback&response_type=code&scope=public_profile+email&state=http%3A%2F%2Fexample.org%7C456")
}

func TestGetAuthURL_Google(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	authURL, err := svc.GetAuthURL(oauth.GoogleProvider, "http://example.org", "123")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://accounts.google.com/o/oauth2/auth?client_id=&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Fgoogle%2Fcallback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&state=http%3A%2F%2Fexample.org%7C123")
}

func TestGetAuthURL_GitHub(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	authURL, err := svc.GetAuthURL(oauth.GitHubProvider, "http://example.org", "456")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://github.com/login/oauth/authorize?client_id=&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Fgithub%2Fcallback&response_type=code&scope=user%3Aemail&state=http%3A%2F%2Fexample.org%7C456")
}

func TestParseProfileResponse_AllFields(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	profile, err := svc.ParseProfileResponse(
		`{"name":"Jon Snow","email":"jon\u0040got.com","id":"789654"}`,
		&models.OAuthConfig{
			JSONUserIDPath: "id",
			JSONNamePath:   "name",
			JSONEmailPath:  "email",
		},
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("789654")
	Expect(profile.Name).Equals("Jon Snow")
	Expect(profile.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_WithoutEmail(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	profile, err := svc.ParseProfileResponse(
		`{"name":"Jon Snow","id":"1"}`,
		&models.OAuthConfig{
			JSONUserIDPath: "id",
			JSONNamePath:   "name",
			JSONEmailPath:  "email",
		},
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("1")
	Expect(profile.Name).Equals("Jon Snow")
	Expect(profile.Email).Equals("")
}

func TestParseProfileResponse_NestedData(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	profile, err := svc.ParseProfileResponse(
		`{
			"id": "123",
			"profile": {
				"name": "Jon Snow",
				"email": "jon@got.com"
			}
		}`,
		&models.OAuthConfig{
			JSONUserIDPath: "id",
			JSONNamePath:   "profile.name",
			JSONEmailPath:  "profile.email",
		},
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("123")
	Expect(profile.Name).Equals("Jon Snow")
	Expect(profile.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_WithFallback(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	profile, err := svc.ParseProfileResponse(
		`{
			"id": 123,
			"profile": {
				"login": "jonny",
				"email": "jon@got.com"
			}
		}`,
		&models.OAuthConfig{
			JSONUserIDPath: "id",
			JSONNamePath:   "profile.name, profile.login",
			JSONEmailPath:  "profile.email",
		},
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("123")
	Expect(profile.Name).Equals("jonny")
	Expect(profile.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_UseEmailWhenNameIsEmpty(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	profile, err := svc.ParseProfileResponse(
		`{
			"id": "123",
			"profile": {
				"email": "jon@got.com"
			}
		}`,
		&models.OAuthConfig{
			JSONUserIDPath: "id",
			JSONNamePath:   "profile.name",
			JSONEmailPath:  "profile.email",
		},
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("123")
	Expect(profile.Name).Equals("jon")
	Expect(profile.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_InvalidEmail(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	profile, err := svc.ParseProfileResponse(
		`{
			"id": "AB123",
			"profile": {
				"name": "Jon Snow",
				"email": "jon"
			}
		}`,
		&models.OAuthConfig{
			JSONUserIDPath: "id",
			JSONNamePath:   "profile.name",
			JSONEmailPath:  "profile.email",
		},
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("AB123")
	Expect(profile.Name).Equals("Jon Snow")
	Expect(profile.Email).Equals("")
}

func TestParseProfileResponse_EmptyID(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	profile, err := svc.ParseProfileResponse(
		`{}`,
		&models.OAuthConfig{
			JSONUserIDPath: "id",
			JSONNamePath:   "name",
			JSONEmailPath:  "email",
		},
	)

	Expect(errors.Cause(err)).Equals(oauth.ErrUserIDRequired)
	Expect(profile).IsNil()
}

func TestParseProfileResponse_EmptyName(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000")
	profile, err := svc.ParseProfileResponse(
		`{ "id": "A0" }`,
		&models.OAuthConfig{
			JSONUserIDPath: "id",
			JSONNamePath:   "name",
			JSONEmailPath:  "email",
		},
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("A0")
	Expect(profile.Name).Equals("Anonymous")
}
