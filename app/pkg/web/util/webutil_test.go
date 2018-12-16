package webutil_test

import (
	"net/url"
	"testing"

	"github.com/getfider/fider/app/pkg/env"
	webutil "github.com/getfider/fider/app/pkg/web/util"

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

	env.Config.HostMode = "multi"
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("https://login.test.fider.io")

	env.Config.HostMode = "single"
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("https://mydomain.com")
}

func TestGetOAuthBaseURL_WithPort(t *testing.T) {
	RegisterT(t)

	ctx := newContext("http://demo.test.fider.io:3000/hello-world")

	env.Config.HostMode = "multi"
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("http://login.test.fider.io:3000")

	env.Config.HostMode = "single"
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("http://demo.test.fider.io:3000")
}
