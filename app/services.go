package app

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/billing"
	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage"
)

// Services holds reference to all Fider services
type Services struct {
	Logger        log.Logger
	OAuth         oauth.Service
	Users         storage.User
	Tags          storage.Tag
	Tenants       storage.Tenant
	Notifications storage.Notification
	Posts         storage.Post
	Emailer       email.Sender
	Blobs         blob.Storage
	Billing       *billing.Client
}

// SetCurrentTenant to current context
func (s *Services) SetCurrentTenant(tenant *models.Tenant) {
	s.Users.SetCurrentTenant(tenant)
	s.Tags.SetCurrentTenant(tenant)
	s.Tenants.SetCurrentTenant(tenant)
	s.Posts.SetCurrentTenant(tenant)
	s.Notifications.SetCurrentTenant(tenant)
	s.Blobs.SetCurrentTenant(tenant)
	s.Billing.SetCurrentTenant(tenant)
}

// SetCurrentUser to current context
func (s *Services) SetCurrentUser(user *models.User) {
	s.Users.SetCurrentUser(user)
	s.Tags.SetCurrentUser(user)
	s.Tenants.SetCurrentUser(user)
	s.Posts.SetCurrentUser(user)
	s.Notifications.SetCurrentUser(user)
	s.Billing.SetCurrentUser(user)
}
