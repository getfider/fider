package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/getfider/fider/app/pkg/dbx"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/pkg/worker"
	"github.com/getfider/fider/app/services/blob"
)

// Map defines a generic map of type `map[string]interface{}`
type Map map[string]interface{}

// StringMap defines a map of type `map[string]string`
type StringMap map[string]string

// Props defines the data required to render rages
type Props struct {
	Title       string
	Description string
	ChunkName   string
	Data        Map
}

// HTMLMimeType is the mimetype for HTML responses
var (
	PlainContentType     = "text/plain"
	HTMLContentType      = "text/html"
	JSONContentType      = "application/json"
	XMLContentType       = "application/xml"
	UTF8PlainContentType = PlainContentType + "; charset=utf-8"
	UTF8HTMLContentType  = HTMLContentType + "; charset=utf-8"
	UTF8XMLContentType   = XMLContentType + "; charset=utf-8"
	UTF8JSONContentType  = JSONContentType + "; charset=utf-8"
)

// CookieSessionName is the name of the cookie that holds the session ID
const CookieSessionName = "user_session_id"

// CookieAuthName is the name of the cookie that holds the Authentication Token
const CookieAuthName = "auth"

// CookieSignUpAuthName is the name of the cookie that holds the temporary Authentication Token
const CookieSignUpAuthName = "__signup_auth"

//Context shared between http pipeline
type Context struct {
	context.Context
	Response           http.ResponseWriter
	Request            Request
	ResponseStatusCode int
	id                 string
	sessionID          string
	engine             *Engine
	params             StringMap
	tasks              []worker.Task
}

//NewContext creates a new web Context
func NewContext(engine *Engine, req *http.Request, res http.ResponseWriter, params StringMap) *Context {
	contextID := rand.String(32)

	wrappedRequest := WrapRequest(req)

	ctx := context.WithValue(req.Context(), app.RequestCtxKey, wrappedRequest)

	ctx = log.WithProperties(ctx, dto.Props{
		log.PropertyKeyContextID: contextID,
		log.PropertyKeyTag:       "WEB",
	})

	return &Context{
		Context:  ctx,
		id:       contextID,
		engine:   engine,
		Request:  wrappedRequest,
		Response: res,
		params:   params,
		tasks:    make([]worker.Task, 0),
	}
}

//Engine returns main HTTP engine
func (c *Context) Engine() *Engine {
	return c.engine
}

//SessionID returns the current session ID
func (c *Context) SessionID() string {
	return c.sessionID
}

//SetSessionID sets the session ID on current context
func (c *Context) SetSessionID(id string) {
	c.sessionID = id
	c.Context = log.WithProperty(c.Context, log.PropertyKeySessionID, id)
}

//ContextID returns the unique id for this context
func (c *Context) ContextID() string {
	return c.id
}

//Commit everything that is pending on current context
func (c *Context) Commit() error {
	trx, ok := c.Value(app.TransactionCtxKey).(*dbx.Trx)
	if ok && trx != nil {
		if err := trx.Commit(); err != nil {
			return err
		}
	}

	for _, task := range c.tasks {
		c.engine.worker.Enqueue(task)
	}

	return nil
}

//Rollback everything that is pending on current context
func (c *Context) Rollback() {
	trx, ok := c.Value(app.TransactionCtxKey).(*dbx.Trx)
	if ok && trx != nil {
		trx.MustRollback()
	}
}

//Enqueue given task to be processed in background
func (c *Context) Enqueue(task worker.Task) {
	task.OriginContext = c
	c.tasks = append(c.tasks, task)
}

//Tenant returns current tenant
func (c *Context) Tenant() *entity.Tenant {
	tenant, ok := c.Value(app.TenantCtxKey).(*entity.Tenant)
	if ok {
		return tenant
	}
	return nil
}

//SetTenant update HTTP context with current tenant
func (c *Context) SetTenant(tenant *entity.Tenant) {
	if tenant != nil {
		c.Context = log.WithProperty(c.Context, log.PropertyKeyTenantID, tenant.ID)
	}
	c.Set(app.TenantCtxKey, tenant)
}

//Bind context values into given model
func (c *Context) Bind(i interface{}) error {
	err := c.engine.binder.Bind(i, c)
	if err != nil {
		return errors.Wrap(err, "failed to bind request to model")
	}
	return nil
}

//BindTo context values into given model
func (c *Context) BindTo(i actions.Actionable) *validate.Result {
	err := c.engine.binder.Bind(i, c)
	if err != nil {
		if err == ErrContentTypeNotAllowed {
			return validate.Failed(err.Error())
		}
		return validate.Error(errors.Wrap(err, "failed to bind request to action"))
	}

	if v, ok := i.(actions.PreExecuteAction); ok {
		if err := v.OnPreExecute(c); err != nil {
			return validate.Error(err)
		}
	}

	if !i.IsAuthorized(c, c.User()) {
		return validate.Unauthorized()
	}

	return i.Validate(c, c.User())
}

//IsAuthenticated returns true if user is authenticated
func (c *Context) IsAuthenticated() bool {
	return c.Value(app.UserCtxKey) != nil
}

//IsAjax returns true if request is AJAX
func (c *Context) IsAjax() bool {
	accept := c.Request.GetHeader("Accept")
	contentType := c.Request.GetHeader("Content-Type")
	return strings.Contains(accept, JSONContentType) || strings.Contains(contentType, JSONContentType)
}

//Unauthorized returns a 403 response
func (c *Context) Unauthorized() error {
	return c.Render(http.StatusForbidden, "403.html", Props{
		Title:       "Not Authorized",
		Description: "You are not authorized to view this page.",
	})
}

//NotFound returns a 404 page
func (c *Context) NotFound() error {
	return c.Render(http.StatusNotFound, "404.html", Props{
		Title:       "Page not found",
		Description: "The link you clicked may be broken or the page may have been removed.",
	})
}

//Gone returns a 410 page
func (c *Context) Gone() error {
	return c.Render(http.StatusGone, "410.html", Props{
		Title:       "Expired",
		Description: "The link you clicked has expired.",
	})
}

//Failure returns a 500 page
func (c *Context) Failure(err error) error {
	err = errors.StackN(err, 1)
	cause := errors.Cause(err)

	if cause == context.Canceled {
		return nil
	}

	if cause == app.ErrNotFound || cause == blob.ErrNotFound {
		return c.NotFound()
	}

	if renderErr := c.Render(http.StatusInternalServerError, "500.html", Props{
		Title:       "Shoot! Well, this is unexpected…",
		Description: "An error has occurred and we're working to fix the problem!",
	}); renderErr != nil {
		return renderErr
	}
	return err
}

//HandleValidation handles given validation result property to return 400 or 500
func (c *Context) HandleValidation(result *validate.Result) error {
	if result.Err != nil {
		return c.Failure(result.Err)
	}

	if !result.Authorized {
		return c.Unauthorized()
	}

	return c.BadRequest(Map{
		"errors": result.Errors,
	})
}

//Attachment returns an attached file
func (c *Context) Attachment(fileName, contentType string, file []byte) error {
	c.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	return c.Blob(http.StatusOK, contentType, file)
}

