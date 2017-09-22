package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/pkg/mock"
	. "github.com/onsi/gomega"
)

func TestUpdateUserSettingsHandler_EmptyInput(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		AsUser(mock.JonSnow).
		ExecutePost(handlers.UpdateUserSettings(), `{ }`)

	Expect(code).To(Equal(http.StatusBadRequest))
}

func TestUpdateUserSettingsHandler_ValidName(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.UpdateUserSettings(), `{ "name": "Jon Stark" }`)

	user, _ := services.Users.GetByEmail(mock.DemoTenant.ID, "jon.snow@got.com")

	Expect(code).To(Equal(http.StatusOK))
	Expect(user.Name).To(Equal("Jon Stark"))
}
