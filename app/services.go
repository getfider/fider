package app

import (
	"context"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage"
)

// Services holds reference to all Fider services
type Services struct {
	Context context.Context //EXPERIMENTAL-BUS: temporary
	OAuth   oauth.Service
	Tenants storage.Tenant
}

// SetCurrentTenant to current context
func (s *Services) SetCurrentTenant(tenant *models.Tenant) {
	s.Context = context.WithValue(s.Context, TenantCtxKey, tenant)
	s.Tenants.SetCurrentTenant(tenant)
}

// SetCurrentUser to current context
func (s *Services) SetCurrentUser(user *models.User) {
	s.Context = context.WithValue(s.Context, UserCtxKey, user)
	s.Tenants.SetCurrentUser(user)
}
