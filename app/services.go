package app

import (
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage"
)

// Services holds reference to all Fider services
type Services struct {
	OAuth   oauth.Service
	Users   storage.User
	Tenants storage.Tenant
	Ideas   storage.Idea
	Emailer email.Sender
}
