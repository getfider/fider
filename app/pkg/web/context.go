package web

import (
	"net/http"
	"strconv"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/labstack/echo"
)

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
	ctx.Set(tenantContextKey, tenant)
}

//IsAuthenticated returns true if user is authenticated
func (ctx *Context) IsAuthenticated() bool {
	return ctx.Get(userContextKey) != nil
}

//NotFound returns a 404 page
func (ctx *Context) NotFound() error {
	return echo.NewHTTPError(http.StatusNotFound)
}

//Failure returns a 500 page
func (ctx *Context) Failure(err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, err)
}

//Page returns a page with given variables
func (ctx *Context) Page(dict echo.Map) error {
	return ctx.Render(200, "index.html", dict)
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
func (ctx *Context) SetUser(claims *models.User) {
	ctx.Set(userContextKey, claims)
}

//Services returns current app.Services from context
func (ctx *Context) Services() *app.Services {
	return ctx.Get(servicesContextKey).(*app.Services)
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
