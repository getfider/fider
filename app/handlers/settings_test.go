package handlers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app"
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

func TestVerifyChangeEmailKeyHandler_Success(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	request := &models.ChangeUserEmail{
		Requestor: mock.JonSnow,
		Email:     "jon.stark@got.com",
	}

	services.Tenants.SaveVerificationKey("th3-s3cr3t", 24*time.Hour, request)
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithURL("/change-email/verify?k=th3-s3cr3t").
		Execute(handlers.VerifyChangeEmailKey())

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	user, err := services.Users.GetByEmail(mock.DemoTenant.ID, "jon.stark@got.com")
	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(mock.JonSnow.ID))
	Expect(user.Name).To(Equal(mock.JonSnow.Name))
	Expect(user.Email).To(Equal(mock.JonSnow.Email))

	result, err := services.Tenants.FindVerificationByKey(models.EmailVerificationKindChangeEmail, "th3-s3cr3t")
	Expect(err).To(BeNil())
	Expect(result.VerifiedOn).NotTo(BeNil())
}

func TestVerifyChangeEmailKeyHandler_DifferentUser(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	request := &models.ChangeUserEmail{
		Requestor: mock.JonSnow,
		Email:     "jon.stark@got.com",
	}
	services.Tenants.SaveVerificationKey("th3-s3cr3t", 24*time.Hour, request)
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithURL("/change-email/verify?k=th3-s3cr3t").
		Execute(handlers.VerifyChangeEmailKey())

	Expect(code).To(Equal(http.StatusTemporaryRedirect))

	_, err := services.Users.GetByEmail(mock.DemoTenant.ID, "jon.snow@got.com")
	Expect(err).To(BeNil())
	_, err = services.Users.GetByEmail(mock.DemoTenant.ID, "arya.stark@got.com")
	Expect(err).To(BeNil())

	_, err = services.Users.GetByEmail(mock.DemoTenant.ID, "jon.stark@got.com")
	Expect(err).To(Equal(app.ErrNotFound))
}
