package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/models"

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

func TestChangeRoleHandler_Valid(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("user_id", mock.AryaStark.ID).
		ExecutePost(handlers.ChangeUserRole(), fmt.Sprintf(`{ "role": %d }`, models.RoleAdministrator))

	user, _ := services.Users.GetByID(mock.AryaStark.ID)

	Expect(code).To(Equal(http.StatusOK))
	Expect(user.Role).To(Equal(models.RoleAdministrator))
}

func TestChangeUserEmailHandler_Valid(t *testing.T) {
	RegisterTestingT(t)

	for _, email := range []string{
		"jon.another@got.com",
		"another.snow@got.com",
	} {
		server, _ := mock.NewServer()
		code, _ := server.
			OnTenant(mock.DemoTenant).
			AsUser(mock.JonSnow).
			ExecutePost(handlers.ChangeUserEmail(), fmt.Sprintf(`{ "email": "%s" }`, email))

		Expect(code).To(Equal(http.StatusOK))
	}
}

func TestChangeUserEmailHandler_Invalid(t *testing.T) {
	RegisterTestingT(t)

	for _, email := range []string{
		"",
		"jon.snow@got.com",
		"jon.snow",
		"arya.stark@got.com",
	} {
		server, _ := mock.NewServer()
		code, _ := server.
			OnTenant(mock.DemoTenant).
			AsUser(mock.JonSnow).
			ExecutePost(handlers.ChangeUserEmail(), fmt.Sprintf(`{ "email": "%s" }`, email))

		Expect(code).To(Equal(http.StatusBadRequest))
	}
}
