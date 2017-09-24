package actions_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestSignInByEmail_EmptyEmail(t *testing.T) {
	RegisterTestingT(t)

	action := actions.SignInByEmail{Model: &models.SignInByEmail{Email: " "}}
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

func TestCompleteProfile_EmptyNameAndKey(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CompleteProfile{Model: &models.CompleteProfile{}}
	result := action.Validate(services)
	ExpectFailed(result, "name")
	ExpectFailed(result, "key")
}

func TestCompleteProfile_LongName(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CompleteProfile{Model: &models.CompleteProfile{
		Name: "123456789012345678901234567890123456789012345678901", // 51 chars
	}}
	result := action.Validate(services)
	ExpectFailed(result, "name")
	ExpectFailed(result, "key")
}

func TestCompleteProfile_UnknownKey(t *testing.T) {
	RegisterTestingT(t)
	action := actions.CompleteProfile{Model: &models.CompleteProfile{Name: "Jon Snow", Key: "1234567890"}}
	result := action.Validate(services)
	ExpectFailed(result, "key")
}

func TestCompleteProfile_ValidKey(t *testing.T) {
	RegisterTestingT(t)
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, "jon.snow@got.com", "")
	action := actions.CompleteProfile{Model: &models.CompleteProfile{Name: "Jon Snow", Key: "1234567890"}}
	result := action.Validate(services)
	ExpectSuccess(result)
	Expect(action.Model.Email).To(Equal("jon.snow@got.com"))
}
