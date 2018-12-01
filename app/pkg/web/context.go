package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/pkg/worker"
)

// Map defines a generic map of type `map[string]interface{}`
type Map map[string]interface{}

// StringMap defines a map of type `map[string]string`
type StringMap map[string]string

// Props defines the data required to render rages
type Props struct {
	Title       string
	Description string
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

var (
	prefixKey             = "__CTX_"
	tenantContextKey      = prefixKey + "TENANT"
	userContextKey        = prefixKey + "USER"
	claimsContextKey      = prefixKey + "CLAIMS"
	transactionContextKey = prefixKey + "TRANSACTION"
	servicesContextKey    = prefixKey + "SERVICES"
	tasksContextKey       = prefixKey + "TASKS"
)

//Context shared between http pipeline
type Context struct {
	id        string
	sessionID string
	Response  http.ResponseWriter
	Request   Request
	engine    *Engine
	logger    log.Logger
	params    StringMap
	store     Map
	worker    worker.Worker
}

//Engine returns main HTTP engine
func (ctx *Context) Engine() *Engine {
	return ctx.engine
}

//SessionID returns the current session ID
func (ctx *Context) SessionID() string {
	return ctx.sessionID
}

//SetSessionID sets the session ID on current context
func (ctx *Context) SetSessionID(id string) {
	ctx.sessionID = id
	ctx.logger.SetProperty(log.PropertyKeySessionID, id)
}

//ContextID returns the unique id for this context
func (ctx *Context) ContextID() string {
	return ctx.id
}

//Commit everything that is pending on current context
func (ctx *Context) Commit() error {
	if trx := ctx.ActiveTransaction(); trx != nil {
		if err := trx.Commit(); err != nil {
			return err
		}
	}

	tasks, ok := ctx.Get(tasksContextKey).([]worker.Task)
	if ok {
		for _, task := range tasks {
			ctx.worker.Enqueue(task)
		}
	}

	return nil
}

//Rollback everything that is pending on current context
func (ctx *Context) Rollback() error {
	if trx := ctx.ActiveTransaction(); trx != nil {
		return trx.Rollback()
	}

	return nil
}

//Enqueue given task to be processed in background
func (ctx *Context) Enqueue(task worker.Task) {
	wrap := func(c *Context) worker.Job {
		return func(wc *worker.Context) error {
			wc.SetUser(c.User())
			wc.SetTenant(c.Tenant())
			wc.SetBaseURL(c.BaseURL())
			wc.SetLogoURL(c.LogoURL())
			return task.Job(wc)
		}
	}

	tasks, ok := ctx.Get(tasksContextKey).([]worker.Task)
	if !ok {
		tasks = make([]worker.Task, 0)
	}

	ctx.Set(tasksContextKey, append(tasks, worker.Task{
		OriginSessionID: ctx.SessionID(),
		Name:            task.Name,
		Job:             wrap(ctx),
	}))
}

//Tenant returns current tenant
func (ctx Context) Tenant() *models.Tenant {
	tenant, ok := ctx.Get(tenantContextKey).(*models.Tenant)
	if ok {
		return tenant
	}
	return nil
}

//SetTenant update HTTP context with current tenant
func (ctx *Context) SetTenant(tenant *models.Tenant) {
	if tenant != nil {
		ctx.logger.SetProperty(log.PropertyKeyTenantID, tenant.ID)
	}
	if ctx.Services() != nil {
		ctx.Services().SetCurrentTenant(tenant)
	}
	ctx.Set(tenantContextKey, tenant)
}

//Bind context values into given model
func (ctx *Context) Bind(i interface{}) error {
	err := ctx.engine.binder.Bind(i, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to bind request to model")
	}
	return nil
}

//BindTo context values into given model
func (ctx *Context) BindTo(i actions.Actionable) *validate.Result {
	err := ctx.engine.binder.Bind(i.Initialize(), ctx)
	if err != nil {
		if err == ErrContentTypeNotAllowed {
			return validate.Failed(err.Error())
		}
		return validate.Error(errors.Wrap(err, "failed to bind request to action"))
	}
	if !i.IsAuthorized(ctx.User(), ctx.Services()) {
		return validate.Unauthorized()
	}
	return i.Validate(ctx.User(), ctx.Services())
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
	accept := ctx.Request.GetHeader("Accept")
	contentType := ctx.Request.GetHeader("Content-Type")
	return strings.Contains(accept, JSONContentType) || strings.Contains(contentType, JSONContentType)
}

//Unauthorized returns a 403 response
func (ctx *Context) Unauthorized() error {
	return ctx.Render(http.StatusForbidden, "403.html", Props{
		Title:       "Not Authorized",
		Description: "You are not authorized to view this page.",
	})
}

//NotFound returns a 404 page
func (ctx *Context) NotFound() error {
	return ctx.Render(http.StatusNotFound, "404.html", Props{
		Title:       "Page not found",
		Description: "The link you clicked may be broken or the page may have been removed.",
	})
}

//Gone returns a 410 page
func (ctx *Context) Gone() error {
	return ctx.Render(http.StatusGone, "410.html", Props{
		Title:       "Expired",
		Description: "The link you clicked has expired.",
	})
}

//Failure returns a 500 page
func (ctx *Context) Failure(err error) error {
	err = errors.StackN(err, 1)
	if errors.Cause(err) == app.ErrNotFound {
		return ctx.NotFound()
	}

	ctx.Logger().Errorf(err.Error(), log.Props{
		"Body":       ctx.Request.Body,
		"HttpMethod": ctx.Request.Method,
		"URL":        ctx.Request.URL.String(),
	})

	ctx.Render(http.StatusInternalServerError, "500.html", Props{
		Title:       "Shoot! Well, this is unexpected…",
		Description: "An error has occurred and we're working to fix the problem!",
	})
	return err
}

//HandleValidation handles given validation result property to return 400 or 500
func (ctx *Context) HandleValidation(result *validate.Result) error {
	if result.Err != nil {
		return ctx.Failure(result.Err)
	}

	if !result.Authorized {
		return ctx.Unauthorized()
	}

	return ctx.BadRequest(Map{
		"errors": result.Errors,
	})
}

//Attachment returns an attached file
func (ctx *Context) Attachment(fileName, contentType string, file []byte) error {
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	return ctx.Blob(http.StatusOK, contentType, file)
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
func (ctx *Context) Page(props Props) error {
	return ctx.Render(http.StatusOK, "index.html", props)
}

// Render renders a template with data and sends a text/html response with status
func (ctx *Context) Render(code int, template string, props Props) error {
	if ctx.IsAjax() {
		return ctx.JSON(code, Map{})
	}

	buf := new(bytes.Buffer)
	ctx.engine.renderer.Render(buf, template, props, ctx)
	return ctx.Blob(code, UTF8HTMLContentType, buf.Bytes())
}

//AddParam add a single param to route parameters list
func (ctx *Context) AddParam(name, value string) {
	ctx.params[name] = value
}

//Claims returns current user claims
func (ctx *Context) Claims() *jwt.FiderClaims {
	claims, ok := ctx.Get(claimsContextKey).(*jwt.FiderClaims)
	if ok {
		return claims
	}
	return nil
}

//SetClaims update HTTP context with current user claims
func (ctx *Context) SetClaims(claims *jwt.FiderClaims) {
	ctx.Set(claimsContextKey, claims)
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
		ctx.logger.SetProperty(log.PropertyKeyUserID, user.ID)
	}
	if ctx.Services() != nil {
		ctx.Services().SetCurrentUser(user)
	}
	ctx.Set(userContextKey, user)
}

//Services returns current app.Services from context
func (ctx *Context) Services() *app.Services {
	svc, ok := ctx.Get(servicesContextKey).(*app.Services)
	if ok {
		return svc
	}
	return nil
}

//AddCookie adds a cookie
func (ctx *Context) AddCookie(name, value string, expires time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
		Expires:  expires,
		Secure:   ctx.Request.IsSecure,
	}
	http.SetCookie(ctx.Response, cookie)
	return cookie
}

