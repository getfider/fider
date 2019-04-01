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
	Users   storage.User
	Tenants storage.Tenant
	Posts   storage.Post
}

// SetCurrentTenant to current context
func (s *Services) SetCurrentTenant(tenant *models.Tenant) {
	s.Context = context.WithValue(s.Context, TenantCtxKey, tenant)
	s.Users.SetCurrentTenant(tenant)
	s.Tenants.SetCurrentTenant(tenant)
	s.Posts.SetCurrentTenant(tenant)
}

// SetCurrentUser to current context
func (s *Services) SetCurrentUser(user *models.User) {
	s.Context = context.WithValue(s.Context, UserCtxKey, user)
	s.Users.SetCurrentUser(user)
	s.Tenants.SetCurrentUser(user)
	s.Posts.SetCurrentUser(user)
}
