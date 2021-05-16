package web_test

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/entity"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

func newGetContext(rawurl string, headers map[string]string) *web.Context {
	u, _ := url.Parse(rawurl)
	e := web.New()
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", u.RequestURI(), nil)
	req.Host = u.Host

	if u.Scheme == "https" {
		req.TLS = &tls.ConnectionState{}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return web.NewContext(e, req, res, nil)
}

func newBodyContext(method string, params web.StringMap, body, contentType string) *web.Context {
	e := web.New()
	res := httptest.NewRecorder()
	req := httptest.NewRequest(method, "/some/resource", strings.NewReader(body))
	req.Host = "demo.test.fider.io:3000"
	req.Header.Set("Content-Type", contentType)
	return web.NewContext(e, req, res, params)
}

func TestContextID(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000", nil)

	Expect(ctx.ContextID()).IsNotEmpty()
	Expect(ctx.ContextID()).HasLen(32)
}

func TestBaseURL(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000", nil)

	Expect(ctx.BaseURL()).Equals("http://demo.test.fider.io:3000")
}

func TestBaseURL_HTTPS(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("https://demo.test.fider.io:3000", nil)

	Expect(ctx.BaseURL()).Equals("https://demo.test.fider.io:3000")
}

func TestBaseURL_HTTPS_Proxy(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000", map[string]string{
		"X-Forwarded-Proto": "https",
	})

	Expect(ctx.BaseURL()).Equals("https://demo.test.fider.io:3000")
}

func TestCurrentURL(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000/resource?id=23", nil)

	Expect(ctx.Request.URL.String()).Equals("http://demo.test.fider.io:3000/resource?id=23")
}

func TestTenantURL(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://login.test.fider.io:3000", nil)
	tenant := &entity.Tenant{
		ID:        1,
		Subdomain: "theavengers",
	}
	Expect(web.TenantBaseURL(ctx, tenant)).Equals("http://theavengers.test.fider.io:3000")
}

func TestTenantURL_WithCNAME(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000", nil)
	tenant := &entity.Tenant{
		ID:        1,
		Subdomain: "theavengers",
		CNAME:     "feedback.theavengers.com",
	}
	Expect(web.TenantBaseURL(ctx, tenant)).Equals("http://feedback.theavengers.com:3000")
}

func TestTenantURL_SingleHostMode(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "single"

	ctx := newGetContext("http://demo.test.fider.io:3000", nil)
	tenant := &entity.Tenant{
		ID:        1,
		Subdomain: "theavengers",
	}
	Expect(web.TenantBaseURL(ctx, tenant)).Equals("http://demo.test.fider.io:3000")
}

func TestAssetsURL_SingleHostMode(t *testing.T) {
	RegisterT(t)

	env.Config.HostMode = "single"
	ctx := newGetContext("http://feedback.theavengers.com:3000", nil)
	ctx.SetTenant(&entity.Tenant{
		ID:        1,
		Subdomain: "theavengers",
	})

	Expect(web.AssetsURL(ctx, "/assets/main.js")).Equals("http://feedback.theavengers.com:3000/assets/main.js")
	Expect(web.AssetsURL(ctx, "/assets/main.css")).Equals("http://feedback.theavengers.com:3000/assets/main.css")

	env.Config.CDN.Host = "fidercdn.com"
	Expect(web.AssetsURL(ctx, "/assets/main.js")).Equals("http://fidercdn.com/assets/main.js")
	Expect(web.AssetsURL(ctx, "/assets/main.css")).Equals("http://fidercdn.com/assets/main.css")
}

func TestAssetsURL_MultiHostMode(t *testing.T) {
	RegisterT(t)

	env.Config.HostMode = "multi"
	ctx := newGetContext("http://theavengers.test.fider.io:3000", nil)
	ctx.SetTenant(&entity.Tenant{
		ID:        1,
		Subdomain: "theavengers",
		CNAME:     "feedback.theavengers.com",
	})

	Expect(web.AssetsURL(ctx, "/assets/main.js")).Equals("http://theavengers.test.fider.io:3000/assets/main.js")
	Expect(web.AssetsURL(ctx, "/assets/main.css")).Equals("http://theavengers.test.fider.io:3000/assets/main.css")

	env.Config.CDN.Host = "fidercdn.com"
	Expect(web.AssetsURL(ctx, "/assets/main.js")).Equals("http://theavengers.fidercdn.com/assets/main.js")
	Expect(web.AssetsURL(ctx, "/assets/main.css")).Equals("http://theavengers.fidercdn.com/assets/main.css")
}

func TestCanonicalURL_SameDomain(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://theavengers.test.fider.io:3000", nil)

	ctx.SetCanonicalURL("")
	Expect(ctx.Value("Canonical-URL")).Equals(`http://theavengers.test.fider.io:3000`)

	ctx.SetCanonicalURL("/some-url")
	Expect(ctx.Value("Canonical-URL")).Equals(`http://theavengers.test.fider.io:3000/some-url`)

	ctx.SetCanonicalURL("/some-other-url")
	Expect(ctx.Value("Canonical-URL")).Equals(`http://theavengers.test.fider.io:3000/some-other-url`)

	ctx.SetCanonicalURL("page-b/abc.html")
	Expect(ctx.Value("Canonical-URL")).Equals(`http://theavengers.test.fider.io:3000/page-b/abc.html`)
}

func TestCanonicalURL_DifferentDomain(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://theavengers.test.fider.io:3000", nil)

	ctx.SetCanonicalURL("http://feedback.theavengers.com/some-url")
	Expect(ctx.Value("Canonical-URL")).Equals(`http://feedback.theavengers.com/some-url`)

	ctx.SetCanonicalURL("")
	Expect(ctx.Value("Canonical-URL")).Equals(`http://feedback.theavengers.com`)

	ctx.SetCanonicalURL("/some-other-url")
	Expect(ctx.Value("Canonical-URL")).Equals(`http://feedback.theavengers.com/some-other-url`)

	ctx.SetCanonicalURL("page-b/abc.html")
	Expect(ctx.Value("Canonical-URL")).Equals(`http://feedback.theavengers.com/page-b/abc.html`)
}

func TestTryAgainLater(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000", nil)
	err := ctx.TryAgainLater(24 * time.Hour)
	Expect(err).IsNil()
	resp := ctx.Response.(*httptest.ResponseRecorder)
	Expect(ctx.ResponseStatusCode).Equals(http.StatusServiceUnavailable)
	Expect(resp.Code).Equals(http.StatusServiceUnavailable)
	Expect(resp.Header().Get("Cache-Control")).Equals("no-cache, no-store")
	Expect(resp.Header().Get("Retry-After")).Equals("86400")
}

func TestGetOAuthBaseURL(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("https://mydomain.com/hello-world", nil)

	env.Config.HostMode = "multi"
	Expect(web.OAuthBaseURL(ctx)).Equals("https://login.test.fider.io")

	env.Config.HostMode = "single"
	Expect(web.OAuthBaseURL(ctx)).Equals("https://mydomain.com")
}

func TestGetOAuthBaseURL_WithPort(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000/hello-world", nil)

	env.Config.HostMode = "multi"
	Expect(web.OAuthBaseURL(ctx)).Equals("http://login.test.fider.io:3000")

	env.Config.HostMode = "single"
	Expect(web.OAuthBaseURL(ctx)).Equals("http://demo.test.fider.io:3000")
}
