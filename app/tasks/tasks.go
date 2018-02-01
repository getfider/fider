package tasks

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/worker"
)

func describe(name string, job worker.Job) worker.Task {
	return worker.Task{Name: name, Job: job}
}

//SendSignInEmail is used to send the sign in email to requestor
func SendSignInEmail(model *models.SignInByEmail, baseURL string) worker.Task {
	return describe("Send sign in e-mail", func(c *worker.Context) error {
		return c.Services().Emailer.Send(c.Tenant().Name, model.Email, "signin_email", web.Map{
			"tenantName":      c.Tenant().Name,
			"baseURL":         baseURL,
			"verificationKey": model.VerificationKey,
		})
	})
}

//SendChangeEmailConfirmation is used to send the change e-mail confirmation e-mail to requestor
func SendChangeEmailConfirmation(model *models.ChangeUserEmail, baseURL string) worker.Task {
	return describe("Send change e-mail confirmation", func(c *worker.Context) error {
		previous := c.User().Email
		if previous == "" {
			previous = "(empty)"
		}
		return c.Services().Emailer.Send(c.Tenant().Name, model.Email, "change_emailaddress_email", web.Map{
			"name":            c.User().Name,
			"oldEmail":        previous,
			"newEmail":        model.Email,
			"baseURL":         baseURL,
			"verificationKey": model.VerificationKey,
		})
	})
}
