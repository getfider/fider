package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/validate"
)

// Map defines a generic map of type `map[string]interface{}`.
type Map map[string]interface{}

// StringMap defines a map of type `map[string]string`.
type StringMap map[string]string

// HTMLMimeType is the mimetype for HTML responses
var (
	PlainContentType = "text/plain; charset=utf-8"
	HTMLContentType  = "text/html; charset=utf-8"
	JSONContentType  = "application/json; charset=utf-8"
)

// CookieAuthName is the name of the authentication cookie
const CookieAuthName = "auth"

var (
	preffixKey             = "__CTX_"
	tenantContextKey       = preffixKey + "TENANT"
	userContextKey         = preffixKey + "USER"
	authEndpointContextKey = preffixKey + "AUTH_ENDPOINT"
	transactionContextKey  = preffixKey + "TRANSACTION"
	servicesContextKey     = preffixKey + "SERVICES"
)

type Context struct {
	engine *Engine
	res    http.ResponseWriter
	req    *http.Request
	logger log.Logger
	params StringMap
	store  Map
}

//Tenant returns current tenant
func (ctx *Context) Tenant() *models.Tenant {
	tenant, ok := ctx.Get(tenantContextKey).(*models.Tenant)
	if ok {
		return tenant
	}
	return nil
}

//SetTenant update HTTP context with current tenant
func (ctx *Context) SetTenant(tenant *models.Tenant) {
	if tenant != nil {
		ctx.Logger().Debugf("Current tenant: %v (ID: %v)", tenant.Name, tenant.ID)
	} else {
		ctx.Logger().Debugf("Current tenant: nil")
	}
	ctx.Services().Tenants.SetCurrentTenant(tenant)
	ctx.Set(tenantContextKey, tenant)
}

//BindTo context values into given model
func (ctx *Context) BindTo(i actions.Actionable) *validate.Result {
	err := ctx.Bind(i.Initialize())
	if err != nil {
		return validate.Error(err)
	}
	if !i.IsAuthorized(ctx.User()) {
		return validate.Unauthorized()
	}
	return i.Validate(ctx.Services())
}

//Logger returns current logger
func (ctx *Context) Logger() log.Logger {
	return ctx.logger
}

//IsAuthenticated returns true if user is authenticated
func (ctx *Context) IsAuthenticated() bool {
	return ctx.Get(userContextKey) != nil
}

//IsAjax returns true if request is AJAX
func (ctx *Context) IsAjax() bool {
	return strings.Contains(ctx.Request().Header.Get("Accept"), "application/json")
}

//NotFound returns a 404 page
func (ctx *Context) NotFound() error {
	return ctx.Render(http.StatusNotFound, "404.html", Map{})
}

//Failure returns a 500 page
func (ctx *Context) Failure(err error) error {
	panic(err)
}

//HandleValidation handles given validation result property to return 400 or 500
func (ctx *Context) HandleValidation(result *validate.Result) error {
	if result.Error != nil {
		return ctx.Failure(result.Error)
	}

	if !result.Authorized {
		return ctx.Unauthorized()
	}

	return ctx.BadRequest(Map{
		"message":  result.Messages,
		"failures": result.Failures,
	})
}

//Unauthorized returns a 401 response
func (ctx *Context) Unauthorized() error {
	return ctx.JSON(http.StatusUnauthorized, Map{})
}

//Ok returns 200 OK with JSON result
func (ctx *Context) Ok(data interface{}) error {
	return ctx.JSON(http.StatusOK, data)
}

//BadRequest returns 400 BadRequest with JSON result
func (ctx *Context) BadRequest(dict Map) error {
	return ctx.JSON(http.StatusBadRequest, dict)
}

//Page returns a page with given variables
func (ctx *Context) Page(dict Map) error {
	return ctx.Render(http.StatusOK, "index.html", dict)
}

// Error logs and returns a 500 response
func (ctx *Context) Error(err error) error {
	ctx.Logger().Error(err)
	return ctx.Render(http.StatusInternalServerError, "500.html", Map{})
}

// Render renders a template with data and sends a text/html response with status
func (ctx *Context) Render(code int, template string, data Map) error {
	buf := new(bytes.Buffer)
	if err := ctx.engine.renderer.Render(buf, template, data, ctx); err != nil {
		return err
	}
	return ctx.Blob(code, HTMLContentType, buf.Bytes())
}

//SetParams sets path parameter names and values.
func (ctx *Context) SetParams(dict StringMap) {
	if dict == nil {
		return
	}

	ctx.params = dict
}

//User returns authenticated user
func (ctx *Context) User() *models.User {
	user, ok := ctx.Get(userContextKey).(*models.User)
	if ok {
		return user
	}
	return nil
}

//SetUser update HTTP context with current user
func (ctx *Context) SetUser(user *models.User) {
	if user != nil {
		ctx.Logger().Debugf("Logged as: %v [%v] (ID: %v)", user.Name, user.Email, user.ID)
	} else {
		ctx.Logger().Debugf("Logged as: nil")
	}
	ctx.Set(userContextKey, user)
}

//Services returns current app.Services from context
func (ctx *Context) Services() *app.Services {
	return ctx.Get(servicesContextKey).(*app.Services)
}

//AddAuthCookie generates and adds a cookie
func (ctx *Context) AddAuthCookie(user *models.User) (string, error) {
	token, err := jwt.Encode(models.FiderClaims{
		UserID:    user.ID,
		UserName:  user.Name,
		UserEmail: user.Email,
	})

	if err != nil {
		return token, err
	}

	ctx.AddCookie(CookieAuthName, token, time.Now().Add(365*24*time.Hour))
	return token, nil
}

