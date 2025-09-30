package tasks

import (
	"fmt"
	"slices"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/i18n"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/webhook"
	"github.com/getfider/fider/app/pkg/worker"
)

// NotifyAboutNewComment sends a notification (web and email) to subscribers
func NotifyAboutNewComment(comment *entity.Comment, post *entity.Post) worker.Task {
	return describe("Notify about new comment", func(c *worker.Context) error {

		// comment.ParseMentions()
		contentString := entity.CommentString(comment.Content)
		mentions := contentString.ParseMentions()
		var mentionNotifications []*entity.MentionNotification

		// Web notification
		users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventNewComment)
		if err != nil {
			return c.Failure(err)
		}

		author := c.User()
		title := fmt.Sprintf("**%s** left a comment on **%s**", author.Name, post.Title)
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

		// Web notification - mentions
		title = fmt.Sprintf("**%s** mentioned you in **%s**", author.Name, post.Title)

		if mentions != nil {

			users, err = getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventMention)
			if err != nil {
				return c.Failure(err)
			}

			// Get the existing mentions that have been sent for this comment
			mN := &query.GetMentionNotifications{
				CommentID: comment.ID,
			}
			err := bus.Dispatch(c, mN)
			if err != nil {
				return c.Failure(err)
			}
			mentionNotifications = mN.Result

			// Iterate the mentions
			for _, mention := range mentions {
				// Check if the user is in the list of mention subscribers (users)
				for _, u := range users {

					if u.Name == mention && !slices.ContainsFunc(mentionNotifications,
						func(n *entity.MentionNotification) bool {
							return n.UserID == u.ID
						}) {
						err = bus.Dispatch(c, &cmd.AddNewNotification{
							User:   u,
							Title:  title,
							Link:   link,
							PostID: post.ID,
						})
						if err != nil {
							return c.Failure(err)
						}

						// Also send the notification log
						err = bus.Dispatch(c, &cmd.AddMentionNotification{
							UserID:    u.ID,
							CommentID: comment.ID,
						})
						if err != nil {
							return c.Failure(err)
						}

					}

				}
			}

		}

		// Standard email notitifications
		users, err = getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventNewComment)
		if err != nil {
			return c.Failure(err)
		}

		to := make([]dto.Recipient, 0)
		for _, user := range users {
			if user.ID != author.ID {
				to = append(to, dto.NewRecipient(user.Name, user.Email, dto.Props{}))
			}
		}

		sendEmailNotifications(c, post, to, contentString.SanitizeMentions(), enum.NotificationEventNewComment, "new_comment")

		// Mentions
		to = make([]dto.Recipient, 0)
		if mentions != nil {

			users, err = getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventMention)
			if err != nil {
				return c.Failure(err)
			}

			for _, mention := range mentions {
				for _, u := range users {

					if u.Name == mention && !slices.ContainsFunc(mentionNotifications,
						func(n *entity.MentionNotification) bool {
							return n.UserID == u.ID
						}) {
						to = append(to, dto.NewRecipient(u.Name, u.Email, dto.Props{}))

						// Also send the notification log
						err = bus.Dispatch(c, &cmd.AddMentionNotification{
							UserID:    u.ID,
							CommentID: comment.ID,
						})
						if err != nil {
							return c.Failure(err)
						}
					}
				}
			}

		}

		sendEmailNotifications(c, post, to, contentString.SanitizeMentions(), enum.NotificationEventMention, "new_comment")

		tenant := c.Tenant()
		baseURL, logoURL := web.BaseURL(c), web.LogoURL(c)

		webhookProps := webhook.Props{"comment": contentString.SanitizeMentions()}
		webhookProps["comment_id"] = comment.ID
		webhookProps.SetPost(post, "post", baseURL, true, true)
		webhookProps.SetUser(author, "author")
		webhookProps.SetTenant(tenant, "tenant", baseURL, logoURL)

		err = bus.Dispatch(c, &cmd.TriggerWebhooks{
			Type:  enum.WebhookNewComment,
			Props: webhookProps,
		})
		if err != nil {
			return c.Failure(err)
		}

		return nil
	})
}

