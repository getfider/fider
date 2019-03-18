package di

import (
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/email/mailgun"
	"github.com/getfider/fider/app/pkg/email/noop"
	"github.com/getfider/fider/app/pkg/email/smtp"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

// NewEmailer creates a new emailer instance based on current configuration
func NewEmailer(logger log.Logger) email.Sender {
	if env.IsTest() {
		return noop.NewSender()
	}

	if env.Config.Email.Mailgun.APIKey != "" {
		return mailgun.NewSender(
			logger,
			env.Config.Email.Mailgun.Domain,
			env.Config.Email.Mailgun.APIKey,
		)
	}

	return smtp.NewSender(
		logger,
		env.Config.Email.SMTP.Host,
		env.Config.Email.SMTP.Port,
		env.Config.Email.SMTP.Username,
		env.Config.Email.SMTP.Password,
	)
}
