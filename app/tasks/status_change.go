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
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/webhook"
	"github.com/getfider/fider/app/pkg/worker"
)

// NotifyAboutStatusChange sends a notification (web and email) to subscribers
// and fires the WebhookChangeStatus payload. Identity is the tenant-defined
// status slug; the payload's labels come from the live tenant catalogue.
func NotifyAboutStatusChange(post *entity.Post, prevStatusSlug string) worker.Task {
	return describe("Notify about post status change", func(c *worker.Context) error {
		if prevStatusSlug == post.StatusSlug {
			return nil
		}

		// Resolve human labels from the tenant catalogue. Built-in slugs go
		// through i18n for the email subject line; webhook receivers get the
		// admin-defined label verbatim (no locale catalog for custom statuses).
		statusList := &query.ListActiveStatusesForTenant{}
		_ = bus.Dispatch(c, statusList)
		labelFor := func(slug string, fallback string) string {
			for _, s := range statusList.Result {
				if s.Slug == slug {
					return s.Label
				}
			}
			return fallback
		}

		users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		newLabel := labelFor(post.StatusSlug, post.StatusSlug)
		oldLabel := labelFor(prevStatusSlug, prevStatusSlug)

		author := c.User()
		title := fmt.Sprintf("**%s** changed status of **%s** to **%s**", author.Name, post.Title, newLabel)
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

		users, err = getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		baseURL := web.BaseURL(c)
		var duplicate string
		if post.StatusSlug == "duplicate" && post.Response != nil && post.Response.Original != nil {
			duplicate = linkWithText(post.Response.Original.Title, baseURL, "/posts/%d/%s", post.Response.Original.Number, post.Response.Original.Slug)
		}

		to := make([]dto.Recipient, 0)
		for _, user := range users {
			if user.ID != author.ID {
				to = append(to, dto.NewRecipient(user.Name, user.Email, dto.Props{}))
			}
		}

		tenant := c.Tenant()
		logoURL := web.LogoURL(c)

		// Built-in slugs have locale entries; custom slugs use the admin label.
		statusForEmail := newLabel
		if isBuiltinSlug(post.StatusSlug) {
			statusForEmail = i18n.T(c, fmt.Sprintf("enum.poststatus.%s", post.StatusSlug))
		}

		var responseText string
		if post.Response != nil {
			responseText = post.Response.Text
		}

		props := dto.Props{
			"title":       post.Title,
			"postLink":    linkWithText(fmt.Sprintf("#%d", post.Number), baseURL, "/posts/%d/%s", post.Number, post.Slug),
			"siteName":    tenant.Name,
			"content":     markdown.Full(responseText, true),
			"status":      statusForEmail,
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

		webhookProps := webhook.Props{
			"post_old_status_slug":  prevStatusSlug,
			"post_old_status_label": oldLabel,
			"post_status_label":     newLabel,
		}
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

func isBuiltinSlug(slug string) bool {
	switch slug {
	case "open", "started", "completed", "declined", "planned", "duplicate", "deleted":
		return true
	}
	return false
}
