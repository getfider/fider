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
	Expect(action.VerificationCode).IsNotEmpty()
	Expect(len(action.VerificationCode)).Equals(6)
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

func TestVerifySignInCode_EmptyEmail(t *testing.T) {
	RegisterT(t)

	action := actions.VerifySignInCode{Email: "", Code: "123456"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "email")
}

func TestVerifySignInCode_InvalidEmail(t *testing.T) {
	RegisterT(t)

	action := actions.VerifySignInCode{Email: "invalid", Code: "123456"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "email")
}

func TestVerifySignInCode_EmptyCode(t *testing.T) {
	RegisterT(t)

	action := actions.VerifySignInCode{Email: "jon.snow@got.com", Code: ""}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "code")
}

func TestVerifySignInCode_InvalidCodeLength(t *testing.T) {
	RegisterT(t)

	action := actions.VerifySignInCode{Email: "jon.snow@got.com", Code: "12345"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "code")

	action2 := actions.VerifySignInCode{Email: "jon.snow@got.com", Code: "1234567"}
	result2 := action2.Validate(context.Background(), nil)
	ExpectFailed(result2, "code")
}

func TestVerifySignInCode_NonNumericCode(t *testing.T) {
	RegisterT(t)

	action := actions.VerifySignInCode{Email: "jon.snow@got.com", Code: "12345A"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "code")
}

func TestVerifySignInCode_ValidCodeAndEmail(t *testing.T) {
	RegisterT(t)

	action := actions.VerifySignInCode{Email: "jon.snow@got.com", Code: "123456"}
	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
}
