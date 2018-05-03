package tasks_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/mock"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/tasks"
)

func TestSendSignUpEmailTask(t *testing.T) {
	RegisterT(t)

	worker, _ := mock.NewWorker()
	task := tasks.SendSignUpEmail(&models.CreateTenant{}, "http://anywhere.com")
	err := worker.
		AsUser(mock.JonSnow).
		Execute(task)
	Expect(err).IsNil()
}

func TestSendSignInEmailTask(t *testing.T) {
	RegisterT(t)

	worker, _ := mock.NewWorker()
	task := tasks.SendSignInEmail(&models.SignInByEmail{})
	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(task)
	Expect(err).IsNil()
}

func TestSendChangeEmailConfirmationTask(t *testing.T) {
	RegisterT(t)

	worker, _ := mock.NewWorker()
	task := tasks.SendChangeEmailConfirmation(&models.ChangeUserEmail{
		Requestor: &models.User{
			Name:  "Some User",
			Email: "some@domain.com",
		},
	})
	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(task)
	Expect(err).IsNil()
}
