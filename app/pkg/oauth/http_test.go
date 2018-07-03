package oauth_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/oauth"
)

func TestGetAuthURL_Facebook(t *testing.T) {
	RegisterT(t)

	svc := oauth.NewHTTPService("http://login.test.fider.io:3000")
	authURL, err := svc.GetAuthURL(oauth.FacebookProvider, "")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://www.facebook.com/dialog/oauth?client_id=&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Ffacebook%2Fcallback&response_type=code&scope=public_profile+email&state=")
}

func TestGetAuthURL_Google(t *testing.T) {
	RegisterT(t)

	svc := oauth.NewHTTPService("http://login.test.fider.io:3000")
	authURL, err := svc.GetAuthURL(oauth.GoogleProvider, "")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://accounts.google.com/o/oauth2/auth?client_id=&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Fgoogle%2Fcallback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&state=")
}

func TestGetAuthURL_GitHub(t *testing.T) {
	RegisterT(t)

	svc := oauth.NewHTTPService("http://login.test.fider.io:3000")
	authURL, err := svc.GetAuthURL(oauth.GitHubProvider, "")

	Expect(err).IsNil()
	Expect(authURL).Equals("https://github.com/login/oauth/authorize?client_id=&redirect_uri=http%3A%2F%2Flogin.test.fider.io%3A3000%2Foauth%2Fgithub%2Fcallback&response_type=code&scope=user%3Aemail&state=")
}
