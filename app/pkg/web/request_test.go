package web_test

import (
	"crypto/tls"
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
			Method:     "GET",
			Header:     header,
			RemoteAddr: "172.10.10.10:5555",
			Host:       "helloworld.com",
		},
	)

	Expect(req.Method).Equals("GET")
	Expect(req.GetHeader("Content-Type")).Equals("application/json")
	Expect(req.ClientIP).Equals("172.10.10.10")
	Expect(req.URL.Hostname()).Equals("helloworld.com")
	Expect(req.URL.Scheme).Equals("http")
	Expect(req.URL.RequestURI()).Equals("/")
	Expect(req.URL.String()).Equals("http://helloworld.com")
	Expect(req.IsSecure).Equals(false)
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
	Expect(req.ClientIP).Equals("")
	Expect(req.URL.Hostname()).Equals("helloworld.com")
	Expect(req.URL.Scheme).Equals("http")
	Expect(req.URL.Port()).Equals("3000")
	Expect(req.URL.RequestURI()).Equals("/echo")
	Expect(req.URL.String()).Equals("http://helloworld.com:3000/echo")
}

func TestRequest_BehindTLSTerminationProxy(t *testing.T) {
	RegisterT(t)

	header := make(http.Header, 0)
	header.Set("X-Forwarded-Host", "feedback.mycompany.com")
	header.Set("X-Forwarded-Proto", "https")
	header.Set("X-Forwarded-For", "127.5.5.5, 129.2.2.2, 121.2.2.5")

	req := web.WrapRequest(
		&http.Request{
			Method: "GET",
			Header: header,
			Host:   "demo.test.fider.io",
		},
	)

	Expect(req.Method).Equals("GET")
	Expect(req.URL.Hostname()).Equals("feedback.mycompany.com")
	Expect(req.URL.Scheme).Equals("https")
	Expect(req.ClientIP).Equals("127.5.5.5")
	Expect(req.IsSecure).Equals(true)
	Expect(req.IsAPI()).IsFalse()
}

func TestRequest_FullURL(t *testing.T) {
	RegisterT(t)

	req := web.WrapRequest(
		&http.Request{
			TLS:        &tls.ConnectionState{},
			Host:       "demo.test.fider.io",
			RequestURI: "/api/hello?value=Jon",
		},
	)

	Expect(req.URL.String()).Equals("https://demo.test.fider.io/api/hello?value=Jon")
	Expect(req.URL.Path).Equals("/api/hello")
	Expect(req.URL.Query().Get("value")).Equals("Jon")
	Expect(req.URL.RequestURI()).Equals("/api/hello?value=Jon")
	Expect(req.IsSecure).Equals(true)
	Expect(req.IsAPI()).IsTrue()
}
