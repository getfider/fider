package tasks

import (
	"fmt"
	"html/template"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/worker"
)

func describe(name string, job worker.Job) worker.Task {
	return worker.Task{Name: name, Job: job}
}

func link(baseURL, path string, args ...interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("<a href='%[1]s%[2]s'>%[1]s%[2]s</a>", baseURL, fmt.Sprintf(path, args...)))
}

//SendSignUpEmail is used to send the sign up email to requestor
func SendSignUpEmail(model *models.CreateTenant, baseURL string) worker.Task {
	return describe("Send sign up e-mail", func(c *worker.Context) error {
		to := email.NewRecipient(model.Email, web.Map{
			"link": link(baseURL, "/signup/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send("signup_email", "Fider", to)
	})
}

//SendSignInEmail is used to send the sign in email to requestor
func SendSignInEmail(model *models.SignInByEmail, baseURL string) worker.Task {
	return describe("Send sign in e-mail", func(c *worker.Context) error {
		to := email.NewRecipient(model.Email, web.Map{
			"tenantName": c.Tenant().Name,
			"link":       link(baseURL, "/signin/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send("signin_email", c.Tenant().Name, to)
	})
}

//SendChangeEmailConfirmation is used to send the change e-mail confirmation e-mail to requestor
func SendChangeEmailConfirmation(model *models.ChangeUserEmail, baseURL string) worker.Task {
	return describe("Send change e-mail confirmation", func(c *worker.Context) error {
		previous := c.User().Email
		if previous == "" {
			previous = "(empty)"
		}

		to := email.NewRecipient(model.Email, web.Map{
			"name":     c.User().Name,
			"oldEmail": previous,
			"newEmail": model.Email,
			"link":     link(baseURL, "/change-email/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send("change_emailaddress_email", c.Tenant().Name, to)
	})
}
