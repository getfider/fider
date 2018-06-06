package web_test

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func newGetContext(params web.StringMap) *web.Context {
	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/some/resource", nil)
	req.Host = "demo.test.fider.io:3000"
	ctx := e.NewContext(res, req, params)
	return &ctx
}

func newBodyContext(method string, params web.StringMap, body, contentType string) *web.Context {
	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(method, "/some/resource", strings.NewReader(body))
	req.Host = "demo.test.fider.io:3000"
	req.Header.Set("Content-Type", contentType)
	ctx := e.NewContext(res, req, params)
	return &ctx
}

func TestContextID(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext(nil)

	Expect(ctx.ContextID()).IsNotEmpty()
	Expect(ctx.ContextID()).HasLen(32)
}

func TestBaseURL(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext(nil)

	Expect(ctx.BaseURL()).Equals("http://demo.test.fider.io:3000")
}

func TestTenantURL(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext(nil)
	tenant := &models.Tenant{
		ID:        1,
		Subdomain: "theavengers",
	}
	Expect(ctx.TenantBaseURL(tenant)).Equals("http://theavengers.test.fider.io:3000")
}

func TestTenantURL_WithCNAME(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext(nil)
	tenant := &models.Tenant{
		ID:        1,
		Subdomain: "theavengers",
		CNAME:     "ideas.theavengers.com",
	}
	Expect(ctx.TenantBaseURL(tenant)).Equals("http://ideas.theavengers.com:3000")
}

func TestTenantURL_SingleTenantMode(t *testing.T) {
	RegisterT(t)
	os.Setenv("HOST_MODE", "single")

	ctx := newGetContext(nil)
	tenant := &models.Tenant{
		ID:        1,
		Subdomain: "theavengers",
	}
	Expect(ctx.TenantBaseURL(tenant)).Equals("http://demo.test.fider.io:3000")
}
