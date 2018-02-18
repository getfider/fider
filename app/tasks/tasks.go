package tasks

import (
	"fmt"
	"html/template"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/worker"
)

func describe(name string, job worker.Job) worker.Task {
	return worker.Task{Name: name, Job: job}
}

func link(baseURL, path string, args ...interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("<a href='%[1]s%[2]s'>%[1]s%[2]s</a>", baseURL, fmt.Sprintf(path, args...)))
}

func linkWithText(text, baseURL, path string, args ...interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("<a href='%s%s'>%s</a>", baseURL, fmt.Sprintf(path, args...), text))
}

//SendSignUpEmail is used to send the sign up email to requestor
func SendSignUpEmail(model *models.CreateTenant, baseURL string) worker.Task {
	return describe("Send sign up e-mail", func(c *worker.Context) error {
		to := email.NewRecipient(model.Name, model.Email, email.Params{
			"link": link(baseURL, "/signup/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send("signup_email", email.Params{}, "Fider", to)
	})
}

//SendSignInEmail is used to send the sign in email to requestor
func SendSignInEmail(model *models.SignInByEmail) worker.Task {
	return describe("Send sign in e-mail", func(c *worker.Context) error {
		to := email.NewRecipient("", model.Email, email.Params{
			"tenantName": c.Tenant().Name,
			"link":       link(c.BaseURL(), "/signin/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send("signin_email", email.Params{}, c.Tenant().Name, to)
	})
}

//SendChangeEmailConfirmation is used to send the change e-mail confirmation e-mail to requestor
func SendChangeEmailConfirmation(model *models.ChangeUserEmail) worker.Task {
	return describe("Send change e-mail confirmation", func(c *worker.Context) error {
		previous := c.User().Email
		if previous == "" {
			previous = "(empty)"
		}

		to := email.NewRecipient(model.Requestor.Name, model.Email, email.Params{
			"name":     c.User().Name,
			"oldEmail": previous,
			"newEmail": model.Email,
			"link":     link(c.BaseURL(), "/change-email/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send("change_emailaddress_email", email.Params{}, c.Tenant().Name, to)
	})
}

//NotifyAboutNewIdea sends a notification (web and e-mail) to subscribers
func NotifyAboutNewIdea(idea *models.Idea) worker.Task {
	return describe("Notify about new idea", func(c *worker.Context) error {
		if env.GetEnvOrDefault("NOTIFICATIONS_ENABLED", "true") == "false" {
			c.Logger().Warnf("Notifications is currently disabled.")
			return nil
		}

		users, err := c.Services().Ideas.GetActiveSubscribers(idea.Number, models.NotificationChannelEmail, models.NotificationEventNewIdea)
		if err != nil {
			return err
		}

		to := make([]email.Recipient, 0)
		for _, user := range users {
			if user.ID != c.User().ID {
				to = append(to, email.NewRecipient(user.Name, user.Email, email.Params{}))
			}
		}

		params := email.Params{
			"title":   fmt.Sprintf("[%s] %s", c.Tenant().Name, idea.Title),
			"content": markdown.Parse(idea.Description),
			"view":    linkWithText("View it on your browser", c.BaseURL(), "/ideas/%d/%s", idea.Number, idea.Slug),
			"change":  linkWithText("change your notification settings", c.BaseURL(), "/settings"),
		}

		return c.Services().Emailer.BatchSend("new_idea", params, c.User().Name, to)
	})
}

//NotifyAboutNewComment sends a notification (web and e-mail) to subscribers
func NotifyAboutNewComment(idea *models.Idea, comment *models.NewComment) worker.Task {
	return describe("Notify about new comment", func(c *worker.Context) error {
		if env.GetEnvOrDefault("NOTIFICATIONS_ENABLED", "true") == "false" {
			c.Logger().Warnf("Notifications is currently disabled.")
			return nil
		}

		users, err := c.Services().Ideas.GetActiveSubscribers(comment.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
		if err != nil {
			return err
		}

		to := make([]email.Recipient, 0)
		for _, user := range users {
			if user.ID != c.User().ID {
				to = append(to, email.NewRecipient(user.Name, user.Email, email.Params{}))
			}
		}

		params := email.Params{
			"title":       fmt.Sprintf("[%s] %s", c.Tenant().Name, idea.Title),
			"content":     markdown.Parse(comment.Content),
			"view":        linkWithText("View it on your browser", c.BaseURL(), "/ideas/%d/%s", idea.Number, idea.Slug),
			"unsubscribe": linkWithText("unsubscribe from it", c.BaseURL(), "/ideas/%d/%s", idea.Number, idea.Slug),
			"change":      linkWithText("change your notification settings", c.BaseURL(), "/settings"),
		}

		return c.Services().Emailer.BatchSend("new_comment", params, c.User().Name, to)
	})
}

//NotifyAboutStatusChange sends a notification (web and e-mail) to subscribers
func NotifyAboutStatusChange(idea *models.Idea, response *models.SetResponse) worker.Task {
	return describe("Notify about idea status change", func(c *worker.Context) error {
		if env.GetEnvOrDefault("NOTIFICATIONS_ENABLED", "true") == "false" {
			c.Logger().Warnf("Notifications is currently disabled.")
			return nil
		}

		//Don't notify if status is the same
		if idea.Status == response.Status {
			return nil
		}

		users, err := c.Services().Ideas.GetActiveSubscribers(response.Number, models.NotificationChannelEmail, models.NotificationEventChangeStatus)
		if err != nil {
			return err
		}

		var duplicate template.HTML
		if response.Status == models.IdeaDuplicate {
			originalIdea, err := c.Services().Ideas.GetByNumber(response.OriginalNumber)
			if err != nil {
				return err
			}
			duplicate = linkWithText(originalIdea.Title, c.BaseURL(), "/ideas/%d/%s", originalIdea.Number, originalIdea.Slug)
		}

		to := make([]email.Recipient, 0)
		for _, user := range users {
			if user.ID != c.User().ID {
				to = append(to, email.NewRecipient(user.Name, user.Email, email.Params{}))
			}
		}

		params := email.Params{
			"title":       fmt.Sprintf("[%s] %s", c.Tenant().Name, idea.Title),
			"content":     markdown.Parse(response.Text),
			"status":      models.GetIdeaStatusName(response.Status),
			"duplicate":   duplicate,
			"view":        linkWithText("View it on your browser", c.BaseURL(), "/ideas/%d/%s", idea.Number, idea.Slug),
			"unsubscribe": linkWithText("unsubscribe from it", c.BaseURL(), "/ideas/%d/%s", idea.Number, idea.Slug),
			"change":      linkWithText("change your notification settings", c.BaseURL(), "/settings"),
		}

		return c.Services().Emailer.BatchSend("change_status", params, c.User().Name, to)
	})
}
