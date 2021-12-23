package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/actions"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestSignInByEmail_EmptyEmail(t *testing.T) {
	RegisterT(t)

	action := actions.SignInByEmail{Email: " "}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "email")
}

func TestSignInByEmail_InvalidEmail(t *testing.T) {
	RegisterT(t)

	action := actions.SignInByEmail{Email: "Hi :)"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "email")
}

func TestSignInByEmail_ShouldHaveVerificationKey(t *testing.T) {
	RegisterT(t)

	action := actions.NewSignInByEmail()
	action.Email = "jon.snow@got.com"

	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
	Expect(action.VerificationKey).IsNotEmpty()
}

func TestCompleteProfile_EmptyNameAndKey(t *testing.T) {
	RegisterT(t)

	action := actions.CompleteProfile{}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "name", "key")
}

func TestCompleteProfile_LongName(t *testing.T) {
	RegisterT(t)

	action := actions.CompleteProfile{
		Name: "123456789012345678901234567890123456789012345678901", // 51 chars
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "name", "key")
}
