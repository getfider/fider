package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestSignInByEmail_EmptyEmail(t *testing.T) {
	RegisterTestingT(t)

	action := actions.SignInByEmail{Model: &models.SignInByEmail{Email: ""}}
	result := action.Validate(services)
	ExpectFailed(result, "email")
}

func TestSignInByEmail_InvalidEmail(t *testing.T) {
	RegisterTestingT(t)

	action := actions.SignInByEmail{Model: &models.SignInByEmail{Email: "Hi :)"}}
	result := action.Validate(services)
	ExpectFailed(result, "email")
}

func TestSignInByEmail_ShouldHaveVerificationKey(t *testing.T) {
	RegisterTestingT(t)

	action := actions.SignInByEmail{}
	action.Initialize()
	action.Model.Email = "jon.snow@got.com"

	result := action.Validate(services)
	ExpectSuccess(result)
	Expect(action.Model.VerificationKey).NotTo(Equal(""))
}
