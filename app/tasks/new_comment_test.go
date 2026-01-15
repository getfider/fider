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

func TestNotifyAboutNewCommentTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	addNotificationLogs := make([]*cmd.AddMentionNotification, 0)
	bus.AddListener(func(ctx context.Context, c *cmd.AddMentionNotification) error {
		addNotificationLogs = append(addNotificationLogs, c)
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetMentionNotifications) error {
		q.Result = []*entity.MentionNotification{}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		if q.Event.UserSettingsKeyName == "event_notification_new_comment" {
			q.Result = []*entity.User{
				mock.JonSnow,
			}
		} else {
			q.Result = []*entity.User{}
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
		User:        mock.JonSnow,
	}
	task := tasks.NotifyAboutNewComment(&entity.Comment{Content: "I agree"}, post)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("new_comment")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":               "Add support for TypeScript",
		"messageLocaleString": "email.new_comment.text",
		"postLink":            "<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>",
		"siteName":            "Demonstration",
		"userName":            "Arya Stark",
		"content":             template.HTML("<p>I agree</p>"),
		"view":                "<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>",
		"change":              "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"unsubscribe":         "<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>",
		"logo":                "https://login.fider.io/static/assets/logo.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(dto.Recipient{
		Name: "Arya Stark",
	})
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Jon Snow",
		Address: "jon.snow@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/1/add-support-for-typescript")
	Expect(addNewNotification.Title).Equals("**Arya Stark** left a comment on **Add support for TypeScript**")
	Expect(addNewNotification.User).Equals(mock.JonSnow)

	Expect(triggerWebhooks).IsNotNil()
	Expect(triggerWebhooks.Type).Equals(enum.WebhookNewComment)
	Expect(triggerWebhooks.Props).ContainsProps(webhook.Props{
		"comment":           "I agree",
		"post_id":           post.ID,
		"post_number":       post.Number,
		"post_title":        post.Title,
		"post_slug":         post.Slug,
		"post_description":  post.Description,
		"post_url":          "http://domain.com/posts/1/add-support-for-typescript",
		"post_author_id":    mock.JonSnow.ID,
		"post_author_name":  mock.JonSnow.Name,
		"post_author_email": mock.JonSnow.Email,
		"post_author_role":  mock.JonSnow.Role.String(),
		"author_id":         mock.AryaStark.ID,
		"author_name":       mock.AryaStark.Name,
		"author_email":      mock.AryaStark.Email,
		"author_role":       mock.AryaStark.Role.String(),
		"tenant_id":         mock.DemoTenant.ID,
		"tenant_name":       mock.DemoTenant.Name,
		"tenant_subdomain":  mock.DemoTenant.Subdomain,
		"tenant_status":     mock.DemoTenant.Status.String(),
		"tenant_url":        "http://domain.com",
	})

	Expect(addNotificationLogs).HasLen(0)
}

