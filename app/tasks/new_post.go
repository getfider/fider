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
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/webhook"
	"github.com/getfider/fider/app/pkg/worker"
)

// NotifyAboutNewPost sends a notification (web and email) to subscribers
func NotifyAboutNewPost(post *entity.Post) worker.Task {
	return describe("Notify about new post", func(c *worker.Context) error {
		// Parse mentions from post description
		contentString := entity.CommentString(post.Description)
		mentions := contentString.ParseMentions()
		var mentionNotifications []*entity.MentionNotification

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

		// Web notification - mentions
		if len(mentions) > 0 {
			title = fmt.Sprintf("**%s** mentioned you in **%s**", author.Name, post.Title)

			users, err = getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventMention)
			if err != nil {
				return c.Failure(err)
			}

			// Get the existing mentions that have been sent for this post
			mN := &query.GetMentionNotifications{
				PostID: post.ID,
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
							UserID: u.ID,
							PostID: post.ID,
						})
						if err != nil {
							return c.Failure(err)
						}
					}
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

		sendEmailNotifications(c, post, to, contentString.SanitizeMentions(), enum.NotificationEventNewComment, "new_post")

		// Email notification - mentions
		to = make([]dto.Recipient, 0)
		if len(mentions) > 0 {
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
							UserID: u.ID,
							PostID: post.ID,
						})
						if err != nil {
							return c.Failure(err)
						}
					}
				}
			}
		}

		// Send mention email notifications
		sendEmailNotifications(c, post, to, contentString.SanitizeMentions(), enum.NotificationEventMention, "new_comment")

		tenant := c.Tenant()
		baseURL, logoURL := web.BaseURL(c), web.LogoURL(c)

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

// NotifyAboutUpdatedPost sends notifications about mentions in an updated post
func NotifyAboutUpdatedPost(post *entity.Post) worker.Task {
	return describe("Notify about updated post", func(c *worker.Context) error {
		contentString := entity.CommentString(post.Description)
		mentions := contentString.ParseMentions()
		var mentionNotifications []*entity.MentionNotification

		log.Infof(c, "Post updated: @{Post:Yellow}. Mentions @{MentionsCount}", dto.Props{
			"Post":          post.Title,
			"MentionsCount": len(mentions),
		})

		author := c.User()
		title := fmt.Sprintf("**%s** mentioned you in **%s**", author.Name, post.Title)
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		mentionNotificationSent := false
		if len(mentions) > 0 {
			users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventMention)
			if err != nil {
				return c.Failure(err)
			}

			// Get the existing mentions that have been sent for this post
			mN := &query.GetMentionNotifications{
				PostID: post.ID,
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
							UserID: u.ID,
							PostID: post.ID,
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
		if len(mentions) > 0 {
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
								UserID: u.ID,
								PostID: post.ID,
							})
							if err != nil {
								return c.Failure(err)
							}
						}
					}
				}
			}
		}

		// Send email notifications for mentions
		sendEmailNotifications(c, post, to, contentString.SanitizeMentions(), enum.NotificationEventMention, "new_comment")

		return nil
	})
}