//RemoveCookie removes a cookie
func (ctx *Context) RemoveCookie(name string) {
	http.SetCookie(ctx.Response, &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now().Add(-100 * time.Hour),
		Secure:   ctx.Request.IsSecure,
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

//BaseURL returns base URL
func (ctx Context) BaseURL() string {
	address := ctx.Request.URL.Scheme + "://" + ctx.Request.URL.Hostname()

	if ctx.Request.URL.Port() != "" {
		address += ":" + ctx.Request.URL.Port()
	}

	return address
}

//TenantBaseURL returns base URL for a given tenant
func (ctx *Context) TenantBaseURL(tenant *models.Tenant) string {
	if env.IsSingleHostMode() {
		return ctx.BaseURL()
	}

	address := ctx.Request.URL.Scheme + "://"
	if tenant.CNAME != "" {
		address += tenant.CNAME
	} else {
		address += tenant.Subdomain + env.MultiTenantDomain()
	}

	if ctx.Request.URL.Port() != "" {
		address += ":" + ctx.Request.URL.Port()
	}

	return address
}

//QueryParam returns querystring parameter for given key
func (ctx *Context) QueryParam(key string) string {
	return ctx.Request.URL.Query().Get(key)
}

//QueryParamAsArray returns querystring parameter for given key as an array
func (ctx *Context) QueryParamAsArray(key string) []string {
	param := ctx.QueryParam(key)
	if param != "" {
		return strings.Split(param, ",")
	}
	return []string{}
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
	value := ctx.Param(name)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse %s to integer", value)
	}
	return intValue, nil
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
	return ctx.Blob(code, UTF8PlainContentType, []byte(text))
}

// XML returns a XML response with status code.
func (ctx *Context) XML(code int, text string) error {
	return ctx.Blob(code, UTF8XMLContentType, []byte(text))
}

// JSON returns a JSON response with status code.
func (ctx *Context) JSON(code int, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return errors.Wrap(err, "failed to marshal response to JSON")
	}
	return ctx.Blob(code, UTF8JSONContentType, b)
}

