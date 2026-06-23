package tasks_test

import (
	"context"
	"html/template"
	"testing"
	"time"

	"github.com/getfider/fider/app/pkg/webhook"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/services/email/emailmock"
	"github.com/getfider/fider/app/tasks"
)

func TestNotifyAboutStatusChangeTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{
			mock.AryaStark,
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveStatusesForTenant) error {
		q.Result = []*entity.Status{
			{Slug: "open", Label: "Open", Kind: "open"},
			{Slug: "planned", Label: "Planned", Kind: "active"},
		}
		return nil
	})

	var triggerWebhooks *cmd.TriggerWebhooks
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooks) error {
		triggerWebhooks = c
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "Add support for TypeScript",
		Slug:        "add-support-for-typescript",
		Description: "TypeScript is great, please add support for it",
		User:        mock.AryaStark,
		StatusSlug:  "planned",
		StatusKind:  "active",
		Response: &entity.PostResponse{
			RespondedAt: time.Now(),
			Text:        "Planned for next release.",
			User:        mock.JonSnow,
		},
	}

	task := tasks.NotifyAboutStatusChange(post, "open")

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("change_status")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":       "Add support for TypeScript",
		"postLink":    "<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>",
		"siteName":    "Demonstration",
		"content":     template.HTML("<p>Planned for next release.</p>"),
		"duplicate":   "",
		"status":      "Planned",
		"view":        "<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>",
		"change":      "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"unsubscribe": "<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>",
		"logo":        "https://login.fider.io/static/assets/logo.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(dto.Recipient{
		Name: "Jon Snow",
	})
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/1/add-support-for-typescript")
	Expect(addNewNotification.Title).Equals("**Jon Snow** changed status of **Add support for TypeScript** to **Planned**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)

	Expect(triggerWebhooks).IsNotNil()
	Expect(triggerWebhooks.Type).Equals(enum.WebhookChangeStatus)
	Expect(triggerWebhooks.Props).ContainsProps(webhook.Props{
		"post_old_status_slug":       "open",
		"post_old_status_label":      "Open",
		"post_status_label":          "Planned",
		"post_id":                    post.ID,
		"post_number":                post.Number,
		"post_title":                 post.Title,
		"post_slug":                  post.Slug,
		"post_description":           post.Description,
		"post_status_slug":           post.StatusSlug,
		"post_status_kind":           post.StatusKind,
		"post_url":                   "http://domain.com/posts/1/add-support-for-typescript",
		"post_author_id":             mock.AryaStark.ID,
		"post_author_name":           mock.AryaStark.Name,
		"post_author_email":          mock.AryaStark.Email,
		"post_author_role":           mock.AryaStark.Role.String(),
		"post_response":              true,
		"post_response_text":         post.Response.Text,
		"post_response_responded_at": post.Response.RespondedAt,
		"post_response_author_id":    mock.JonSnow.ID,
		"post_response_author_name":  mock.JonSnow.Name,
		"post_response_author_email": mock.JonSnow.Email,
		"post_response_author_role":  mock.JonSnow.Role.String(),
		"author_id":                  mock.JonSnow.ID,
		"author_name":                mock.JonSnow.Name,
		"author_email":               mock.JonSnow.Email,
		"author_role":                mock.JonSnow.Role.String(),
		"tenant_id":                  mock.DemoTenant.ID,
		"tenant_name":                mock.DemoTenant.Name,
		"tenant_subdomain":           mock.DemoTenant.Subdomain,
		"tenant_status":              mock.DemoTenant.Status.String(),
		"tenant_url":                 "http://domain.com",
	})
}

