package tasks

import (
	"github.com/Spicy-Bush/fider-tarkov-community/app/actions"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/cmd"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/env"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/worker"
)

// SendSignUpEmail is used to send the sign up email to requestor
func SendSignUpEmail(action *actions.CreateTenant, baseURL string) worker.Task {
	return describe("Send sign up email", func(c *worker.Context) error {
		to := dto.NewRecipient(action.Name, action.Email, dto.Props{
			"link": link(baseURL, "/signup/verify?k=%s", action.VerificationKey),
		})

		bus.Publish(c, &cmd.SendMail{
			From:         dto.Recipient{Name: "Fider"},
			To:           []dto.Recipient{to},
			TemplateName: "signup_email",
			Props: dto.Props{
				"logo": web.LogoURL(c),
			},
		})

		return nil
	})
}

// SendWelcomeEmail is used to send a welcome email to new tenant admin
// This email is not sent on self hosted instaces
func SendWelcomeEmail(name, email, baseURL string) worker.Task {
	return describe("Send welcome email", func(c *worker.Context) error {
		if env.IsSingleHostMode() {
			return nil
		}

		to := dto.NewRecipient(name, email, dto.Props{
			"name": name,
			"url":  link(baseURL, "/"),
		})

		bus.Publish(c, &cmd.SendMail{
			From: dto.Recipient{
				Name:    "Fider",
				Address: "contact@fider.io",
			},
			To:           []dto.Recipient{to},
			TemplateName: "welcome_email",
			Props: dto.Props{
				"logo": web.LogoURL(c),
			},
		})

		return nil
	})
}