func TestNotifyAboutNewCommentTask_WithMention(t *testing.T) {
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
		Description: "TypeScript is great, please add support for it",
		User:        mock.JonSnow,
	}
	task := tasks.NotifyAboutNewComment(&entity.Comment{Content: "I agree with @[Jon Snow]"}, post)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("new_comment")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":               "Add support for TypeScript",
		"messageLocaleString": "email.new_mention.text",
		"postLink":            "<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>",
		"siteName":            "Demonstration",
		"userName":            "Arya Stark",
		"content":             template.HTML("<p>I agree with @Jon Snow</p>"),
		"view":                "<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>",
		"change":              "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"unsubscribe":         "<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>",
		"logo":                "https://login.fider.io/static/assets/logo.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(dto.Recipient{
		Name: "Arya Stark",
	})
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Jon Snow",
		Address: "jon.snow@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/1/add-support-for-typescript")
	Expect(addNewNotification.Title).Equals("**Arya Stark** mentioned you in **Add support for TypeScript**")
	Expect(addNewNotification.User).Equals(mock.JonSnow)

	Expect(triggerWebhooks).IsNotNil()
	Expect(triggerWebhooks.Type).Equals(enum.WebhookNewComment)
	Expect(triggerWebhooks.Props).ContainsProps(webhook.Props{
		"comment":           "I agree with @Jon Snow",
		"post_id":           post.ID,
		"post_number":       post.Number,
		"post_title":        post.Title,
		"post_slug":         post.Slug,
		"post_description":  post.Description,
		"post_url":          "http://domain.com/posts/1/add-support-for-typescript",
		"post_author_id":    mock.JonSnow.ID,
		"post_author_name":  mock.JonSnow.Name,
		"post_author_email": mock.JonSnow.Email,
		"post_author_role":  mock.JonSnow.Role.String(),
		"author_id":         mock.AryaStark.ID,
		"author_name":       mock.AryaStark.Name,
		"author_email":      mock.AryaStark.Email,
		"author_role":       mock.AryaStark.Role.String(),
		"tenant_id":         mock.DemoTenant.ID,
		"tenant_name":       mock.DemoTenant.Name,
		"tenant_subdomain":  mock.DemoTenant.Subdomain,
		"tenant_status":     mock.DemoTenant.Status.String(),
		"tenant_url":        "http://domain.com",
	})

	Expect(addNotificationLogs).HasLen(2)
	Expect(addNotificationLogs[0].UserID).Equals(mock.JonSnow.ID)
	Expect(addNotificationLogs[0].CommentID).Equals(0)

	Expect(addNotificationLogs[1].UserID).Equals(mock.JonSnow.ID)
	Expect(addNotificationLogs[1].CommentID).Equals(0)

}

func TestNotifyAboutUpdatedComment(t *testing.T) {
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

	bus.AddHandler(func(ctx context.Context, q *query.GetMentionNotifications) error {
		q.Result = []*entity.MentionNotification{}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{
			mock.JonSnow,
		}
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "Add support for TypeScript",
		Slug:        "add-support-for-typescript",
		Description: "TypeScript is great, please add support for it",
		User:        mock.JonSnow,
	}

	comment := &entity.Comment{
		ID:        1,
		Content:   "I agree with @[Jon Snow] but not @[Arya Stark]",
		CreatedAt: time.Now(),
		User:      mock.AryaStark,
	}

	task := tasks.NotifyAboutUpdatedComment(post, comment)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("new_comment")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":               "Add support for TypeScript",
		"messageLocaleString": "email.new_mention.text",
		"postLink":            "<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>",
		"siteName":            "Demonstration",
		"userName":            "Arya Stark",
		"content":             template.HTML("<p>I agree with @Jon Snow but not @Arya Stark</p>"),
		"view":                "<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>",
		"change":              "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"unsubscribe":         "<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>",
		"logo":                "https://login.fider.io/static/assets/logo.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(dto.Recipient{
		Name: "Arya Stark",
	})
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Jon Snow",
		Address: "jon.snow@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/1/add-support-for-typescript")
	Expect(addNewNotification.Title).Equals("**Arya Stark** mentioned you in **Add support for TypeScript**")
	Expect(addNewNotification.User).Equals(mock.JonSnow)

	Expect(addNotificationLogs).HasLen(1)
	Expect(addNotificationLogs[0].UserID).Equals(mock.JonSnow.ID)
	Expect(addNotificationLogs[0].CommentID).Equals(1)

}

func TestNotifyAboutUpdatedComment_UserAlreadyMentioned(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetMentionNotifications) error {
		q.Result = []*entity.MentionNotification{
			{
				UserID:    mock.JonSnow.ID,
				CommentID: 1,
				TenantID:  mock.DemoTenant.ID,
			},
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{
			mock.JonSnow,
		}
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "Add support for TypeScript",
		Slug:        "add-support-for-typescript",
		Description: "TypeScript is great, please add support for it",
		User:        mock.JonSnow,
	}

	comment := &entity.Comment{
		ID:        1,
		Content:   "I agree with @[Jon Snow] but not @[Arya Stark]",
		CreatedAt: time.Now(),
		User:      mock.AryaStark,
	}

	task := tasks.NotifyAboutUpdatedComment(post, comment)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(0)

	Expect(addNewNotification).IsNil()
}