func TestNotifyAboutStatusChangeTask_Duplicate(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{
			mock.AryaStark,
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveStatusesForTenant) error {
		q.Result = []*entity.Status{
			{Slug: "open", Label: "Open", Kind: "open"},
			{Slug: "duplicate", Label: "Duplicate", Kind: "duplicate"},
		}
		return nil
	})

	var triggerWebhooks *cmd.TriggerWebhooks
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooks) error {
		triggerWebhooks = c
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:         2,
		Number:     2,
		Title:      "I need TypeScript",
		Slug:       "i-need-typescript",
		User:       mock.AryaStark,
		StatusSlug: "duplicate",
		StatusKind: "duplicate",
		Response: &entity.PostResponse{
			RespondedAt: time.Now(),
			User:        mock.JonSnow,
			Original: &entity.OriginalPost{
				Number:     1,
				Title:      "Add support for TypeScript",
				Slug:       "add-support-for-typescript",
				StatusSlug: "open",
			},
		},
	}

	task := tasks.NotifyAboutStatusChange(post, "open")

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("change_status")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":       "I need TypeScript",
		"postLink":    "<a href='http://domain.com/posts/2/i-need-typescript'>#2</a>",
		"siteName":    "Demonstration",
		"content":     template.HTML(""),
		"duplicate":   "<a href='http://domain.com/posts/1/add-support-for-typescript'>Add support for TypeScript</a>",
		"status":      "Duplicate",
		"view":        "<a href='http://domain.com/posts/2/i-need-typescript'>view it on your browser</a>",
		"change":      "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"unsubscribe": "<a href='http://domain.com/posts/2/i-need-typescript'>unsubscribe from it</a>",
		"logo":        "https://login.fider.io/static/assets/logo.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(dto.Recipient{
		Name: "Jon Snow",
	})
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/2/i-need-typescript")
	Expect(addNewNotification.Title).Equals("**Jon Snow** changed status of **I need TypeScript** to **Duplicate**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)

	Expect(triggerWebhooks).IsNotNil()
	Expect(triggerWebhooks.Type).Equals(enum.WebhookChangeStatus)
	Expect(triggerWebhooks.Props).ContainsProps(webhook.Props{
		"post_old_status_slug":          "open",
		"post_old_status_label":         "Open",
		"post_status_label":             "Duplicate",
		"post_id":                       post.ID,
		"post_number":                   post.Number,
		"post_title":                    post.Title,
		"post_slug":                     post.Slug,
		"post_description":              post.Description,
		"post_status_slug":              post.StatusSlug,
		"post_status_kind":              post.StatusKind,
		"post_url":                      "http://domain.com/posts/2/i-need-typescript",
		"post_author_id":                mock.AryaStark.ID,
		"post_author_name":              mock.AryaStark.Name,
		"post_author_email":             mock.AryaStark.Email,
		"post_author_role":              mock.AryaStark.Role.String(),
		"post_response":                 true,
		"post_response_text":            post.Response.Text,
		"post_response_responded_at":    post.Response.RespondedAt,
		"post_response_author_id":       mock.JonSnow.ID,
		"post_response_author_name":     mock.JonSnow.Name,
		"post_response_author_email":    mock.JonSnow.Email,
		"post_response_author_role":     mock.JonSnow.Role.String(),
		"post_response_original_number": post.Response.Original.Number,
		"post_response_original_title":  post.Response.Original.Title,
		"post_response_original_slug":   post.Response.Original.Slug,
		"post_response_original_status_slug": post.Response.Original.StatusSlug,
		"post_response_original_url":    "http://domain.com/posts/1/add-support-for-typescript",
		"author_id":                     mock.JonSnow.ID,
		"author_name":                   mock.JonSnow.Name,
		"author_email":                  mock.JonSnow.Email,
		"author_role":                   mock.JonSnow.Role.String(),
		"tenant_id":                     mock.DemoTenant.ID,
		"tenant_name":                   mock.DemoTenant.Name,
		"tenant_subdomain":              mock.DemoTenant.Subdomain,
		"tenant_status":                 mock.DemoTenant.Status.String(),
		"tenant_url":                    "http://domain.com",
	})
}

// TestNotifyAboutStatusChange_CustomSlugBothEnumZero covers the regression
// where two tenant-defined statuses both backed by legacy enum 0 (the default)
// would have falsely short-circuited the webhook. Now that the guard compares
// slugs directly, the change fires.
func TestNotifyAboutStatusChange_CustomSlugBothEnumZero(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveStatusesForTenant) error {
		q.Result = []*entity.Status{
			{Slug: "triage", Label: "Triage", Kind: "open"},
			{Slug: "parked", Label: "Parked", Kind: "closed-declined"},
		}
		return nil
	})

	var triggerWebhooks *cmd.TriggerWebhooks
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooks) error {
		triggerWebhooks = c
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:         9,
		Number:     9,
		Title:      "Custom slug shift",
		Slug:       "custom-slug-shift",
		User:       mock.AryaStark,
		StatusSlug: "parked",
		StatusKind: "closed-declined",
	}

	task := tasks.NotifyAboutStatusChange(post, "triage")
	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(triggerWebhooks).IsNotNil()
	Expect(triggerWebhooks.Type).Equals(enum.WebhookChangeStatus)
	Expect(triggerWebhooks.Props).ContainsProps(webhook.Props{
		"post_old_status_slug":  "triage",
		"post_old_status_label": "Triage",
		"post_status_label":     "Parked",
		"post_status_slug":      "parked",
		"post_status_kind":      "closed-declined",
	})
}