func NotifyAboutUpdatedComment(post *entity.Post, comment *entity.Comment) worker.Task {
	return describe("Notify about updated comment", func(c *worker.Context) error {

		contentString := entity.CommentString(comment.Content)
		mentions := contentString.ParseMentions()
		var mentionNotifications []*entity.MentionNotification

		log.Infof(c, "Comment updated: @{Comment:Yellow}. Mentions @{MentionsCount}", dto.Props{
			"Comment":       contentString,
			"MentionsCount": len(mentions),
		})

		author := c.User()
		title := fmt.Sprintf("**%s** mentioned you in **%s**", author.Name, post.Title)
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		mentionNotificationSent := false
		if mentions != nil {

			users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventMention)
			if err != nil {
				return c.Failure(err)
			}

			// Get the existing mentions that have been sent for this comment
			mN := &query.GetMentionNotifications{
				CommentID: comment.ID,
			}
			err = bus.Dispatch(c, mN)
			if err != nil {
				return c.Failure(err)
			}
			mentionNotifications = mN.Result

			// Iterate the mentions
			for _, mention := range mentions {
				// Check if the user is in the list of mention subscribers (users)
				for _, u := range users {
					if u.Name == mention && !slices.ContainsFunc(mentionNotifications,
						func(n *entity.MentionNotification) bool {
							return n.UserID == u.ID
						}) {
						err = bus.Dispatch(c, &cmd.AddNewNotification{
							User:   u,
							Title:  title,
							Link:   link,
							PostID: post.ID,
						})
						if err != nil {
							return c.Failure(err)
						}

						// Also send the notification log
						err = bus.Dispatch(c, &cmd.AddMentionNotification{
							UserID:    u.ID,
							CommentID: comment.ID,
						})
						if err != nil {
							return c.Failure(err)
						}
						mentionNotificationSent = true
					}
				}
			}

		}

		to := make([]dto.Recipient, 0)
		if mentions != nil {

			users, err := getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventMention)
			if err != nil {
				return c.Failure(err)
			}

			for _, mention := range mentions {
				// Check if the user is in the list of mention subscribers (users)
				for _, u := range users {
					if u.Name == mention && !slices.ContainsFunc(mentionNotifications,
						func(n *entity.MentionNotification) bool {
							return n.UserID == u.ID
						}) {
						to = append(to, dto.NewRecipient(u.Name, u.Email, dto.Props{}))

						// Also send the notification log
						if !mentionNotificationSent {
							err = bus.Dispatch(c, &cmd.AddMentionNotification{
								UserID:    u.ID,
								CommentID: comment.ID,
							})
							if err != nil {
								return c.Failure(err)
							}
						}
					}
				}
			}
		}

		sendEmailNotifications(c, post, to, contentString.SanitizeMentions(), enum.NotificationEventMention, "new_comment")

		return nil
	})
}

func sendEmailNotifications(c *worker.Context, post *entity.Post, to []dto.Recipient, comment string, event enum.NotificationEvent, templateName string) {
	// Short circuit if there is no one to notify
	if len(to) == 0 {
		return
	}

	author := c.User()
	tenant := c.Tenant()
	baseURL, logoURL := web.BaseURL(c), web.LogoURL(c)
	messaleLocaleString := "email.new_comment.text"
	if event.UserSettingsKeyName == enum.NotificationEventMention.UserSettingsKeyName {
		messaleLocaleString = "email.new_mention.text"
	}

	mailProps := dto.Props{
		"title":               post.Title,
		"messageLocaleString": messaleLocaleString,
		"siteName":            tenant.Name,
		"userName":            author.Name,
		"content":             markdown.Full(comment, false),
		"postLink":            linkWithText(fmt.Sprintf("#%d", post.Number), baseURL, "/posts/%d/%s", post.Number, post.Slug),
		"view":                linkWithText(i18n.T(c, "email.subscription.view"), baseURL, "/posts/%d/%s", post.Number, post.Slug),
		"unsubscribe":         linkWithText(i18n.T(c, "email.subscription.unsubscribe"), baseURL, "/posts/%d/%s", post.Number, post.Slug),
		"change":              linkWithText(i18n.T(c, "email.subscription.change"), baseURL, "/settings"),
		"logo":                logoURL,
	}

	bus.Publish(c, &cmd.SendMail{
		From:         dto.Recipient{Name: author.Name},
		To:           to,
		TemplateName: templateName,
		Props:        mailProps,
	})
}
