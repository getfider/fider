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
	Expect(emailer.Requests[0].Tenant).IsNil()
	Expect(emailer.Requests[0].Params).Equals(email.Params{})
	Expect(emailer.Requests[0].From).Equals("Fider")
	Expect(emailer.Requests[0].To).HasLen(1)
	Expect(emailer.Requests[0].To[0]).Equals(email.Recipient{
		Params: email.Params{
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
	Expect(emailer.Requests[0].Tenant).Equals(mock.DemoTenant)
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
	Expect(emailer.Requests[0].Tenant).Equals(mock.DemoTenant)
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

func TestNotifyAboutNewIdeaTask(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("Add support for TypeScript", "TypeScript is great, please add support for it")

	services.Ideas.AddSubscriber(idea, mock.AryaStark)
	emailer := services.Emailer.(*noop.Sender)
	task := tasks.NotifyAboutNewIdea(idea)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("new_idea")
	Expect(emailer.Requests[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{
		"title":   "[Demonstration] Add support for TypeScript",
		"content": template.HTML("<p>TypeScript is great, please add support for it</p>"),
		"view":    template.HTML("<a href='http://domain.com/ideas/1/add-support-for-typescript'>View it on your browser</a>"),
		"change":  template.HTML("<a href='http://domain.com/settings'>change your notification settings</a>"),
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
	Expect(notifications[0].CreatedOn).TemporarilySimilar(time.Now(), 5*time.Second)
	Expect(notifications[0].Link).Equals("/ideas/1/add-support-for-typescript")
	Expect(notifications[0].Read).IsFalse()
	Expect(notifications[0].Title).Equals("New idea: **Add support for TypeScript**")
}

func TestNotifyAboutNewCommentTask(t *testing.T) {
	RegisterT(t)

	worker, services := mock.NewWorker()
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("Add support for TypeScript", "TypeScript is great, please add support for it")
	comment := &models.NewComment{
		Number:  idea.Number,
		Content: "I agree",
	}

	emailer := services.Emailer.(*noop.Sender)
	task := tasks.NotifyAboutNewComment(idea, comment)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailer.Requests).HasLen(1)
	Expect(emailer.Requests[0].TemplateName).Equals("new_comment")
	Expect(emailer.Requests[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailer.Requests[0].Params).Equals(email.Params{
		"title":       "[Demonstration] Add support for TypeScript",
		"content":     template.HTML("<p>I agree</p>"),
		"view":        template.HTML("<a href='http://domain.com/ideas/1/add-support-for-typescript'>View it on your browser</a>"),
		"change":      template.HTML("<a href='http://domain.com/settings'>change your notification settings</a>"),
		"unsubscribe": template.HTML("<a href='http://domain.com/ideas/1/add-support-for-typescript'>unsubscribe from it</a>"),
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
	Expect(notifications[0].CreatedOn).TemporarilySimilar(time.Now(), 5*time.Second)
	Expect(notifications[0].Link).Equals("/ideas/1/add-support-for-typescript")
	Expect(notifications[0].Read).IsFalse()
	Expect(notifications[0].Title).Equals("**Arya Stark** left a comment on **Add support for TypeScript**")
}
