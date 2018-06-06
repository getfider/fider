package web_test

import (
	"net/http"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func TestRequest_Basic(t *testing.T) {
	RegisterT(t)

	header := make(http.Header, 0)
	header.Set("Content-Type", "application/json")

	req := web.WrapRequest(
		&http.Request{
			Method: "GET",
			Header: header,
			Host:   "helloworld.com",
		},
	)

	Expect(req.Method).Equals("GET")
	Expect(req.GetHeader("Content-Type")).Equals("application/json")
	Expect(req.Host).Equals("helloworld.com")
	Expect(req.Protocol).Equals("http")
	Expect(req.URL).Equals("/")
	Expect(req.FullURL).Equals("http://helloworld.com")
}

func TestRequest_WithPort(t *testing.T) {
	RegisterT(t)

	header := make(http.Header, 0)
	header.Set("Content-Type", "application/json")

	req := web.WrapRequest(
		&http.Request{
			Method:     "GET",
			Header:     header,
			Host:       "helloworld.com:3000",
			RequestURI: "/echo",
		},
	)

	Expect(req.Method).Equals("GET")
	Expect(req.GetHeader("Content-Type")).Equals("application/json")
	Expect(req.Host).Equals("helloworld.com")
	Expect(req.Protocol).Equals("http")
	Expect(req.Port).Equals("3000")
	Expect(req.URL).Equals("/echo")
	Expect(req.FullURL).Equals("http://helloworld.com:3000/echo")
}

func TestRequest_BehindTLSTerminationProxy(t *testing.T) {
	RegisterT(t)

	header := make(http.Header, 0)
	header.Set("X-Forwarded-Host", "feedback.mycompany.com")
	header.Set("X-Forwarded-Proto", "https")

	req := web.WrapRequest(
		&http.Request{
			Method: "GET",
			Header: header,
			Host:   "demo.test.fider.io",
		},
	)

	Expect(req.Method).Equals("GET")
	Expect(req.Host).Equals("feedback.mycompany.com")
	Expect(req.Subdomain).Equals("")
	Expect(req.Protocol).Equals("https")
}

func TestRequest_FullURL(t *testing.T) {
	RegisterT(t)

	req := web.WrapRequest(
		&http.Request{
			Host:       "demo.test.fider.io",
			RequestURI: "/echo?value=Jon",
		},
	)

	Expect(req.FullURL).Equals("http://demo.test.fider.io/echo?value=Jon")
	Expect(req.Path).Equals("/echo")
	Expect(req.Subdomain).Equals("demo")
	Expect(req.Query.Get("value")).Equals("Jon")
	Expect(req.URL).Equals("/echo?value=Jon")
}
