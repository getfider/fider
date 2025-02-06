package tasks

import (
	"fmt"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/cmd"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/i18n"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/markdown"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/webhook"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/worker"
)

// NotifyAboutNewPost sends a notification (web and email) to subscribers
func NotifyAboutNewPost(post *entity.Post) worker.Task {
	return describe("Notify about new post", func(c *worker.Context) error {
		// Web notification
		users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventNewPost)
		if err != nil {
			return c.Failure(err)
		}

		author := c.User()
		title := fmt.Sprintf("New post: **%s**", post.Title)
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

		// Email notification
		users, err = getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventNewPost)
		if err != nil {
			return c.Failure(err)
		}

		to := make([]dto.Recipient, 0)
		for _, user := range users {
			if user.ID != author.ID {
				to = append(to, dto.NewRecipient(user.Name, user.Email, dto.Props{}))
			}
		}

		tenant := c.Tenant()
		baseURL, logoURL := web.BaseURL(c), web.LogoURL(c)

		mailProps := dto.Props{
			"title":    post.Title,
			"siteName": tenant.Name,
			"userName": author.Name,
			"content":  markdown.Full(post.Description),
			"postLink": linkWithText(fmt.Sprintf("#%d", post.Number), baseURL, "/posts/%d/%s", post.Number, post.Slug),
			"view":     linkWithText(i18n.T(c, "email.subscription.view"), baseURL, "/posts/%d/%s", post.Number, post.Slug),
			"change":   linkWithText(i18n.T(c, "email.subscription.change"), baseURL, "/settings"),
			"logo":     logoURL,
		}

		bus.Publish(c, &cmd.SendMail{
			From:         dto.Recipient{Name: author.Name},
			To:           to,
			TemplateName: "new_post",
			Props:        mailProps,
		})

		webhookProps := webhook.Props{}
		webhookProps.SetPost(post, "post", baseURL, false, false)
		webhookProps.SetUser(author, "author")
		webhookProps.SetTenant(tenant, "tenant", baseURL, logoURL)

		err = bus.Dispatch(c, &cmd.TriggerWebhooks{
			Type:  enum.WebhookNewPost,
			Props: webhookProps,
		})
		if err != nil {
			return c.Failure(err)
		}

		return nil
	})
}
