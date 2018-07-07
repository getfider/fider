package webutil_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/getfider/fider/app/pkg/web/util"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func newContext(rawurl string) web.Context {
	url, _ := url.Parse(rawurl)

	return web.Context{
		Request: web.Request{
			URL: url,
		},
	}
}

func TestGetOAuthBaseURL(t *testing.T) {
	RegisterT(t)

	ctx := newContext("https://mydomain.com/hello-world")

	os.Setenv("HOST_MODE", "multi")
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("https://login.test.fider.io")

	os.Setenv("HOST_MODE", "single")
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("https://mydomain.com")
}

func TestGetOAuthBaseURL_WithPort(t *testing.T) {
	RegisterT(t)

	ctx := newContext("http://demo.test.fider.io:3000/hello-world")

	os.Setenv("HOST_MODE", "multi")
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("http://login.test.fider.io:3000")

	os.Setenv("HOST_MODE", "single")
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("http://demo.test.fider.io:3000")
}