// Blob sends a blob response with status code and content type.
func (ctx *Context) Blob(code int, contentType string, b []byte) error {
	ctx.Response.Header().Set("Content-Type", contentType)
	ctx.Response.WriteHeader(code)
	_, err := ctx.Response.Write(b)
	return err
}

// NoContent sends a response with no body and a status code.
func (ctx *Context) NoContent(code int) error {
	ctx.Response.WriteHeader(code)
	return nil
}

// Redirect the request to a provided URL
func (ctx *Context) Redirect(url string) error {
	ctx.Response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response.Header().Set("Location", url)
	ctx.Response.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

// PermanentRedirect the request to a provided URL
func (ctx *Context) PermanentRedirect(url string) error {
	ctx.Response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response.Header().Set("Location", url)
	ctx.Response.WriteHeader(http.StatusMovedPermanently)
	return nil
}

// GlobalAssetsURL return the full URL to a globally shared static asset
func (ctx *Context) GlobalAssetsURL(path string, a ...interface{}) string {
	path = fmt.Sprintf(path, a...)
	if env.IsDefined("CDN_HOST") {
		if env.IsSingleHostMode() {
			return ctx.Request.URL.Scheme + "://" + env.MustGet("CDN_HOST") + path
		}
		return ctx.Request.URL.Scheme + "://cdn." + env.MustGet("CDN_HOST") + path
	}
	return ctx.BaseURL() + path
}

// TenantAssetsURL return the full URL to a tenant-specific static asset
func (ctx *Context) TenantAssetsURL(path string, a ...interface{}) string {
	path = fmt.Sprintf(path, a...)
	if env.IsDefined("CDN_HOST") && ctx.Tenant() != nil {
		if env.IsSingleHostMode() {
			return ctx.Request.URL.Scheme + "://" + env.MustGet("CDN_HOST") + path
		}
		return ctx.Request.URL.Scheme + "://" + ctx.Tenant().Subdomain + "." + env.MustGet("CDN_HOST") + path
	}
	return ctx.BaseURL() + path
}

// LogoURL return the full URL to the tenant-specific logo URL
func (ctx Context) LogoURL() string {
	if ctx.Tenant() != nil && ctx.Tenant().LogoID > 0 {
		return ctx.TenantAssetsURL("/images/200/%d", ctx.Tenant().LogoID)
	}
	return "https://getfider.com/images/logo-100x100.png"
}

// FaviconURL return the full URL to the tenant-specific favicon URL
func (ctx Context) FaviconURL() string {
	if ctx.Tenant() != nil && ctx.Tenant().LogoID > 0 {
		return ctx.TenantAssetsURL("/images/50/%d", ctx.Tenant().LogoID)
	}
	return ctx.GlobalAssetsURL("/favicon.ico")
}
