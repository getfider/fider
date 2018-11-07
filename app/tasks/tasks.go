package tasks

import (
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/worker"
)

var mentionsPattern = regexp.MustCompile(`\[(\d+)\]`)

func describe(name string, job worker.Job) worker.Task {
	return worker.Task{Name: name, Job: job}
}

func link(baseURL, path string, args ...interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("<a href='%[1]s%[2]s'>%[1]s%[2]s</a>", baseURL, fmt.Sprintf(path, args...)))
}

func linkWithText(text, baseURL, path string, args ...interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("<a href='%s%s'>%s</a>", baseURL, fmt.Sprintf(path, args...), text))
}

func replaceMentions(text string, c *worker.Context) string {
	return mentionsPattern.ReplaceAllStringFunc(text, func(match string) string {
		r := []rune(match)
		i, err := strconv.ParseInt(string(r[1:len(r)-1]), 10, 32)
		if err != nil {
			return "???"
		} else {
			user, err2 := c.Services().Users.GetByID(int(i))
			if err2 == nil {
				return user.Name
			} else {
				return "???"
			}
		}
	})
}

//SendSignUpEmail is used to send the sign up email to requestor
func SendSignUpEmail(model *models.CreateTenant, baseURL string) worker.Task {
	return describe("Send sign up email", func(c *worker.Context) error {
		to := email.NewRecipient(model.Name, model.Email, email.Params{
			"logo": "https://getfider.com/images/logo-100x100.png",
			"link": link(baseURL, "/signup/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send(c, "signup_email", email.Params{}, "Fider", to)
	})
}

//SendSignInEmail is used to send the sign in email to requestor
func SendSignInEmail(model *models.SignInByEmail) worker.Task {
	return describe("Send sign in email", func(c *worker.Context) error {
		to := email.NewRecipient("", model.Email, email.Params{
			"tenantName": c.Tenant().Name,
			"link":       link(c.BaseURL(), "/signin/verify?k=%s", model.VerificationKey),
		})
		return c.Services().Emailer.Send(c, "signin_email", email.Params{}, c.Tenant().Name, to)
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
		return c.Services().Emailer.Send(c, "change_emailaddress_email", email.Params{}, c.Tenant().Name, to)
	})
}

//NotifyAboutNewPost sends a notification (web and email) to subscribers
func NotifyAboutNewPost(post *models.Post) worker.Task {
	return describe("Notify about new post", func(c *worker.Context) error {
		// Web notification
		users, err := c.Services().Posts.GetActiveSubscribers(post.Number, models.NotificationChannelWeb, models.NotificationEventNewPost)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("New post: **%s**", post.Title)
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		for _, user := range users {
			if user.ID != c.User().ID {
				if _, err = c.Services().Notifications.Insert(user, title, link, post.ID); err != nil {
					return c.Failure(err)
				}
			}
		}

		// Email notification
		users, err = c.Services().Posts.GetActiveSubscribers(post.Number, models.NotificationChannelEmail, models.NotificationEventNewPost)
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
			"title":      post.Title,
			"tenantName": c.Tenant().Name,
			"userName":   c.User().Name,
			"content":    markdown.Parse(post.Description),
			"postLink":   linkWithText(fmt.Sprintf("#%d", post.Number), c.BaseURL(), "/posts/%d/%s", post.Number, post.Slug),
			"view":       linkWithText("View it on your browser", c.BaseURL(), "/posts/%d/%s", post.Number, post.Slug),
			"change":     linkWithText("change your notification settings", c.BaseURL(), "/settings"),
		}

		return c.Services().Emailer.BatchSend(c, "new_post", params, c.User().Name, to)
	})
}

//NotifyAboutNewComment sends a notification (web and email) to subscribers
func NotifyAboutNewComment(post *models.Post, comment *models.NewComment) worker.Task {
	return describe("Notify about new comment", func(c *worker.Context) error {

		title := fmt.Sprintf("**%s** left a comment on **%s**", c.User().Name, post.Title)
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		// Mention web notifications
		mentionsMatch := mentionsPattern.FindAllStringSubmatch(comment.Content, -1)

		for _, match := range mentionsMatch {
			id, _ := strconv.ParseInt(match[1], 10, 32)
			user, _ := c.Services().Users.GetByID(int(id))

			if user.ID != c.User().ID {
				title := fmt.Sprintf("**%s** mentioned you", c.User().Name)
				_, _ = c.Services().Notifications.Insert(user, title, link, post.ID)
			}
		}

		// Here we will also send only one notification, as the second insert will fail

		// Web notification
		users, err := c.Services().Posts.GetActiveSubscribers(post.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
		if err != nil {
			return c.Failure(err)
		}

		for _, user := range users {
			if user.ID != c.User().ID {
				c.Services().Notifications.Insert(user, title, link, post.ID)
			}
		}

		// Email notification
		users, err = c.Services().Posts.GetActiveSubscribers(post.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
		if err != nil {
			return c.Failure(err)
		}

		toMap := make(map[string]email.Recipient)
		for _, user := range users {
			if user.ID != c.User().ID {
				toMap[user.Email] = email.NewRecipient(user.Name, user.Email, email.Params{})
			}
		}

		// Mentions Email notifications
		for _, match := range mentionsMatch {
			id, _ := strconv.ParseInt(match[1], 10, 32)
			user, _ := c.Services().Users.GetByID(int(id))

			if user.ID != c.User().ID {
				toMap[user.Email] = email.NewRecipient(user.Name, user.Email, email.Params{"mention": true})
			}
		}

		params := email.Params{
			"title":       post.Title,
			"tenantName":  c.Tenant().Name,
			"userName":    c.User().Name,
			"content":     markdown.Parse(replaceMentions(comment.Content, c)),
			"postLink":    linkWithText(fmt.Sprintf("#%d", post.Number), c.BaseURL(), "/posts/%d/%s", post.Number, post.Slug),
			"view":        linkWithText("View it on your browser", c.BaseURL(), "/posts/%d/%s", post.Number, post.Slug),
			"unsubscribe": linkWithText("unsubscribe from it", c.BaseURL(), "/posts/%d/%s", post.Number, post.Slug),
			"change":      linkWithText("change your notification settings", c.BaseURL(), "/settings"),
		}

		// This is done so that we don't send two emails to a user for the same post, the mention one will prevail
		to := make([]email.Recipient, 0)
		for _, value := range toMap {
			to = append(to, value)
		}

		return c.Services().Emailer.BatchSend(c, "new_comment", params, c.User().Name, to)
	})
}

//NotifyAboutStatusChange sends a notification (web and email) to subscribers
func NotifyAboutStatusChange(post *models.Post, prevStatus models.PostStatus) worker.Task {
	return describe("Notify about post status change", func(c *worker.Context) error {
		//Don't notify if previous status is the same
		if prevStatus == post.Status {
			return nil
		}

		// Web notification
		users, err := c.Services().Posts.GetActiveSubscribers(post.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("**%s** changed status of **%s** to **%s**", c.User().Name, post.Title, post.Status.Name())
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		for _, user := range users {
			if user.ID != c.User().ID {
				if _, err = c.Services().Notifications.Insert(user, title, link, post.ID); err != nil {
					return c.Failure(err)
				}
			}
		}

		// Email notification
		users, err = c.Services().Posts.GetActiveSubscribers(post.Number, models.NotificationChannelEmail, models.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		var duplicate template.HTML
		if post.Status == models.PostDuplicate {
			originalPost, err := c.Services().Posts.GetByNumber(post.Response.Original.Number)
			if err != nil {
				return c.Failure(err)
			}
			duplicate = linkWithText(originalPost.Title, c.BaseURL(), "/posts/%d/%s", originalPost.Number, originalPost.Slug)
		}

		to := make([]email.Recipient, 0)
		for _, user := range users {
			if user.ID != c.User().ID {
				to = append(to, email.NewRecipient(user.Name, user.Email, email.Params{}))
			}
		}

		params := email.Params{
			"title":       post.Title,
			"postLink":    linkWithText(fmt.Sprintf("#%d", post.Number), c.BaseURL(), "/posts/%d/%s", post.Number, post.Slug),
			"tenantName":  c.Tenant().Name,
			"content":     markdown.Parse(post.Response.Text),
			"status":      post.Status.Name(),
			"duplicate":   duplicate,
			"view":        linkWithText("View it on your browser", c.BaseURL(), "/posts/%d/%s", post.Number, post.Slug),
			"unsubscribe": linkWithText("unsubscribe from it", c.BaseURL(), "/posts/%d/%s", post.Number, post.Slug),
			"change":      linkWithText("change your notification settings", c.BaseURL(), "/settings"),
		}

		return c.Services().Emailer.BatchSend(c, "change_status", params, c.User().Name, to)
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

			url := fmt.Sprintf("%s/invite/verify?k=%s", c.BaseURL(), invite.VerificationKey)
			toMessage := strings.Replace(message, app.InvitePlaceholder, string(url), -1)
			to[i] = email.NewRecipient("", invite.Email, email.Params{
				"message": markdown.Parse(toMessage),
			})
		}
		return c.Services().Emailer.BatchSend(c, "invite_email", email.Params{
			"subject": subject,
		}, c.User().Name, to)
	})
}
