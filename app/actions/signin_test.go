package actions_test

import (
	"context"
	"encoding/json"
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

func TestSignInByEmail_VerificationCodeNotOverwrittenByJSON(t *testing.T) {
	RegisterT(t)

	action := actions.NewSignInByEmail()
	originalCode := action.VerificationCode
	originalKey := action.LinkKey

	// Simulate attacker trying to overwrite internal fields via JSON binding
	payload := []byte(`{"email":"attacker@evil.com","VerificationCode":"000000","LinkKey":"attacker-key"}`)
	err := json.Unmarshal(payload, action)
	Expect(err).IsNil()

	// VerificationCode and LinkKey must NOT be overwritten
	Expect(action.VerificationCode).Equals(originalCode)
	Expect(action.LinkKey).Equals(originalKey)
	// Email should still bind normally
	Expect(action.Email).Equals("attacker@evil.com")
}

func TestSignInByEmailWithName_VerificationCodeNotOverwrittenByJSON(t *testing.T) {
	RegisterT(t)

	action := actions.NewSignInByEmailWithName()
	originalCode := action.VerificationCode
	originalKey := action.LinkKey

	payload := []byte(`{"email":"attacker@evil.com","name":"Attacker","VerificationCode":"000000","LinkKey":"attacker-key"}`)
	err := json.Unmarshal(payload, action)
	Expect(err).IsNil()

	Expect(action.VerificationCode).Equals(originalCode)
	Expect(action.LinkKey).Equals(originalKey)
	Expect(action.Email).Equals("attacker@evil.com")
	Expect(action.Name).Equals("Attacker")
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

func TestSignInByEmailWithName_EmptyEmailAndName(t *testing.T) {
	RegisterT(t)

	action := actions.SignInByEmailWithName{Email: "", Name: ""}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "email", "name")
}

func TestSignInByEmailWithName_InvalidEmail(t *testing.T) {
	RegisterT(t)

	action := actions.SignInByEmailWithName{Email: "invalid", Name: "Jon Snow"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "email")
}

func TestSignInByEmailWithName_EmptyName(t *testing.T) {
	RegisterT(t)

	action := actions.SignInByEmailWithName{Email: "jon.snow@got.com", Name: ""}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "name")
}

func TestSignInByEmailWithName_LongName(t *testing.T) {
	RegisterT(t)

	// 101 characters
	longName := "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"
	action := actions.SignInByEmailWithName{Email: "jon.snow@got.com", Name: longName}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "name")
}

func TestSignInByEmailWithName_ValidEmailAndName(t *testing.T) {
	RegisterT(t)

	action := actions.NewSignInByEmailWithName()
	action.Email = "jon.snow@got.com"
	action.Name = "Jon Snow"

	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
	Expect(action.VerificationCode).IsNotEmpty()
	Expect(len(action.VerificationCode)).Equals(6)
}
