package web

import (
	"net/http"
	"strconv"

	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/labstack/echo"
)

var (
	preffixKey       = "__CTX_"
	tenantContextKey = preffixKey + "TENANT"
	userContextKey   = preffixKey + "USER"
)

//Context wraps echo.context to provide userful WeCHY information
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

//ParamAsInt returns parameter as int
func (ctx *Context) ParamAsInt(name string) (int, error) {
	val, err := strconv.Atoi(ctx.Param(name))
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

//HandlerFunc represents an HTTP handler
type HandlerFunc func(Context) error

//MiddlewareFunc represents an HTTP middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc
