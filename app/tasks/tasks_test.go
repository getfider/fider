package tasks_test

import (
	"html/template"
	"testing"
	"time"

	"github.com/getfider/fider/app/pkg/email"

	"github.com/getfider/fider/app/pkg/email/noop"
	"github.com/getfider/fider/app/pkg/mock"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/tasks"
)

func TestSendSignUpEmailTask(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	emailer := services.Emailer.(*noop.Sender)
	task := tasks.SendSignUpEmail(&models.CreateTenant{
		VerificationKey: "1234",
	}, "http://domain.com")

	err := worker.
		AsUser(mock.JonSnow).
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("signup_email")
	Expect(emailer.Requests[0].Context.Tenant()).IsNil()
	Expect(emailer.Requests[0].Params).Equals(email.Params{})
	Expect(emailer.Requests[0].From).Equals("Fider")
	Expect(emailer.Requests[0].To).HasLen(1)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Params: email.Params{
			"logo": "https://getfider.com/images/logo-100x100.png",
			"link": template.HTML("<a href='http://domain.com/signup/verify?k=1234'>http://domain.com/signup/verify?k=1234</a>"),
		},
	})
}

func TestSendSignInEmailTask(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	emailer := services.Emailer.(*noop.Sender)
	task := tasks.SendSignInEmail(&models.SignInByEmail{
		VerificationKey: "9876",
	})

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("signin_email")
	Expect(emailer.Requests[0].Context.Tenant()).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{})
	Expect(emailer.Requests[0].From).Equals(mock.DemoTenant.Name)
	Expect(emailer.Requests[0].To).HasLen(1)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Params: email.Params{
			"tenantName": mock.DemoTenant.Name,
			"link":       template.HTML("<a href='http://domain.com/signin/verify?k=9876'>http://domain.com/signin/verify?k=9876</a>"),
		},
	})
}

func TestSendChangeEmailConfirmationTask(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	emailer := services.Emailer.(*noop.Sender)
	task := tasks.SendChangeEmailConfirmation(&models.ChangeUserEmail{
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
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("change_emailaddress_email")
	Expect(emailer.Requests[0].Context.Tenant()).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{})
	Expect(emailer.Requests[0].From).Equals(mock.DemoTenant.Name)
	Expect(emailer.Requests[0].To).HasLen(1)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Name:    "Jon Snow",
		Address: "newemail@domain.com",
		Params: email.Params{
			"name":     "Jon Snow",
			"oldEmail": "jon.snow@got.com",
			"newEmail": "newemail@domain.com",
			"link":     template.HTML("<a href='http://domain.com/change-email/verify?k=13579'>http://domain.com/change-email/verify?k=13579</a>"),
		},
	})
}

func TestNotifyAboutNewPostTask(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("Add support for TypeScript", "TypeScript is great, please add support for it")

	services.Posts.AddSubscriber(post, mock.AryaStark)
	emailer := services.Emailer.(*noop.Sender)
	task := tasks.NotifyAboutNewPost(post)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("new_post")
	Expect(emailer.Requests[0].Context.Tenant()).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{
		"title":      "Add support for TypeScript",
		"postLink":   template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>"),
		"tenantName": "Demonstration",
		"userName":   "Jon Snow",
		"content":    template.HTML("<p>TypeScript is great, please add support for it</p>"),
		"view":       template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>View it on your browser</a>"),
		"change":     template.HTML("<a href='http://domain.com/settings'>change your notification settings</a>"),
	})
	Expect(emailer.Requests[0].From).Equals("Jon Snow")
	Expect(emailer.Requests[0].To).HasLen(1)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Params:  email.Params{},
	})

	services.SetCurrentUser(mock.AryaStark)
	notifications, err := services.Notifications.GetActiveNotifications()
	Expect(err).IsNil()
	Expect(notifications).HasLen(1)
	Expect(notifications[0].ID).Equals(1)
	Expect(notifications[0].CreatedAt).TemporarilySimilar(time.Now(), 5*time.Second)
	Expect(notifications[0].Link).Equals("/posts/1/add-support-for-typescript")
	Expect(notifications[0].Read).IsFalse()
	Expect(notifications[0].Title).Equals("New post: **Add support for TypeScript**")
}

func TestNotifyAboutNewCommentTask(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("Add support for TypeScript", "TypeScript is great, please add support for it")
	comment := &models.NewComment{
		Number:  post.Number,
		Content: "I agree",
	}

	emailer := services.Emailer.(*noop.Sender)
	task := tasks.NotifyAboutNewComment(post, comment)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("new_comment")
	Expect(emailer.Requests[0].Context.Tenant()).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{
		"title":       "Add support for TypeScript",
		"postLink":    template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>"),
		"tenantName":  "Demonstration",
		"userName":    "Arya Stark",
		"content":     template.HTML("<p>I agree</p>"),
		"view":        template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>View it on your browser</a>"),
		"change":      template.HTML("<a href='http://domain.com/settings'>change your notification settings</a>"),
		"unsubscribe": template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>"),
	})
	Expect(emailer.Requests[0].From).Equals("Arya Stark")
	Expect(emailer.Requests[0].To).HasLen(1)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Name:    "Jon Snow",
		Address: "jon.snow@got.com",
		Params:  email.Params{},
	})

	services.SetCurrentUser(mock.JonSnow)
	notifications, err := services.Notifications.GetActiveNotifications()
	Expect(err).IsNil()
	Expect(notifications).HasLen(1)
	Expect(notifications[0].ID).Equals(1)
	Expect(notifications[0].CreatedAt).TemporarilySimilar(time.Now(), 5*time.Second)
	Expect(notifications[0].Link).Equals("/posts/1/add-support-for-typescript")
	Expect(notifications[0].Read).IsFalse()
	Expect(notifications[0].Title).Equals("**Arya Stark** left a comment on **Add support for TypeScript**")
}

func TestNotifyAboutStatusChangeTask(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("Add support for TypeScript", "TypeScript is great, please add support for it")
	services.Posts.SetResponse(post, "Planned for next release.", models.PostPlanned)

	emailer := services.Emailer.(*noop.Sender)
	task := tasks.NotifyAboutStatusChange(post, models.PostOpen)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("change_status")
	Expect(emailer.Requests[0].Context.Tenant()).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{
		"title":       "Add support for TypeScript",
		"postLink":    template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>"),
		"tenantName":  "Demonstration",
		"content":     template.HTML("<p>Planned for next release.</p>"),
		"duplicate":   template.HTML(""),
		"status":      "planned",
		"view":        template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>View it on your browser</a>"),
		"change":      template.HTML("<a href='http://domain.com/settings'>change your notification settings</a>"),
		"unsubscribe": template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>"),
	})
	Expect(emailer.Requests[0].From).Equals("Jon Snow")
	Expect(emailer.Requests[0].To).HasLen(1)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Params:  email.Params{},
	})

	services.SetCurrentUser(mock.AryaStark)
	notifications, err := services.Notifications.GetActiveNotifications()
	Expect(err).IsNil()
	Expect(notifications).HasLen(1)
	Expect(notifications[0].ID).Equals(1)
	Expect(notifications[0].CreatedAt).TemporarilySimilar(time.Now(), 5*time.Second)
	Expect(notifications[0].Link).Equals("/posts/1/add-support-for-typescript")
	Expect(notifications[0].Read).IsFalse()
	Expect(notifications[0].Title).Equals("**Jon Snow** changed status of **Add support for TypeScript** to **planned**")
}

