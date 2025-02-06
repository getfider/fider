package tasks_test

import (
	"testing"

	"github.com/Spicy-Bush/fider-tarkov-community/app/actions"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	. "github.com/Spicy-Bush/fider-tarkov-community/app/pkg/assert"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/mock"
	"github.com/Spicy-Bush/fider-tarkov-community/app/services/email/emailmock"
	"github.com/Spicy-Bush/fider-tarkov-community/app/tasks"
)

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
		"logo": "https://fider.io/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(dto.Recipient{
		Name: mock.DemoTenant.Name,
	})
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Jon Snow",
		Address: "newemail@domain.com",
		Props: dto.Props{
			"name":     "Jon Snow",
			"oldEmail": "jon.snow@got.com",
			"newEmail": "newemail@domain.com",
			"link":     "<a href='http://domain.com/change-email/verify?k=13579'>http://domain.com/change-email/verify?k=13579</a>",
		},
	})
}
