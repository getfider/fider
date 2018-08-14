package actions_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestSignInByEmail_EmptyEmail(t *testing.T) {
	RegisterT(t)

	action := actions.SignInByEmail{Model: &models.SignInByEmail{Email: " "}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "email")
}

func TestSignInByEmail_InvalidEmail(t *testing.T) {
	RegisterT(t)

	action := actions.SignInByEmail{Model: &models.SignInByEmail{Email: "Hi :)"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "email")
}

func TestSignInByEmail_ShouldHaveVerificationKey(t *testing.T) {
	RegisterT(t)

	action := actions.SignInByEmail{}
	action.Initialize()
	action.Model.Email = "jon.snow@got.com"

	result := action.Validate(nil, services)
	ExpectSuccess(result)
	Expect(action.Model.VerificationKey).IsNotEmpty()
}

func TestCompleteProfile_EmptyNameAndKey(t *testing.T) {
	RegisterT(t)

	action := actions.CompleteProfile{Model: &models.CompleteProfile{}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "name", "key")
}

func TestCompleteProfile_LongName(t *testing.T) {
	RegisterT(t)

	action := actions.CompleteProfile{Model: &models.CompleteProfile{
		Name: "123456789012345678901234567890123456789012345678901", // 51 chars
	}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "name", "key")
}

func TestCompleteProfile_UnknownKey(t *testing.T) {
	RegisterT(t)
	action := actions.CompleteProfile{Model: &models.CompleteProfile{Name: "Jon Snow", Key: "1234567890"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "key")
}

func TestCompleteProfile_ValidKey(t *testing.T) {
	RegisterT(t)

	e := &models.SignInByEmail{Email: "jon.snow@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 30*time.Minute, e)
	action := actions.CompleteProfile{Model: &models.CompleteProfile{Name: "Jon Snow", Key: "1234567890"}}
	result := action.Validate(nil, services)

	ExpectSuccess(result)
	Expect(action.Model.Email).Equals("jon.snow@got.com")
}

func TestCompleteProfile_UserInvitation_ValidKey(t *testing.T) {
	RegisterT(t)

	e := &models.UserInvitation{Email: "jon.snow@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 30*time.Minute, e)
	action := actions.CompleteProfile{Model: &models.CompleteProfile{Name: "Jon Snow", Key: "1234567890"}}
	result := action.Validate(nil, services)

	ExpectSuccess(result)
	Expect(action.Model.Email).Equals("jon.snow@got.com")
}