func TestNotifyAboutDeletePostTask(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("Add support for TypeScript", "TypeScript is great, please add support for it")
	services.Posts.SetResponse(post, "Invalid post!", models.PostDeleted)

	emailer := services.Emailer.(*noop.Sender)
	task := tasks.NotifyAboutDeletedPost(post)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("delete_post")
	Expect(emailer.Requests[0].Context.Tenant()).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{
		"title":      "Add support for TypeScript",
		"tenantName": "Demonstration",
		"content":    template.HTML("<p>Invalid post!</p>"),
		"change":     template.HTML("<a href='http://domain.com/settings'>change your notification settings</a>"),
	})
	Expect(emailer.Requests[0].From).Equals("Jon Snow")
	Expect(emailer.Requests[0].To).HasLen(1)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Params:  email.Params{},
	})

	services.SetCurrentUser(mock.AryaStark)
	notifications, err := services.Notifications.GetActiveNotifications()
	Expect(err).IsNil()
	Expect(notifications).HasLen(1)
	Expect(notifications[0].ID).Equals(1)
	Expect(notifications[0].CreatedAt).TemporarilySimilar(time.Now(), 5*time.Second)
	Expect(notifications[0].Link).Equals("")
	Expect(notifications[0].Read).IsFalse()
	Expect(notifications[0].Title).Equals("**Jon Snow** deleted **Add support for TypeScript**")
}

func TestNotifyAboutStatusChangeTask_Duplicate(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	services.SetCurrentUser(mock.AryaStark)
	post1, _ := services.Posts.Add("Add support for TypeScript", "TypeScript is great, please add support for it")
	post2, _ := services.Posts.Add("I need TypeScript", "")
	services.Posts.MarkAsDuplicate(post2, post1)

	emailer := services.Emailer.(*noop.Sender)
	task := tasks.NotifyAboutStatusChange(post2, models.PostOpen)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("change_status")
	Expect(emailer.Requests[0].Context.Tenant()).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{
		"title":       "I need TypeScript",
		"postLink":    template.HTML("<a href='http://domain.com/posts/2/i-need-typescript'>#2</a>"),
		"tenantName":  "Demonstration",
		"content":     template.HTML(""),
		"duplicate":   template.HTML("<a href='http://domain.com/posts/1/add-support-for-typescript'>Add support for TypeScript</a>"),
		"status":      "duplicate",
		"view":        template.HTML("<a href='http://domain.com/posts/2/i-need-typescript'>View it on your browser</a>"),
		"change":      template.HTML("<a href='http://domain.com/settings'>change your notification settings</a>"),
		"unsubscribe": template.HTML("<a href='http://domain.com/posts/2/i-need-typescript'>unsubscribe from it</a>"),
	})
	Expect(emailer.Requests[0].From).Equals("Jon Snow")
	Expect(emailer.Requests[0].To).HasLen(1)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Params:  email.Params{},
	})

	services.SetCurrentUser(mock.AryaStark)
	notifications, err := services.Notifications.GetActiveNotifications()
	Expect(err).IsNil()
	Expect(notifications).HasLen(1)
	Expect(notifications[0].ID).Equals(1)
	Expect(notifications[0].CreatedAt).TemporarilySimilar(time.Now(), 5*time.Second)
	Expect(notifications[0].Link).Equals("/posts/2/i-need-typescript")
	Expect(notifications[0].Read).IsFalse()
	Expect(notifications[0].Title).Equals("**Jon Snow** changed status of **I need TypeScript** to **duplicate**")
}

func TestSendInvites(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	emailer := services.Emailer.(*noop.Sender)
	task := tasks.SendInvites("My Subject", "Click here: %invite%", []*models.UserInvitation{
		&models.UserInvitation{Email: "user1@domain.com", VerificationKey: "1234"},
		&models.UserInvitation{Email: "user2@domain.com", VerificationKey: "5678"},
	})

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("invite_email")
	Expect(emailer.Requests[0].Context.Tenant()).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{
		"subject": "My Subject",
	})
	Expect(emailer.Requests[0].From).Equals("Jon Snow")
	Expect(emailer.Requests[0].To).HasLen(2)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Address: "user1@domain.com",
		Params: email.Params{
			"message": template.HTML(`<p>Click here: <a href="http://domain.com/invite/verify?k=1234">http://domain.com/invite/verify?k=1234</a></p>`),
		},
	})
	Expect(emailer.Requests[0].To[1]).Equals(email.Recipient{
		Address: "user2@domain.com",
		Params: email.Params{
			"message": template.HTML(`<p>Click here: <a href="http://domain.com/invite/verify?k=5678">http://domain.com/invite/verify?k=5678</a></p>`),
		},
	})
}
