package tasks_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/mock"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/tasks"
	. "github.com/onsi/gomega"
)

func TestSendSignInEmailTask(t *testing.T) {
	RegisterTestingT(t)

	worker, _ := mock.NewWorker()
	task := tasks.SendSignInEmail(&models.SignInByEmail{}, "http://anywhere.com")
	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(task)
	Expect(err).To(BeNil())
}

func TestSendChangeEmailConfirmationTask(t *testing.T) {
	RegisterTestingT(t)

	worker, _ := mock.NewWorker()
	task := tasks.SendChangeEmailConfirmation(&models.ChangeUserEmail{}, "http://anywhere.com")
	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(task)
	Expect(err).To(BeNil())
}