//Ok returns 200 OK with JSON result
func (c *Context) Ok(data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

//BadRequest returns 400 BadRequest with JSON result
func (c *Context) BadRequest(dict Map) error {
	return c.JSON(http.StatusBadRequest, dict)
}

//Page returns a page with given variables
func (c *Context) Page(props Props) error {
	if len(env.Config.Rendergun.URL) > 0 && c.Request.IsCrawler() {
		html := new(bytes.Buffer)
		c.engine.renderer.Render(html, http.StatusOK, "index.html", props, c)
		return c.prerender(http.StatusOK, html)
	}

	return c.Render(http.StatusOK, "index.html", props)
}

// Render renders a template with data and sends a text/html response with status
func (c *Context) Render(code int, template string, props Props) error {
	if c.IsAjax() {
		return c.JSON(code, Map{})
	}

	buf := new(bytes.Buffer)
	c.engine.renderer.Render(buf, code, template, props, c)

	return c.Blob(code, UTF8HTMLContentType, buf.Bytes())
}

func (c *Context) prerender(code int, html io.Reader) error {
	req := &cmd.HTTPRequest{
		Method: "POST",
		URL:    fmt.Sprintf("%s/render?url=%s", env.Config.Rendergun.URL, c.Request.URL.String()),
		Body:   html,
		Headers: map[string]string{
			"Content-Type":              "text/html",
			"x-rendergun-wait-until":    "networkidle0",
			"x-rendergun-block-ads":     "true",
			"x-rendergun-abort-request": "assets\\/css\\/(common|vendor|main)\\.",
		},
	}
	err := bus.Dispatch(c, req)
	if err != nil {
		log.Error(c, errors.Wrap(err, "failed to execute rendergun"))
		return c.TryAgainLater(24 * time.Hour)
	}

	return c.Blob(code, UTF8HTMLContentType, req.ResponseBody)
}

//TryAgainLater returns a service unavailable response with Retry-After header
func (c *Context) TryAgainLater(d time.Duration) error {
	c.Response.Header().Set("Cache-Control", "no-cache, no-store")
	c.Response.Header().Set("Retry-After", fmt.Sprintf("%.0f", d.Seconds()))
	return c.NoContent(http.StatusServiceUnavailable)
}

//AddParam add a single param to route parameters list
func (c *Context) AddParam(name, value string) {
	c.params[name] = value
}

//User returns authenticated user
func (c *Context) User() *entity.User {
	user, ok := c.Value(app.UserCtxKey).(*entity.User)
	if ok {
		return user
	}
	return nil
}

//SetUser update HTTP context with current user
func (c *Context) SetUser(user *entity.User) {
	if user != nil {
		c.Context = log.WithProperty(c.Context, log.PropertyKeyUserID, user.ID)
	}
	c.Set(app.UserCtxKey, user)
}

//AddCookie adds a cookie
func (c *Context) AddCookie(name, value string, expires time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
		Expires:  expires,
		Secure:   c.Request.IsSecure,
	}
	http.SetCookie(c.Response, cookie)
	return cookie
}

//RemoveCookie removes a cookie
func (c *Context) RemoveCookie(name string) {
	http.SetCookie(c.Response, &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now().Add(-100 * time.Hour),
		Secure:   c.Request.IsSecure,
	})
}

//BaseURL returns base URL
func (c *Context) BaseURL() string {
	return c.Request.BaseURL()
}

//QueryParam returns querystring parameter for given key
func (c *Context) QueryParam(key string) string {
	return c.Request.URL.Query().Get(key)
}

//QueryParamAsInt returns querystring parameter for given key
func (c *Context) QueryParamAsInt(key string) (int, error) {
	value := c.QueryParam(key)
	if value == "" {
		return 0, nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse %s to integer", value)
	}
	return intValue, nil
}

//QueryParamAsArray returns querystring parameter for given key as an array
func (c *Context) QueryParamAsArray(key string) []string {
	param := c.QueryParam(key)
	if param != "" {
		return strings.Split(param, ",")
	}
	return []string{}
}

//Param returns parameter as string
func (c *Context) Param(name string) string {
	if c.params == nil {
		return ""
	}
	return strings.TrimPrefix(c.params[name], "/")
}

//ParamAsInt returns parameter as int
func (c *Context) ParamAsInt(name string) (int, error) {
	value := c.Param(name)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse %s to integer", value)
	}
	return intValue, nil
}

// Set saves data in the context.
func (c *Context) Set(key interface{}, val interface{}) {
	c.Context = context.WithValue(c.Context, key, val)
}

// String returns a text response with status code.
func (c *Context) String(code int, text string) error {
	return c.Blob(code, UTF8PlainContentType, []byte(text))
}

// XML returns a XML response with status code.
func (c *Context) XML(code int, text string) error {
	return c.Blob(code, UTF8XMLContentType, []byte(text))
}

// JSON returns a JSON response with status code.
func (c *Context) JSON(code int, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return errors.Wrap(err, "failed to marshal response to JSON")
	}
	return c.Blob(code, UTF8JSONContentType, b)
}

// Image sends an image blob response with status code and content type.
func (c *Context) Image(contentType string, b []byte) error {
	if !strings.HasPrefix(contentType, "image/") {
		return c.Failure(errors.New("'%s' is not an image", c.Request.URL.String()))
	}
	return c.Blob(http.StatusOK, contentType, b)
}