//AddCookie adds a cookie
func (ctx *Context) AddCookie(name, value string, expires time.Time) {
	ctx.SetCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
		Expires:  expires,
	})
}

//RemoveCookie removes a cookie
func (ctx *Context) RemoveCookie(name string) {
	ctx.SetCookie(&http.Cookie{
		Name:    name,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
	})
}

//SetServices update current context with app.Services
func (ctx *Context) SetServices(services *app.Services) {
	ctx.Set(servicesContextKey, services)
}

//SetActiveTransaction adds transaction to context
func (ctx *Context) SetActiveTransaction(trx *dbx.Trx) {
	ctx.Set(transactionContextKey, trx)
}

//ActiveTransaction returns current active Database transaction
func (ctx *Context) ActiveTransaction() *dbx.Trx {
	return ctx.Get(transactionContextKey).(*dbx.Trx)
}

//BaseURL returns base URL as string
func (ctx *Context) BaseURL() string {
	protocol := "http"
	if ctx.Request().TLS != nil {
		protocol = "https"
	}
	return protocol + "://" + ctx.req.Host
}

//TenantBaseURL returns base URL for a given tenant
func (ctx *Context) TenantBaseURL(tenant *models.Tenant) string {
	if env.IsSingleHostMode() {
		return ctx.BaseURL()
	}

	protocol := "http"
	if ctx.req.TLS != nil {
		protocol = "https"
	}
	return protocol + "://" + tenant.Subdomain + env.MultiTenantDomain()
}

//AuthEndpoint auth endpoint
func (ctx *Context) AuthEndpoint() string {
	endpoint, ok := ctx.Get(authEndpointContextKey).(string)
	if !ok {
		if env.IsSingleHostMode() {
			endpoint = ctx.BaseURL()
		} else {
			endpoint = env.MustGet("AUTH_ENDPOINT")
		}
		ctx.Set(authEndpointContextKey, endpoint)
	}
	return endpoint
}

//QueryParam returns querystring parameter for given key
func (ctx *Context) QueryParam(key string) string {
	return ctx.req.URL.Query().Get(key)
}

//Param returns parameter as string
func (ctx *Context) Param(name string) string {
	if ctx.params == nil {
		return ""
	}
	return ctx.params[name]
}

//ParamAsInt returns parameter as int
func (ctx *Context) ParamAsInt(name string) (int, error) {
	val, err := strconv.Atoi(ctx.Param(name))
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

//RenderVars returns all registered RenderVar
func (ctx *Context) RenderVars() Map {
	vars := ctx.Get("__renderVars")
	if vars != nil {
		return vars.(Map)
	}
	return nil
}

//AddRenderVar register given key/value to RenderVar map
func (ctx *Context) AddRenderVar(name string, value interface{}) {
	var renderVars = ctx.Get("__renderVars")
	if renderVars == nil {
		renderVars = make(Map)
		ctx.Set("__renderVars", renderVars)
	}

	ctx.Logger().Debugf("storage.%v: %v", name, value)
	renderVars.(Map)[name] = value
}

//Request returns current request
func (ctx *Context) Request() *http.Request {
	return ctx.req
}

//Response returns current response writer
func (ctx *Context) Response() http.ResponseWriter {
	return ctx.res
}

// Get retrieves data from the context.
func (ctx *Context) Get(key string) interface{} {
	return ctx.store[key]
}

// Set saves data in the context.
func (ctx *Context) Set(key string, val interface{}) {
	if ctx.store == nil {
		ctx.store = make(Map)
	}
	ctx.store[key] = val
}

// String returns a text response with status code.
func (ctx *Context) String(code int, text string) error {
	return ctx.Blob(code, PlainContentType, []byte(text))
}

// HTML returns a HTML response with status code.
func (ctx *Context) HTML(code int, html string) error {
	return ctx.Blob(code, HTMLContentType, []byte(html))
}

// JSON returns a JSON response with status code.
func (ctx *Context) JSON(code int, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return ctx.Blob(code, JSONContentType, b)
}

// Blob sends a blob response with status code and content type.
func (ctx *Context) Blob(code int, contentType string, b []byte) (err error) {
	ctx.res.Header().Set("Content-Type", contentType)
	ctx.res.WriteHeader(code)
	_, err = ctx.res.Write(b)
	return
}

// Bind binds the request body into provided type `i`. The default binder
// does it based on Content-Type header.
func (ctx *Context) Bind(i interface{}) error {
	req := ctx.Request()
	if req.Method == "POST" {
		if err := json.NewDecoder(req.Body).Decode(i); err != nil {
			return err
		}
	}
	return nil
}

// Cookie returns the named cookie provided in the request.
func (ctx *Context) Cookie(name string) (*http.Cookie, error) {
	return ctx.req.Cookie(name)
}

// SetCookie adds a `Set-Cookie` header in HTTP response.
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.res, cookie)
}

// Cookies returns the HTTP cookies sent with the request.
func (ctx *Context) Cookies() []*http.Cookie {
	return ctx.req.Cookies()
}

// NoContent sends a response with no body and a status code.
func (ctx *Context) NoContent(code int) error {
	ctx.res.WriteHeader(code)
	return nil
}

// Redirect redirects the request to a provided URL with status code.
func (ctx *Context) Redirect(code int, url string) error {
	ctx.res.Header().Set("Location", url)
	ctx.res.WriteHeader(code)
	return nil
}
