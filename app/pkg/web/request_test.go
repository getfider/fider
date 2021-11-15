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

	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	req := web.WrapRequest(
		&http.Request{
			Method:     "GET",
			Header:     header,
			Host:       "helloworld.com",
		},
	)

	Expect(req.Method).Equals("GET")
	Expect(req.GetHeader("Content-Type")).Equals("application/json")
	Expect(req.URL.Hostname()).Equals("helloworld.com")
	Expect(req.URL.Scheme).Equals("http")
	Expect(req.URL.RequestURI()).Equals("/")
	Expect(req.URL.String()).Equals("http://helloworld.com")
	Expect(req.IsSecure).Equals(false)
}

func TestRequest_WithPort(t *testing.T) {
	RegisterT(t)

	header := make(http.Header)
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
	Expect(req.URL.Hostname()).Equals("helloworld.com")
	Expect(req.URL.Scheme).Equals("http")
	Expect(req.URL.Port()).Equals("3000")
	Expect(req.URL.RequestURI()).Equals("/echo")
	Expect(req.URL.String()).Equals("http://helloworld.com:3000/echo")
}

func TestRequest_BehindTLSTerminationProxy(t *testing.T) {
	RegisterT(t)

	header := make(http.Header)
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
	Expect(req.URL.Hostname()).Equals("feedback.mycompany.com")
	Expect(req.URL.Scheme).Equals("https")
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

func TestRequest_IsCrawler(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		userAgent string
		isCrawler bool
	}{
		{"Baidu Union Search", true},
		{"msnbot/2.0b (+http://search.msn.com/msnbot.htm)", true},
		{"Mozilla/5.0 (compatible; bingbot/2.0 +http://www.bing.com/bingbot.htm)", true},
		{"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534+ (KHTML, like Gecko) BingPreview/1.0b", true},
		{"DuckDuckBot/1.0; (+http://duckduckgo.com/duckduckbot.html)", true},
		{"Googlebot/2.1 (+http://www.googlebot.com/bot.html)", true},
		{"AdsBot-Google (+http://www.google.com/adsbot.html)", true},
		{"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)", true},
		{"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)", true},
		{"Mozilla/5.0 (compatible; YandexMetrika/3.0; +http://yandex.com/bots)", true},
		{"Mozilla/5.0 (compatible; AhrefsBot/6.1; +http://ahrefs.com/robot/)", true},
		{"Googlebot-News", true},
		{"Twitterbot/1.0", true},
		{"Slackbot 1.0(+https://api.slack.com/robots)", true},
		{"Slackbot-LinkExpanding 1.0 (+https://api.slack.com/robots)", true},
		{"Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)", true},
		{"SemrushBot", true},
		{"Mozilla/5.0 (compatible; Konqueror/3.5; Linux) KHTML/3.5.5 (like Gecko) (Exabot-Thumbnails); Mozilla/5.0 (compatible; Exabot/3.0; +http://www.exabot.com/go/robot)", true},
		{"Google Chrome", false},
		{"Google", false},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299", false},
		{"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)", false},
		{"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:50.0) Gecko/20100101 Firefox/50.0", false},
		{"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)", false},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36", false},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134", false},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36", false},
		{"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.90 Safari/537.36", false},
		{"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.112 Safari/537.36", false},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:63.0) Gecko/20100101 Firefox/63.0", false},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36", false},
		{"Mozilla/5.0 (iPhone; CPU iPhone OS 11_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.0 Mobile/15E148 Safari/604.1", false},
		{"Safari", false},
		{"Firefox", false},
		{"Anything", false},
	}

	for _, tt := range testCases {
		header := make(http.Header)
		header.Set("User-Agent", tt.userAgent)

		req := web.WrapRequest(
			&http.Request{
				Method: "GET",
				Header: header,
				Host:   "demo.test.fider.io",
			},
		)
		Expect(req.IsCrawler()).Equals(tt.isCrawler)
	}
}
