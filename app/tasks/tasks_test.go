package tasks_test

import (
	"context"
	"html/template"
	"testing"
	"time"

	"github.com/getfider/fider/app/actions"
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

func TestSendSignUpEmailTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	worker := mock.NewWorker()
	task := tasks.SendSignUpEmail(&actions.CreateTenant{
		VerificationKey: "1234",
	}, "http://domain.com")

	err := worker.
		AsUser(mock.JonSnow).
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("signup_email")
	Expect(emailmock.MessageHistory[0].Tenant).IsNil()
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"logo": "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Fider")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Props: dto.Props{
			"link": template.HTML("<a href='http://domain.com/signup/verify?k=1234'>http://domain.com/signup/verify?k=1234</a>"),
		},
	})
}

func TestSendSignInEmailTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	worker := mock.NewWorker()
	task := tasks.SendSignInEmail("jon@got.com", "9876")

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("signin_email")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"logo": "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(mock.DemoTenant.Name)
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Address: "jon@got.com",
		Props: dto.Props{
			"tenantName": mock.DemoTenant.Name,
			"link":       template.HTML("<a href='http://domain.com/signin/verify?k=9876'>http://domain.com/signin/verify?k=9876</a>"),
		},
	})
}

func TestSendChangeEmailConfirmationTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	worker := mock.NewWorker()
	task := tasks.SendChangeEmailConfirmation(&actions.ChangeUserEmail{
		Email:           "newemail@domain.com",
		VerificationKey: "13579",
		Requestor:       mock.JonSnow,
	})

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("change_emailaddress_email")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"logo": "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(mock.DemoTenant.Name)
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Jon Snow",
		Address: "newemail@domain.com",
		Props: dto.Props{
			"name":     "Jon Snow",
			"oldEmail": "jon.snow@got.com",
			"newEmail": "newemail@domain.com",
			"link":     template.HTML("<a href='http://domain.com/change-email/verify?k=13579'>http://domain.com/change-email/verify?k=13579</a>"),
		},
	})
}

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
		"title":      "Add support for TypeScript",
		"postLink":   template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>"),
		"tenantName": "Demonstration",
		"userName":   "Jon Snow",
		"content":    template.HTML("<p>TypeScript is great, please add support for it</p>"),
		"view":       template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>"),
		"change":     template.HTML("<a href='http://domain.com/settings'>change your notification preferences</a>"),
		"logo":       "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
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
}

func TestNotifyAboutNewCommentTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
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
	}
	task := tasks.NotifyAboutNewComment(post, "I agree")

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
		"title":       "Add support for TypeScript",
		"postLink":    template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>"),
		"tenantName":  "Demonstration",
		"userName":    "Arya Stark",
		"content":     template.HTML("<p>I agree</p>"),
		"view":        template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>"),
		"change":      template.HTML("<a href='http://domain.com/settings'>change your notification preferences</a>"),
		"unsubscribe": template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>"),
		"logo":        "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Arya Stark")
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
}

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

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "Add support for TypeScript",
		Slug:        "add-support-for-typescript",
		Description: "TypeScript is great, please add support for it",
		Status:      enum.PostPlanned,
		Response: &entity.PostResponse{
			RespondedAt: time.Now(),
			Text:        "Planned for next release.",
			User:        mock.JonSnow,
		},
	}

	task := tasks.NotifyAboutStatusChange(post, enum.PostOpen)

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
		"postLink":    template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>"),
		"tenantName":  "Demonstration",
		"content":     template.HTML("<p>Planned for next release.</p>"),
		"duplicate":   template.HTML(""),
		"status":      "planned",
		"view":        template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>"),
		"change":      template.HTML("<a href='http://domain.com/settings'>change your notification preferences</a>"),
		"unsubscribe": template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>"),
		"logo":        "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/1/add-support-for-typescript")
	Expect(addNewNotification.Title).Equals("**Jon Snow** changed status of **Add support for TypeScript** to **planned**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)
}

