package web_test

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

func newGetContext(params web.StringMap) *web.Context {
	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://demo.test.fider.io:3000/some/resource", nil)
	ctx := e.NewContext(res, req, params)
	return &ctx
}

func newPostContext(params web.StringMap, body string) *web.Context {
	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "http://demo.test.fider.io:3000/some/resource", strings.NewReader(body))
	ctx := e.NewContext(res, req, params)
	return &ctx
}

func TestBaseURL(t *testing.T) {
	RegisterTestingT(t)

	ctx := newGetContext(nil)

	Expect(ctx.BaseURL()).To(Equal("http://demo.test.fider.io:3000"))
}

func TestTenantURL(t *testing.T) {
	RegisterTestingT(t)

	ctx := newGetContext(nil)
	tenant := &models.Tenant{
		ID:        1,
		Subdomain: "lotr",
	}
	Expect(ctx.TenantBaseURL(tenant)).To(Equal("http://lotr.test.fider.io:3000"))
}

func TestTenantURL_SingleTenantMode(t *testing.T) {
	RegisterTestingT(t)
	os.Setenv("HOST_MODE", "single")

	ctx := newGetContext(nil)
	tenant := &models.Tenant{
		ID:        1,
		Subdomain: "lotr",
	}
	Expect(ctx.TenantBaseURL(tenant)).To(Equal("http://demo.test.fider.io:3000"))
}
