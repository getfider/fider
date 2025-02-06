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
		"logo": "https://fider.io/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(dto.Recipient{
		Name: "Fider",
	})
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Props: dto.Props{
			"link": "<a href='http://domain.com/signup/verify?k=1234'>http://domain.com/signup/verify?k=1234</a>",
		},
	})
}
