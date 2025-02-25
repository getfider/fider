package tasks

import (
	"fmt"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/cmd"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/env"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/i18n"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/markdown"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/webhook"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/worker"
)

// NotifyAboutStatusChange sends a notification (web and email) to subscribers
func NotifyAboutStatusChange(post *entity.Post, prevStatus enum.PostStatus) worker.Task {
	return describe("Notify about post status change", func(c *worker.Context) error {
		//Don't notify if previous status is the same
		if prevStatus == post.Status {
			return nil
		}

		// Web notification
		users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		author := c.User()
		title := fmt.Sprintf("**%s** changed status of **%s** to **%s**", author.Name, post.Title, post.Status.Name())
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		for _, user := range users {
			if user.ID != author.ID {
				err = bus.Dispatch(c, &cmd.AddNewNotification{
					User:   user,
					Title:  title,
					Link:   link,
					PostID: post.ID,
				})
				if err != nil {
					return c.Failure(err)
				}
			}
		}

		tenant := c.Tenant()
		logoURL := web.LogoURL(c)
		baseURL := web.BaseURL(c)

		// Email notification
		if !env.Config.Email.DisableEmailNotifications {
			users, err = getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventChangeStatus)
			if err != nil {
				return c.Failure(err)
			}

			var duplicate string
			if post.Status == enum.PostDuplicate {
				duplicate = linkWithText(post.Response.Original.Title, baseURL, "/posts/%d/%s", post.Response.Original.Number, post.Response.Original.Slug)
			}

			to := make([]dto.Recipient, 0)
			for _, user := range users {
				if user.ID != author.ID {
					to = append(to, dto.NewRecipient(user.Name, user.Email, dto.Props{}))
				}
			}

			props := dto.Props{
				"title":       post.Title,
				"postLink":    linkWithText(fmt.Sprintf("#%d", post.Number), baseURL, "/posts/%d/%s", post.Number, post.Slug),
				"siteName":    tenant.Name,
				"content":     markdown.Full(post.Response.Text),
				"status":      i18n.T(c, fmt.Sprintf("enum.poststatus.%s", post.Status.Name())),
				"duplicate":   duplicate,
				"view":        linkWithText(i18n.T(c, "email.subscription.view"), baseURL, "/posts/%d/%s", post.Number, post.Slug),
				"unsubscribe": linkWithText(i18n.T(c, "email.subscription.unsubscribe"), baseURL, "/posts/%d/%s", post.Number, post.Slug),
				"change":      linkWithText(i18n.T(c, "email.subscription.change"), baseURL, "/settings"),
				"logo":        logoURL,
			}

			bus.Publish(c, &cmd.SendMail{
				From:         dto.Recipient{Name: author.Name},
				To:           to,
				TemplateName: "change_status",
				Props:        props,
			})

		}

		webhookProps := webhook.Props{"post_old_status": prevStatus.Name()}
		webhookProps.SetPost(post, "post", baseURL, true, true)
		webhookProps.SetUser(author, "author")
		webhookProps.SetTenant(tenant, "tenant", baseURL, logoURL)

		err = bus.Dispatch(c, &cmd.TriggerWebhooks{
			Type:  enum.WebhookChangeStatus,
			Props: webhookProps,
		})
		if err != nil {
			return c.Failure(err)
		}
		return nil
	})
}
