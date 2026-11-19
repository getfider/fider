package tasks

import (
	"fmt"

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

// NotifyAboutNewPost sends a notification (web and email) to subscribers.
// A pending post only tells moderators it needs approval; everyone is notified once it's approved.
func NotifyAboutNewPost(post *entity.Post) worker.Task {
	if !post.IsApproved {
		return notifyAwaitingApproval("Notify about post awaiting approval", func() awaitingApproval { return postAwaitingApproval(post) })
	}
	return notifyAboutNewPost("Notify about new post", post, 0)
}

// NotifyAboutApprovedPost sends the notifications that were held back while the post was pending,
// to everyone but the moderator who approved it
func NotifyAboutApprovedPost(post *entity.Post, approverID int) worker.Task {
	return notifyAboutNewPost("Notify about approved post", post, approverID)
}

func notifyAboutNewPost(name string, post *entity.Post, skipUserID int) worker.Task {
	return describe(name, func(c *worker.Context) error {
		contentString := entity.CommentString(post.Description)
		author := post.User

		rules := recipientRules{
			authorID:   author.ID,
			mentions:   contentString.ParseMentions(),
			general:    &enum.NotificationEventNewPost,
			aud:        everyone,
			skipUserID: skipUserID,
		}
		if len(rules.mentions) > 0 {
			mN := &query.GetMentionNotifications{PostID: post.ID}
			if err := bus.Dispatch(c, mN); err != nil {
				return c.Failure(err)
			}
			rules.alreadyNotified = mN.Result
		}
		onWeb, err := pickRecipients(c, post, enum.NotificationChannelWeb, rules)
		if err != nil {
			return c.Failure(err)
		}
		byEmail, err := pickRecipients(c, post, enum.NotificationChannelEmail, rules)
		if err != nil {
			return c.Failure(err)
		}

		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		mentionTitle := fmt.Sprintf("**%s** mentioned you in **%s**", markdown.Escape(author.Name), markdown.Escape(post.Title))
		if err := addWebNotifications(c, onWeb.mention, mentionTitle, link, post.ID); err != nil {
			return c.Failure(err)
		}
		if err := addWebNotifications(c, onWeb.general, fmt.Sprintf("New post: **%s**", markdown.Escape(post.Title)), link, post.ID); err != nil {
			return c.Failure(err)
		}
		if err := logMentionNotifications(c, onWeb, byEmail, 0, post.ID); err != nil {
			return c.Failure(err)
		}

		tenant := c.Tenant()
		baseURL, logoURL := web.BaseURL(c), web.LogoURL(c)

		content := contentString.SanitizeMentions()
		sendEmailNotifications(c, author, post, toRecipients(byEmail.mention), content, "email.new_mention.text", "new_comment")

		if len(byEmail.general) > 0 {
			bus.Publish(c, &cmd.SendMail{
				From:         dto.Recipient{Name: author.Name},
				To:           toRecipients(byEmail.general),
				TemplateName: "new_post",
				Props: dto.Props{
					"title":               post.Title,
					"messageLocaleString": "email.new_post.text",
					"siteName":            tenant.Name,
					"userName":            author.Name,
					"content":             markdown.Full(c, content, false),
					"postLink":            linkWithText(fmt.Sprintf("#%d", post.Number), baseURL, "/posts/%d/%s", post.Number, post.Slug),
					"view":                linkWithText(i18n.T(c, "email.subscription.view"), baseURL, "/posts/%d/%s", post.Number, post.Slug),
					"change":              linkWithText(i18n.T(c, "email.subscription.change"), baseURL, "/settings"),
					"logo":                logoURL,
				},
			})
		}

		// Webhook targets aren't moderation-aware, so wait until the post is approved
		if !post.IsApproved {
			return nil
		}

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

		log.Infof(c, "Post updated: @{Post:Yellow}. Mentions @{MentionsCount}", dto.Props{
			"Post":          post.Title,
			"MentionsCount": len(mentions),
		})
		// Nothing is sent about a pending post until it's approved
		if len(mentions) == 0 || !post.IsApproved {
			return nil
		}

		mN := &query.GetMentionNotifications{PostID: post.ID}
		if err := bus.Dispatch(c, mN); err != nil {
			return c.Failure(err)
		}

		author := c.User()
		rules := recipientRules{
			authorID:        author.ID,
			mentions:        mentions,
			alreadyNotified: mN.Result,
			aud:             everyone,
		}
		onWeb, err := pickRecipients(c, post, enum.NotificationChannelWeb, rules)
		if err != nil {
			return c.Failure(err)
		}
		byEmail, err := pickRecipients(c, post, enum.NotificationChannelEmail, rules)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("**%s** mentioned you in **%s**", markdown.Escape(author.Name), markdown.Escape(post.Title))
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		if err := addWebNotifications(c, onWeb.mention, title, link, post.ID); err != nil {
			return c.Failure(err)
		}
		if err := logMentionNotifications(c, onWeb, byEmail, 0, post.ID); err != nil {
			return c.Failure(err)
		}

		sendEmailNotifications(c, author, post, toRecipients(byEmail.mention), contentString.SanitizeMentions(), "email.new_mention.text", "new_comment")

		return nil
	})
}
