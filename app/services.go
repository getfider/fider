package app

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage"
)

// Services holds reference to all Fider services
type Services struct {
	OAuth         oauth.Service
	Users         storage.User
	Tags          storage.Tag
	Tenants       storage.Tenant
	Notifications storage.Notification
	Ideas         storage.Idea
	Emailer       email.Sender
}

// SetCurrentTenant to current context
func (s *Services) SetCurrentTenant(tenant *models.Tenant) {
	s.Users.SetCurrentTenant(tenant)
	s.Tags.SetCurrentTenant(tenant)
	s.Tenants.SetCurrentTenant(tenant)
	s.Ideas.SetCurrentTenant(tenant)
	s.Notifications.SetCurrentTenant(tenant)
}

// SetCurrentUser to current context
func (s *Services) SetCurrentUser(user *models.User) {
	s.Users.SetCurrentUser(user)
	s.Tags.SetCurrentUser(user)
	s.Tenants.SetCurrentUser(user)
	s.Ideas.SetCurrentUser(user)
	s.Notifications.SetCurrentUser(user)
}
