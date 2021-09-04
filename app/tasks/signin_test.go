package tasks_test

import (
	"testing"

	"github.com/getfider/fider/app/models/dto"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/services/email/emailmock"
	"github.com/getfider/fider/app/tasks"
)

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
			"siteName": mock.DemoTenant.Name,
			"link":     "<a href='http://domain.com/signin/verify?k=9876'>http://domain.com/signin/verify?k=9876</a>",
		},
	})
}
