package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
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

func TestCompleteProfile_UnknownKey(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		return app.ErrNotFound
	})

	action := actions.CompleteProfile{Name: "Jon Snow", Key: "1234567890"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "key")
}

func TestCompleteProfile_ValidKey(t *testing.T) {
	RegisterT(t)

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		if q.Key == key && q.Kind == enum.EmailVerificationKindSignIn {
			q.Result = &entity.EmailVerification{
				Key:   q.Key,
				Kind:  q.Kind,
				Email: "jon.snow@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	action := actions.CompleteProfile{Name: "Jon Snow", Key: key}
	result := action.Validate(context.Background(), nil)

	ExpectSuccess(result)
	Expect(action.Email).Equals("jon.snow@got.com")
}

func TestCompleteProfile_UserInvitation_ValidKey(t *testing.T) {
	RegisterT(t)

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		if q.Key == key && q.Kind == enum.EmailVerificationKindUserInvitation {
			q.Result = &entity.EmailVerification{
				Key:   q.Key,
				Kind:  q.Kind,
				Email: "jon.snow@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	action := actions.CompleteProfile{Name: "Jon Snow", Key: key}
	result := action.Validate(context.Background(), nil)

	ExpectSuccess(result)
	Expect(action.Email).Equals("jon.snow@got.com")
}
