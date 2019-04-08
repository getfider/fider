package app

import (
	"context"

	"github.com/getfider/fider/app/models"
)

// Services holds reference to all Fider services
type Services struct {
	Context context.Context //EXPERIMENTAL-BUS: temporary
}

// SetCurrentTenant to current context
func (s *Services) SetCurrentTenant(tenant *models.Tenant) {
	s.Context = context.WithValue(s.Context, TenantCtxKey, tenant)
}

// SetCurrentUser to current context
func (s *Services) SetCurrentUser(user *models.User) {
	s.Context = context.WithValue(s.Context, UserCtxKey, user)
}
