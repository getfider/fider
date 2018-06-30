package tasks

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
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
	return describe("Send sign up email", func(c *worker.Context) error {
		to := email.NewRecipient(model.Name, model.Email, email.Params{
			"link": link(baseURL, "/signup/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send(c.Tenant(), "signup_email", email.Params{}, "Fider", to)
	})
}

//SendSignInEmail is used to send the sign in email to requestor
func SendSignInEmail(model *models.SignInByEmail) worker.Task {
	return describe("Send sign in email", func(c *worker.Context) error {
		to := email.NewRecipient("", model.Email, email.Params{
			"tenantName": c.Tenant().Name,
			"link":       link(c.BaseURL(), "/signin/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send(c.Tenant(), "signin_email", email.Params{}, c.Tenant().Name, to)
	})
}

//SendChangeEmailConfirmation is used to send the change email confirmation email to requestor
func SendChangeEmailConfirmation(model *models.ChangeUserEmail) worker.Task {
	return describe("Send change email confirmation", func(c *worker.Context) error {
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
		return c.Services().Emailer.Send(c.Tenant(), "change_emailaddress_email", email.Params{}, c.Tenant().Name, to)
	})
}

//NotifyAboutNewIdea sends a notification (web and email) to subscribers
func NotifyAboutNewIdea(idea *models.Idea) worker.Task {
	return describe("Notify about new idea", func(c *worker.Context) error {
		// Web notification
		users, err := c.Services().Ideas.GetActiveSubscribers(idea.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("New idea: **%s**", idea.Title)
		link := fmt.Sprintf("/ideas/%d/%s", idea.Number, idea.Slug)
		for _, user := range users {
			if user.ID != c.User().ID {
				if _, err = c.Services().Notifications.Insert(user, title, link, idea.ID); err != nil {
					return c.Failure(err)
				}
			}
		}

		// Email notification
		users, err = c.Services().Ideas.GetActiveSubscribers(idea.Number, models.NotificationChannelEmail, models.NotificationEventNewIdea)
		if err != nil {
			return c.Failure(err)
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

		return c.Services().Emailer.BatchSend(c.Tenant(), "new_idea", params, c.User().Name, to)
	})
}

//NotifyAboutNewComment sends a notification (web and email) to subscribers
func NotifyAboutNewComment(idea *models.Idea, comment *models.NewComment) worker.Task {
	return describe("Notify about new comment", func(c *worker.Context) error {
		// Web notification
		users, err := c.Services().Ideas.GetActiveSubscribers(idea.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("**%s** left a comment on **%s**", c.User().Name, idea.Title)
		link := fmt.Sprintf("/ideas/%d/%s", idea.Number, idea.Slug)
		for _, user := range users {
			if user.ID != c.User().ID {
				if _, err = c.Services().Notifications.Insert(user, title, link, idea.ID); err != nil {
					return c.Failure(err)
				}
			}
		}

		// Email notification
		users, err = c.Services().Ideas.GetActiveSubscribers(idea.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
		if err != nil {
			return c.Failure(err)
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

		return c.Services().Emailer.BatchSend(c.Tenant(), "new_comment", params, c.User().Name, to)
	})
}

//NotifyAboutStatusChange sends a notification (web and email) to subscribers
func NotifyAboutStatusChange(idea *models.Idea, prevStatus int) worker.Task {
	return describe("Notify about idea status change", func(c *worker.Context) error {
		//Don't notify if previous status is the same
		if prevStatus == idea.Status {
			return nil
		}

		// Web notification
		users, err := c.Services().Ideas.GetActiveSubscribers(idea.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("**%s** changed status of **%s** to **%s**", c.User().Name, idea.Title, models.GetIdeaStatusName(idea.Status))
		link := fmt.Sprintf("/ideas/%d/%s", idea.Number, idea.Slug)
		for _, user := range users {
			if user.ID != c.User().ID {
				if _, err = c.Services().Notifications.Insert(user, title, link, idea.ID); err != nil {
					return c.Failure(err)
				}
			}
		}

		// Email notification
		users, err = c.Services().Ideas.GetActiveSubscribers(idea.Number, models.NotificationChannelEmail, models.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		var duplicate template.HTML
		if idea.Status == models.IdeaDuplicate {
			originalIdea, err := c.Services().Ideas.GetByNumber(idea.Response.Original.Number)
			if err != nil {
				return c.Failure(err)
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
			"content":     markdown.Parse(idea.Response.Text),
			"status":      models.GetIdeaStatusName(idea.Status),
			"duplicate":   duplicate,
			"view":        linkWithText("View it on your browser", c.BaseURL(), "/ideas/%d/%s", idea.Number, idea.Slug),
			"unsubscribe": linkWithText("unsubscribe from it", c.BaseURL(), "/ideas/%d/%s", idea.Number, idea.Slug),
			"change":      linkWithText("change your notification settings", c.BaseURL(), "/settings"),
		}

		return c.Services().Emailer.BatchSend(c.Tenant(), "change_status", params, c.User().Name, to)
	})
}

//SendInvites sends one email to each invited recipient
func SendInvites(subject, message string, invitations []*models.UserInvitation) worker.Task {
	return describe("Send invites", func(c *worker.Context) error {
		to := make([]email.Recipient, len(invitations))
		for i, invite := range invitations {
			err := c.Services().Tenants.SaveVerificationKey(invite.VerificationKey, 15*24*time.Hour, invite)
			if err != nil {
				return c.Failure(err)
			}

			url := link(c.BaseURL(), "/invite/verify?k=%s", invite.VerificationKey)
			toMessage := strings.Replace(message, app.InvitePlaceholder, string(url), -1)
			to[i] = email.NewRecipient("", invite.Email, email.Params{
				"message": markdown.Parse(toMessage),
			})
		}
		return c.Services().Emailer.BatchSend(c.Tenant(), "invite_email", email.Params{
			"subject": subject,
		}, c.User().Name, to)
	})
}
