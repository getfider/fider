package tasks_test

import (
	"context"
	"html/template"
	"testing"

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

func TestNotifyAboutNewPostTask(t *testing.T) {
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
	}
	task := tasks.NotifyAboutNewPost(post)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("new_post")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":    "Add support for TypeScript",
		"postLink": "<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>",
		"siteName": "Demonstration",
		"userName": "Jon Snow",
		"content":  template.HTML("<p>TypeScript is great, please add support for it</p>"),
		"view":     "<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>",
		"change":   "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"logo":     "https://fider.io/images/logo-100x100.png",
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
	Expect(addNewNotification.Title).Equals("New post: **Add support for TypeScript**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)

	Expect(triggerWebhooks).IsNotNil()
	Expect(triggerWebhooks.Type).Equals(enum.WebhookNewPost)
	Expect(triggerWebhooks.Props).ContainsProps(webhook.Props{
		"post_id":          post.ID,
		"post_number":      post.Number,
		"post_title":       post.Title,
		"post_slug":        post.Slug,
		"post_description": post.Description,
		"post_url":         "http://domain.com/posts/1/add-support-for-typescript",
		"author_id":        mock.JonSnow.ID,
		"author_name":      mock.JonSnow.Name,
		"author_email":     mock.JonSnow.Email,
		"author_role":      mock.JonSnow.Role.String(),
		"tenant_id":        mock.DemoTenant.ID,
		"tenant_name":      mock.DemoTenant.Name,
		"tenant_subdomain": mock.DemoTenant.Subdomain,
		"tenant_status":    mock.DemoTenant.Status.String(),
		"tenant_url":       "http://domain.com",
	})
}

func TestNotifyAboutNewPostTask_WithMention(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	addNotificationLogs := make([]*cmd.AddMentionNotification, 0)
	bus.AddHandler(func(ctx context.Context, c *cmd.AddMentionNotification) error {
		addNotificationLogs = append(addNotificationLogs, c)
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		if q.Event.UserSettingsKeyName == "event_notification_mention" {
			q.Result = []*entity.User{
				mock.JonSnow,
			}
		} else {
			q.Result = []*entity.User{}
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetMentionNotifications) error {
		q.Result = []*entity.MentionNotification{}
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
		Description: "TypeScript is great, please add support for it @[Jon Snow]",
	}
	task := tasks.NotifyAboutNewPost(post)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(2)

	// Check standard notification email
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("new_post")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props["content"]).Equals(template.HTML("<p>TypeScript is great, please add support for it @Jon Snow</p>"))

	// Check mention notification email
	Expect(emailmock.MessageHistory[1].TemplateName).Equals("new_post")
	Expect(emailmock.MessageHistory[1].Props["messageLocaleString"]).Equals("email.new_mention.text")
	Expect(emailmock.MessageHistory[1].To).HasLen(1)
	Expect(emailmock.MessageHistory[1].To[0].Name).Equals("Jon Snow")

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Title).Equals("**Arya Stark** mentioned you in **Add support for TypeScript**")
	Expect(addNewNotification.User).Equals(mock.JonSnow)

	Expect(addNotificationLogs).HasLen(2)
	Expect(addNotificationLogs[0].UserID).Equals(mock.JonSnow.ID)
	Expect(addNotificationLogs[0].PostID).Equals(post.ID)
	Expect(addNotificationLogs[1].UserID).Equals(mock.JonSnow.ID)
	Expect(addNotificationLogs[1].PostID).Equals(post.ID)

	Expect(triggerWebhooks).IsNotNil()
	Expect(triggerWebhooks.Type).Equals(enum.WebhookNewPost)
}
