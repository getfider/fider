package app

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

//Context wraps echo.context to provide userful WeCHY information
type Context struct {
	echo.Context
}

//Tenant returns current tenant
func (ctx *Context) Tenant() *Tenant {
	tenant, ok := ctx.Get("Tenant").(*Tenant)
	if ok {
		return tenant
	}
	return nil
}

//IsAuthenticated returns true if user is authenticated
func (ctx *Context) IsAuthenticated() bool {
	return ctx.Get("Claims") != nil
}

//Failure returns a 500 response
func (ctx *Context) Failure(err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, err)
}

//Claims returns authenticated user claims
func (ctx *Context) Claims() *WechyClaims {
	return ctx.Get("Claims").(*WechyClaims)
}

//ParamAsInt64 returns parameter as int64
func (ctx *Context) ParamAsInt64(name string) (int64, error) {
	val, err := strconv.Atoi(ctx.Param(name))
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

//HandlerFunc represents an HTTP handler
type HandlerFunc func(Context) error

//MiddlewareFunc represents an HTTP middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc
