package web_test

import (
	"testing"

	"github.com/getfider/fider/app/models"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage/inmemory"
)

func TestGetAuthURL_Facebook(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000", inmemory.NewTenantStorage())
	authURL, err := svc.GetAuthURL(oauth.FacebookProvider, "http://example.org", "456")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://www.facebook.com/dialog/oauth?client_id=FB_CL_ID&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Ffacebook%2Fcallback&response_type=code&scope=public_profile+email&state=http%3A%2F%2Fexample.org%7C456")
}

func TestGetAuthURL_Google(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000", inmemory.NewTenantStorage())
	authURL, err := svc.GetAuthURL(oauth.GoogleProvider, "http://example.org", "123")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://accounts.google.com/o/oauth2/auth?client_id=GO_CL_ID&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Fgoogle%2Fcallback&response_type=code&scope=profile+email&state=http%3A%2F%2Fexample.org%7C123")
}

func TestGetAuthURL_GitHub(t *testing.T) {
	RegisterT(t)

	svc := web.NewOAuthService("http://login.test.fider.io:3000", inmemory.NewTenantStorage())
	authURL, err := svc.GetAuthURL(oauth.GitHubProvider, "http://example.org", "456")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://github.com/login/oauth/authorize?client_id=GH_CL_ID&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Fgithub%2Fcallback&response_type=code&scope=user%3Aemail&state=http%3A%2F%2Fexample.org%7C456")
}

func TestGetAuthURL_Custom(t *testing.T) {
	RegisterT(t)

	tenants := inmemory.NewTenantStorage()
	tenants.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		Provider:     "_custom",
		ClientID:     "CU_CL_ID",
		Scope:        "profile email",
		AuthorizeURL: "https://example.org/oauth/authorize",
	})
	svc := web.NewOAuthService("http://login.test.fider.io:3000", tenants)
	authURL, err := svc.GetAuthURL("_custom", "http://example.org", "456")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://example.org/oauth/authorize?client_id=CU_CL_ID&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2F_custom%2Fcallback&response_type=code&scope=profile+email&state=http%3A%2F%2Fexample.org%7C456")
}

func TestParseProfileResponse_AllFields(t *testing.T) {
	RegisterT(t)

	storage := inmemory.NewTenantStorage()
	storage.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		Provider:          "_test1",
		JSONUserIDPath:    "id",
		JSONUserNamePath:  "name",
		JSONUserEmailPath: "email",
	})
	svc := web.NewOAuthService("http://login.test.fider.io:3000", storage)
	profile, err := svc.ParseRawProfile(
		"_test1",
		`{"name":"Jon Snow","email":"jon\u0040got.com","id":"789654"}`,
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("789654")
	Expect(profile.Name).Equals("Jon Snow")
	Expect(profile.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_WithoutEmail(t *testing.T) {
	RegisterT(t)

	storage := inmemory.NewTenantStorage()
	storage.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		Provider:          "_test1",
		JSONUserIDPath:    "id",
		JSONUserNamePath:  "name",
		JSONUserEmailPath: "email",
	})
	svc := web.NewOAuthService("http://login.test.fider.io:3000", storage)
	profile, err := svc.ParseRawProfile(
		"_test1",
		`{"name":"Jon Snow","id":"1"}`,
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("1")
	Expect(profile.Name).Equals("Jon Snow")
	Expect(profile.Email).Equals("")
}

func TestParseProfileResponse_NestedData(t *testing.T) {
	RegisterT(t)

	storage := inmemory.NewTenantStorage()
	storage.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		Provider:          "_test1",
		JSONUserIDPath:    "id",
		JSONUserNamePath:  "profile.name",
		JSONUserEmailPath: "profile.email",
	})
	svc := web.NewOAuthService("http://login.test.fider.io:3000", storage)
	profile, err := svc.ParseRawProfile(
		"_test1",
		`{
			"id": "123",
			"profile": {
				"name": "Jon Snow",
				"email": "jon@got.com"
			}
		}`,
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("123")
	Expect(profile.Name).Equals("Jon Snow")
	Expect(profile.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_WithFallback(t *testing.T) {
	RegisterT(t)

	storage := inmemory.NewTenantStorage()
	storage.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		Provider:          "_test1",
		JSONUserIDPath:    "id",
		JSONUserNamePath:  "profile.name, profile.login",
		JSONUserEmailPath: "profile.email",
	})
	svc := web.NewOAuthService("http://login.test.fider.io:3000", storage)
	profile, err := svc.ParseRawProfile(
		"_test1",
		`{
			"id": 123,
			"profile": {
				"login": "jonny",
				"email": "jon@got.com"
			}
		}`,
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("123")
	Expect(profile.Name).Equals("jonny")
	Expect(profile.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_UseEmailWhenNameIsEmpty(t *testing.T) {
	RegisterT(t)

	storage := inmemory.NewTenantStorage()
	storage.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		Provider:          "_test1",
		JSONUserIDPath:    "id",
		JSONUserNamePath:  "profile.name",
		JSONUserEmailPath: "profile.email",
	})
	svc := web.NewOAuthService("http://login.test.fider.io:3000", storage)
	profile, err := svc.ParseRawProfile(
		"_test1",
		`{
			"id": "123",
			"profile": {
				"email": "jon@got.com"
			}
		}`,
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("123")
	Expect(profile.Name).Equals("jon")
	Expect(profile.Email).Equals("jon@got.com")
}

func TestParseProfileResponse_InvalidEmail(t *testing.T) {
	RegisterT(t)

	storage := inmemory.NewTenantStorage()
	storage.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		Provider:          "_test1",
		JSONUserIDPath:    "id",
		JSONUserNamePath:  "profile.name",
		JSONUserEmailPath: "profile.email",
	})
	svc := web.NewOAuthService("http://login.test.fider.io:3000", storage)
	profile, err := svc.ParseRawProfile(
		"_test1",
		`{
			"id": "AB123",
			"profile": {
				"name": "Jon Snow",
				"email": "jon"
			}
		}`,
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("AB123")
	Expect(profile.Name).Equals("Jon Snow")
	Expect(profile.Email).Equals("")
}

func TestParseProfileResponse_EmptyID(t *testing.T) {
	RegisterT(t)

	storage := inmemory.NewTenantStorage()
	storage.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		Provider:          "_test1",
		JSONUserIDPath:    "id",
		JSONUserNamePath:  "name",
		JSONUserEmailPath: "email",
	})
	svc := web.NewOAuthService("http://login.test.fider.io:3000", storage)
	profile, err := svc.ParseRawProfile(
		"_test1",
		`{}`,
	)

	Expect(errors.Cause(err)).Equals(oauth.ErrUserIDRequired)
	Expect(profile).IsNil()
}

func TestParseProfileResponse_EmptyName(t *testing.T) {
	RegisterT(t)

	storage := inmemory.NewTenantStorage()
	storage.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		Provider:          "_test1",
		JSONUserIDPath:    "id",
		JSONUserNamePath:  "name",
		JSONUserEmailPath: "email",
	})
	svc := web.NewOAuthService("http://login.test.fider.io:3000", storage)
	profile, err := svc.ParseRawProfile(
		"_test1",
		`{ "id": "A0" }`,
	)

	Expect(err).IsNil()
	Expect(profile.ID).Equals("A0")
	Expect(profile.Name).Equals("Anonymous")
}
