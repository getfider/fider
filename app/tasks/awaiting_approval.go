package tasks

import (
	"fmt"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/worker"
)

// awaitingApproval is a post or comment waiting for moderation
type awaitingApproval struct {
	post     *entity.Post
	author   *entity.User
	content  string // markdown shown in the email
	title    string // web notification title
	link     string
	emailKey string // email intro line
	template string
}

// notifyAwaitingApproval tells staff that can moderate about pending content, instead of the usual notifications.
// describeContent runs when the task does, like the other tasks that read the post and author.
func notifyAwaitingApproval(name string, describeContent func() awaitingApproval) worker.Task {
	return describe(name, func(c *worker.Context) error {
		a := describeContent()
		rules := recipientRules{authorID: a.author.ID, general: &enum.NotificationEventAwaitingApproval, aud: staffOnly}

		onWeb, err := pickRecipients(c, a.post, enum.NotificationChannelWeb, rules)
		if err != nil {
			return c.Failure(err)
		}
		if err := addWebNotifications(c, onWeb.general, a.title, a.link, a.post.ID); err != nil {
			return c.Failure(err)
		}

		byEmail, err := pickRecipients(c, a.post, enum.NotificationChannelEmail, rules)
		if err != nil {
			return c.Failure(err)
		}
		sendEmailNotifications(c, a.author, a.post, toRecipients(byEmail.general), a.content, a.emailKey, a.template)

		return nil
	})
}

func postAwaitingApproval(post *entity.Post) awaitingApproval {
	return awaitingApproval{
		post:     post,
		author:   post.User,
		content:  entity.CommentString(post.Description).SanitizeMentions(),
		title:    fmt.Sprintf("**%s** created a new post for approval: **%s**", markdown.Escape(post.User.Name), markdown.Escape(post.Title)),
		link:     fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug),
		emailKey: "email.new_post_approval.text",
		template: "new_post",
	}
}

func commentAwaitingApproval(comment *entity.Comment, post *entity.Post) awaitingApproval {
	return awaitingApproval{
		post:     post,
		author:   comment.User,
		content:  entity.CommentString(comment.Content).SanitizeMentions(),
		title:    fmt.Sprintf("**%s** left a comment that needs approval on **%s**", markdown.Escape(comment.User.Name), markdown.Escape(post.Title)),
		link:     fmt.Sprintf("/posts/%d/%s#comment-%d", post.Number, post.Slug, comment.ID),
		emailKey: "email.new_comment_approval.text",
		template: "new_comment",
	}
}