func TestNotifyAboutDeletePostTask(t *testing.T) {
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

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "Add support for TypeScript",
		Slug:        "add-support-for-typescript",
		Description: "TypeScript is great, please add support for it",
		Status:      enum.PostDeleted,
		Response: &entity.PostResponse{
			RespondedAt: time.Now(),
			Text:        "Invalid post!",
			User:        mock.JonSnow,
		},
	}

	task := tasks.NotifyAboutDeletedPost(post)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("delete_post")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":      "Add support for TypeScript",
		"tenantName": "Demonstration",
		"content":    template.HTML("<p>Invalid post!</p>"),
		"change":     template.HTML("<a href='http://domain.com/settings'>change your notification preferences</a>"),
		"logo":       "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("")
	Expect(addNewNotification.Title).Equals("**Jon Snow** deleted **Add support for TypeScript**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)
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

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:     2,
		Number: 2,
		Title:  "I need TypeScript",
		Slug:   "i-need-typescript",
		Status: enum.PostDuplicate,
		Response: &entity.PostResponse{
			RespondedAt: time.Now(),
			User:        mock.JonSnow,
			Original: &entity.OriginalPost{
				Number: 1,
				Title:  "Add support for TypeScript",
				Slug:   "add-support-for-typescript",
				Status: enum.PostOpen,
			},
		},
	}

	task := tasks.NotifyAboutStatusChange(post, enum.PostOpen)

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
		"postLink":    template.HTML("<a href='http://domain.com/posts/2/i-need-typescript'>#2</a>"),
		"tenantName":  "Demonstration",
		"content":     template.HTML(""),
		"duplicate":   template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>Add support for TypeScript</a>"),
		"status":      "duplicate",
		"view":        template.HTML("<a href='http://domain.com/posts/2/i-need-typescript'>view it on your browser</a>"),
		"change":      template.HTML("<a href='http://domain.com/settings'>change your notification preferences</a>"),
		"unsubscribe": template.HTML("<a href='http://domain.com/posts/2/i-need-typescript'>unsubscribe from it</a>"),
		"logo":        "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/2/i-need-typescript")
	Expect(addNewNotification.Title).Equals("**Jon Snow** changed status of **I need TypeScript** to **duplicate**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)
}

func TestSendInvites(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	savedKeys := make([]*cmd.SaveVerificationKey, 0)
	bus.AddHandler(func(ctx context.Context, c *cmd.SaveVerificationKey) error {
		savedKeys = append(savedKeys, c)
		return nil
	})

	worker := mock.NewWorker()
	task := tasks.SendInvites("My Subject", "Click here: %invite%", []*actions.UserInvitation{
		{Email: "user1@domain.com", VerificationKey: "1234"},
		{Email: "user2@domain.com", VerificationKey: "5678"},
	})

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("invite_email")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"subject": "My Subject",
		"logo":    "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
	Expect(emailmock.MessageHistory[0].To).HasLen(2)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Address: "user1@domain.com",
		Props: dto.Props{
			"message": template.HTML(`<p>Click here: <a href="http://domain.com/invite/verify?k=1234">http://domain.com/invite/verify?k=1234</a></p>`),
		},
	})
	Expect(emailmock.MessageHistory[0].To[1]).Equals(dto.Recipient{
		Address: "user2@domain.com",
		Props: dto.Props{
			"message": template.HTML(`<p>Click here: <a href="http://domain.com/invite/verify?k=5678">http://domain.com/invite/verify?k=5678</a></p>`),
		},
	})

	Expect(savedKeys).HasLen(2)
	Expect(savedKeys[0].Key).Equals("1234")
	Expect(savedKeys[0].Request.GetEmail()).Equals("user1@domain.com")
	Expect(savedKeys[1].Key).Equals("5678")
	Expect(savedKeys[1].Request.GetEmail()).Equals("user2@domain.com")
}
