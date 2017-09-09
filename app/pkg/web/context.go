package web

import (
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
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/labstack/echo"
)

// Map defines a generic map of type `map[string]interface{}`.
type Map map[string]interface{}

var (
	preffixKey             = "__CTX_"
	tenantContextKey       = preffixKey + "TENANT"
	userContextKey         = preffixKey + "USER"
	authEndpointContextKey = preffixKey + "AUTH_ENDPOINT"
	transactionContextKey  = preffixKey + "TRANSACTION"
	servicesContextKey     = preffixKey + "SERVICES"
)

//Context wraps echo.context to provide userful information
type Context struct {
	echo.Context
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
		ctx.Logger().Debug("Current tenant: nil")
	}
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
	return echo.NewHTTPError(http.StatusNotFound)
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

//SetParams sets path parameter names and values.
func (ctx *Context) SetParams(dict Map) {
	if dict == nil {
		return
	}

	names := make([]string, len(dict))
	values := make([]string, len(dict))

	for name, value := range dict {
		names = append(names, name)
		values = append(values, fmt.Sprint(value))
	}
	ctx.SetParamNames(names...)
	ctx.SetParamValues(values...)
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
		ctx.Logger().Debug("Logged as: nil")
	}
	ctx.Set(userContextKey, user)
}

//Services returns current app.Services from context
func (ctx *Context) Services() *app.Services {
	return ctx.Get(servicesContextKey).(*app.Services)
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
	return protocol + "://" + ctx.Request().Host
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