// Blob sends a blob response with status code and content type.
func (c *Context) Blob(code int, contentType string, b []byte) error {
	c.Response.Header().Set("Content-Type", contentType)
	c.ResponseStatusCode = code
	c.Response.WriteHeader(code)
	_, err := c.Response.Write(b)
	return err
}

// NoContent sends a response with no body and a status code.
func (c *Context) NoContent(code int) error {
	c.ResponseStatusCode = code
	c.Response.WriteHeader(code)
	return nil
}

// Redirect the request to a provided URL
func (c *Context) Redirect(url string) error {
	c.Response.Header().Set("Cache-Control", "no-cache, no-store")
	c.Response.Header().Set("Location", url)
	c.ResponseStatusCode = http.StatusTemporaryRedirect
	c.Response.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

// PermanentRedirect the request to a provided URL
func (c *Context) PermanentRedirect(url string) error {
	c.Response.Header().Set("Cache-Control", "no-cache, no-store")
	c.Response.Header().Set("Location", url)
	c.ResponseStatusCode = http.StatusMovedPermanently
	c.Response.WriteHeader(http.StatusMovedPermanently)
	return nil
}

// SetCanonicalURL sets the canonical link on the HTTP Response Headers
func (c *Context) SetCanonicalURL(rawurl string) {
	u, err := url.Parse(rawurl)
	if err == nil {
		if u.Host == "" {
			baseURL, ok := c.Value("Canonical-BaseURL").(string)
			if !ok {
				baseURL = c.BaseURL()
			}
			if len(rawurl) > 0 && rawurl[0] != '/' {
				rawurl = "/" + rawurl
			}
			rawurl = baseURL + rawurl
		} else {
			c.Set("Canonical-BaseURL", u.Scheme+"://"+u.Host)
		}

		c.Set("Canonical-URL", rawurl)
	}
}

//TenantBaseURL returns base URL for a given tenant
func TenantBaseURL(ctx context.Context, tenant *entity.Tenant) string {
	request := ctx.Value(app.RequestCtxKey).(Request)

	if env.IsSingleHostMode() {
		return request.BaseURL()
	}

	address := request.URL.Scheme + "://"
	if tenant.CNAME != "" {
		address += tenant.CNAME
	} else {
		address += tenant.Subdomain + env.MultiTenantDomain()
	}

	if request.URL.Port() != "" {
		address += ":" + request.URL.Port()
	}

	return address
}

// AssetsURL return the full URL to a tenant-specific static asset
func AssetsURL(ctx context.Context, path string, a ...interface{}) string {
	request := ctx.Value(app.RequestCtxKey).(Request)
	tenant, hasTenant := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	path = fmt.Sprintf(path, a...)
	if env.Config.CDN.Host != "" && hasTenant {
		if env.IsSingleHostMode() {
			return request.URL.Scheme + "://" + env.Config.CDN.Host + path
		}
		return request.URL.Scheme + "://" + tenant.Subdomain + "." + env.Config.CDN.Host + path
	}
	return request.BaseURL() + path
}

// LogoURL return the full URL to the tenant-specific logo URL
func LogoURL(ctx context.Context) string {
	tenant, hasTenant := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	if hasTenant && tenant.LogoBlobKey != "" {
		return AssetsURL(ctx, "/static/images/%s?size=200", tenant.LogoBlobKey)
	}
	return "https://getfider.com/images/logo-100x100.png"
}

// BaseURL return the base URL from given context
func BaseURL(ctx context.Context) string {
	request, ok := ctx.Value(app.RequestCtxKey).(Request)
	if ok {
		return request.BaseURL()
	}
	return ""
}

// OAuthBaseURL returns the OAuth base URL used for host-wide OAuth authentication
// For Single Tenant HostMode, BaseURL is the current BaseURL
// For Multi Tenant HostMode, BaseURL is //login.{HOST_DOMAIN}
func OAuthBaseURL(ctx context.Context) string {
	request := ctx.Value(app.RequestCtxKey).(Request)

	if env.IsSingleHostMode() {
		return BaseURL(ctx)
	}

	oauthBaseURL := request.URL.Scheme + "://login" + env.MultiTenantDomain()
	port := request.URL.Port()
	if port != "" {
		oauthBaseURL += ":" + port
	}
	return oauthBaseURL
}
