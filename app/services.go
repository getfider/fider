package app

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/email/mailgun"
	"github.com/getfider/fider/app/pkg/email/smtp"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage"
)

// Services holds reference to all Fider services
type Services struct {
	OAuth   oauth.Service
	Users   storage.User
	Tags    storage.Tag
	Tenants storage.Tenant
	Ideas   storage.Idea
	Emailer email.Sender
}

// SetCurrentTenant to current context
func (s *Services) SetCurrentTenant(tenant *models.Tenant) {
	s.Users.SetCurrentTenant(tenant)
	s.Tags.SetCurrentTenant(tenant)
	s.Tenants.SetCurrentTenant(tenant)
	s.Ideas.SetCurrentTenant(tenant)
}

// SetCurrentUser to current context
func (s *Services) SetCurrentUser(user *models.User) {
	s.Users.SetCurrentUser(user)
	s.Tags.SetCurrentUser(user)
	s.Tenants.SetCurrentUser(user)
	s.Ideas.SetCurrentUser(user)
}

//NewEmailer creates a new emailer based on system configuration
func NewEmailer(logger log.Logger) email.Sender {
	if env.IsTest() {
		return email.NewNoopSender()
	}
	if env.IsDefined("MAILGUN_API") {
		return mailgun.NewSender(logger, env.MustGet("MAILGUN_DOMAIN"), env.MustGet("MAILGUN_API"))
	}
	return smtp.NewSender(
		logger,
		env.MustGet("SMTP_HOST"),
		env.MustGet("SMTP_PORT"),
		env.MustGet("SMTP_USERNAME"),
		env.MustGet("SMTP_PASSWORD"),
	)
}
