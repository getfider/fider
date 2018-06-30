package tasks_test

import (
	"html/template"
	"testing"

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
	}, "http://anywhere.com")

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
			"link": template.HTML("<a href='http://anywhere.com/signup/verify?k=1234'>http://anywhere.com/signup/verify?k=1234</a>"),
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
		WithBaseURL("http://anywhere.com").
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
			"link":       template.HTML("<a href='http://anywhere.com/signin/verify?k=9876'>http://anywhere.com/signin/verify?k=9876</a>"),
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
		WithBaseURL("http://anywhere.com").
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
			"link":     template.HTML("<a href='http://anywhere.com/change-email/verify?k=13579'>http://anywhere.com/change-email/verify?k=13579</a>"),
		},
	})
}
