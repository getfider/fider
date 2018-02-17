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

//NotifyAboutNewIdea sends a notification (web and e-mail) to subscribers
func NotifyAboutNewIdea(idea *models.Idea) worker.Task {
	return describe("Notify about new idea", func(c *worker.Context) error {
		users, err := c.Services().Ideas.GetActiveSubscribers(idea.Number, models.NotificationChannelEmail, models.NotificationEventNewIdea)
		if err != nil {
			return err
		}

		for _, user := range users {
			if user.ID != c.User().ID && email.CanSendTo(user.Email) {
				c.Logger().Infof("Notify %s (%s)  about %s", user.Name, user.Email, idea.Title)
			}
		}

		return nil
	})
}

//NotifyAboutNewComment sends a notification (web and e-mail) to subscribers
func NotifyAboutNewComment(comment *models.NewComment) worker.Task {
	return describe("Notify about new comment", func(c *worker.Context) error {
		users, err := c.Services().Ideas.GetActiveSubscribers(comment.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
		if err != nil {
			return err
		}

		to := make([]email.Recipient, len(users))
		for i, user := range users {
			if user.ID != c.User().ID {
				to[i] = email.Recipient{
					Address: user.Email,
					Params: web.Map{
						"name": user.Name,
					},
				}
			}
		}

		return c.Services().Emailer.BatchSend("echo_test", "Testing...", to)
	})
}

//NotifyAboutStatusChange sends a notification (web and e-mail) to subscribers
func NotifyAboutStatusChange(response *models.SetResponse) worker.Task {
	return describe("Notify about new comment", func(c *worker.Context) error {
		users, err := c.Services().Ideas.GetActiveSubscribers(response.Number, models.NotificationChannelEmail, models.NotificationEventChangeStatus)
		if err != nil {
			return err
		}

		for _, user := range users {
			if user.ID != c.User().ID && email.CanSendTo(user.Email) {
				c.Logger().Infof("Notify %s (%s) about new status %d - %s", user.Name, user.Email, response.Status, response.Text)
			}
		}

		return nil
	})
}
