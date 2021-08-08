package web

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/getfider/fider/app/pkg/errors"
)

//Request wraps the http request object
type Request struct {
	instance      *http.Request
	Method        string
	ContentLength int64
	Body          string
	IsSecure      bool
	StartTime     time.Time
	URL           *url.URL
}

// WrapRequest returns Fider wrapper of HTTP Request
func WrapRequest(request *http.Request) Request {
	protocol := "http"
	if request.TLS != nil || request.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}

	host := request.Host
	if request.Header.Get("X-Forwarded-Host") != "" {
		host = request.Header.Get("X-Forwarded-Host")
	}

	fullURL := protocol + "://" + host + request.RequestURI
	u, err := url.Parse(fullURL)
	if err != nil {
		panic(errors.Wrap(err, "Failed to parse url '%s'", fullURL))
	}

	var bodyBytes []byte
	if request.ContentLength > 0 {
		bodyBytes, err = ioutil.ReadAll(request.Body)
		if err != nil {
			panic(errors.Wrap(err, "failed to read body").Error())
		}
	}

	return Request{
		instance:      request,
		Method:        request.Method,
		ContentLength: request.ContentLength,
		Body:          string(bodyBytes),
		URL:           u,
		IsSecure:      protocol == "https",
		StartTime:     time.Now(),
	}
}

// GetHeader returns the value of HTTP header from given key
func (r *Request) GetHeader(key string) string {
	return r.instance.Header.Get(key)
}

// SetHeader updates the value of HTTP header of given key
func (r *Request) SetHeader(key, value string) {
	r.instance.Header.Set(key, value)
}

// Cookie returns the named cookie provided in the request.
func (r *Request) Cookie(name string) (*http.Cookie, error) {
	cookie, err := r.instance.Cookie(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get '%s' cookie", name)
	}
	return cookie, nil
}

// AddCookie adds a cookie
func (r *Request) AddCookie(cookie *http.Cookie) {
	r.instance.AddCookie(cookie)
}

// IsAPI returns true if its a request for an API resource
func (r *Request) IsAPI() bool {
	return strings.HasPrefix(r.URL.Path, "/api/")
}

var crawlerRegex = regexp.MustCompile("(?i)(baidu)|(msnbot)|(bingbot)|(bingpreview)|(duckduckbot)|(googlebot)|(adsbot-google)|(mediapartners-google)|(slurp)|(yandexbot)|(yandexmetrika)|(ahrefsbot)|(twitterbot)|(slackbot)|(discordbot)|(semrushBot)|(exabot)")

// IsCrawler returns true if the request is coming from a crawler
func (r *Request) IsCrawler() bool {
	return crawlerRegex.MatchString(r.GetHeader("User-Agent"))
}

//BaseURL returns base URL
func (r *Request) BaseURL() string {
	address := r.URL.Scheme + "://" + r.URL.Hostname()

	if r.URL.Port() != "" {
		address += ":" + r.URL.Port()
	}

	return address
}
