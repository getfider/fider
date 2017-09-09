package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestLoginByEmail_EmptyEmail(t *testing.T) {
	RegisterTestingT(t)

	action := actions.LoginByEmail{Model: &models.LoginByEmail{Email: ""}}
	result := action.Validate(services)
	ExpectFailed(result, "email")
}

func TestLoginByEmail_InvalidEmail(t *testing.T) {
	RegisterTestingT(t)

	action := actions.LoginByEmail{Model: &models.LoginByEmail{Email: "Hi :)"}}
	result := action.Validate(services)
	ExpectFailed(result, "email")
}

func TestLoginByEmail_ShouldHaveVerificationKey(t *testing.T) {
	RegisterTestingT(t)

	action := actions.LoginByEmail{}
	action.Initialize()
	action.Model.Email = "jon.snow@got.com"

	result := action.Validate(services)
	ExpectSuccess(result)
	Expect(action.Model.VerificationKey).NotTo(Equal(""))
}
