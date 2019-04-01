package actions

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

type IsAuthorizedHandler func(tenant *models.Tenant, user *models.User, services *app.Services) bool

func isAuthorized(ctx context.Context, handler IsAuthorizedHandler) bool {
	tenant, _ := ctx.Value(app.TenantCtxKey).(*models.Tenant)
	user, _ := ctx.Value(app.UserCtxKey).(*models.User)
	services, _ := ctx.Value(app.ServicesCtxKey).(*app.Services)
	return handler(tenant, user, services)
}

type ValidateHandler func(tenant *models.Tenant, user *models.User, services *app.Services) *validate.Result

func using(ctx context.Context, handler ValidateHandler) *validate.Result {
	tenant, _ := ctx.Value(app.TenantCtxKey).(*models.Tenant)
	user, _ := ctx.Value(app.UserCtxKey).(*models.User)
	services, _ := ctx.Value(app.ServicesCtxKey).(*app.Services)
	return handler(tenant, user, services)
}
