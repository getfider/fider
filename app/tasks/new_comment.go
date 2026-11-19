package tasks

import (
	"context"
	"fmt"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/errors"

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

// NotifyAboutNewComment sends a notification (web and email) to subscribers.
// A pending comment only tells moderators it needs approval. A comment on a pending post only notifies staff.
func NotifyAboutNewComment(comment *entity.Comment, post *entity.Post) worker.Task {
	if !comment.IsApproved {
		return notifyAwaitingApproval("Notify about comment awaiting approval", func() awaitingApproval { return commentAwaitingApproval(comment, post) })
	}
	return notifyAboutNewComment("Notify about new comment", comment, post, audienceFor(post.IsApproved), 0)
}

// NotifyAboutApprovedComment sends the notifications that were held back while the comment was pending,
// to everyone but the moderator who approved it
func NotifyAboutApprovedComment(comment *entity.Comment, post *entity.Post, approverID int) worker.Task {
	return notifyAboutNewComment("Notify about approved comment", comment, post, everyone, approverID)
}

func notifyAboutNewComment(name string, comment *entity.Comment, post *entity.Post, aud audience, skipUserID int) worker.Task {
	return describe(name, func(c *worker.Context) error {
		contentString := entity.CommentString(comment.Content)
		author := comment.User

		replyTo, err := getReplyTo(c, comment)
		if err != nil {
			return c.Failure(err)
		}

		rules := recipientRules{
			authorID:   author.ID,
			replyTo:    replyTo,
			mentions:   contentString.ParseMentions(),
			general:    &enum.NotificationEventNewComment,
			aud:        aud,
			skipUserID: skipUserID,
		}
		if replyTo != nil || len(rules.mentions) > 0 {
			mN := &query.GetMentionNotifications{CommentID: comment.ID}
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

		authorName, postTitle := markdown.Escape(author.Name), markdown.Escape(post.Title)
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		for _, n := range []struct {
			users []*entity.User
			title string
		}{
			{onWeb.reply, fmt.Sprintf("**%s** replied to your comment on **%s**", authorName, postTitle)},
			{onWeb.mention, fmt.Sprintf("**%s** mentioned you in **%s**", authorName, postTitle)},
			{onWeb.general, fmt.Sprintf("**%s** left a comment on **%s**", authorName, postTitle)},
		} {
			if err := addWebNotifications(c, n.users, n.title, link, post.ID); err != nil {
				return c.Failure(err)
			}
		}

		if err := logMentionNotifications(c, onWeb, byEmail, comment.ID, 0); err != nil {
			return c.Failure(err)
		}

		content := contentString.SanitizeMentions()
		sendEmailNotifications(c, author, post, toRecipients(byEmail.reply), content, "email.new_reply.text", "new_comment")
		sendEmailNotifications(c, author, post, toRecipients(byEmail.mention), content, "email.new_mention.text", "new_comment")
		sendEmailNotifications(c, author, post, toRecipients(byEmail.general), content, "email.new_comment.text", "new_comment")

		// Webhook targets aren't moderation-aware, so wait until the comment and its post are approved
		if !comment.IsApproved || !post.IsApproved {
			return nil
		}

		tenant := c.Tenant()
		baseURL, logoURL := web.BaseURL(c), web.LogoURL(c)

		webhookProps := webhook.Props{"comment": content}
		webhookProps["comment_id"] = comment.ID
		if comment.ParentID != nil {
			webhookProps["comment_parent_id"] = *comment.ParentID
		}
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

// getReplyTo returns the author of the comment being replied to, or nil if it isn't a reply or the parent is gone
func getReplyTo(ctx context.Context, comment *entity.Comment) (*entity.User, error) {
	if comment.ParentID == nil {
		return nil, nil
	}
	q := &query.GetCommentByID{CommentID: *comment.ParentID}
	if err := bus.Dispatch(ctx, q); err != nil {
		if errors.Cause(err) == app.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return q.Result.User, nil
}

func NotifyAboutUpdatedComment(post *entity.Post, comment *entity.Comment) worker.Task {
	return describe("Notify about updated comment", func(c *worker.Context) error {
		contentString := entity.CommentString(comment.Content)
		mentions := contentString.ParseMentions()

		log.Infof(c, "Comment updated: @{Comment:Yellow}. Mentions @{MentionsCount}", dto.Props{
			"Comment":       contentString,
			"MentionsCount": len(mentions),
		})
		// Nothing is sent about a pending comment until it's approved
		if len(mentions) == 0 || !comment.IsApproved {
			return nil
		}

		mN := &query.GetMentionNotifications{CommentID: comment.ID}
		if err := bus.Dispatch(c, mN); err != nil {
			return c.Failure(err)
		}

		author := c.User()
		rules := recipientRules{
			authorID:        author.ID,
			mentions:        mentions,
			alreadyNotified: mN.Result,
			aud:             audienceFor(post.IsApproved),
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
		if err := logMentionNotifications(c, onWeb, byEmail, comment.ID, 0); err != nil {
			return c.Failure(err)
		}

		sendEmailNotifications(c, author, post, toRecipients(byEmail.mention), contentString.SanitizeMentions(), "email.new_mention.text", "new_comment")

		return nil
	})
}

// sendEmailNotifications emails a comment (or post description) to recipients. messageLocaleString is the intro line.
func sendEmailNotifications(c *worker.Context, author *entity.User, post *entity.Post, to []dto.Recipient, comment string, messageLocaleString string, templateName string) {
	// Short circuit if there is no one to notify
	if len(to) == 0 {
		return
	}

	tenant := c.Tenant()
	baseURL, logoURL := web.BaseURL(c), web.LogoURL(c)

	mailProps := dto.Props{
		"title":               post.Title,
		"messageLocaleString": messageLocaleString,
		"siteName":            tenant.Name,
		"userName":            author.Name,
		"content":             markdown.Full(c, comment, false),
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
