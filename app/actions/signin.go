package actions

import (
	"context"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/validate"
)

// SignInByEmail happens when user request to sign in by email
type SignInByEmail struct {
	Input *models.SignInByEmail
}

func NewSignInByEmail() *SignInByEmail {
	return &SignInByEmail{
		Input: &models.SignInByEmail{
			VerificationKey: models.GenerateSecretKey(),
		},
	}
}

// Returns the struct to bind the request to
func (action *SignInByEmail) BindTarget() interface{} {
	return action.Input
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *SignInByEmail) IsAuthorized(ctx context.Context, user *models.User) bool {
	return true
}

// Validate if current model is valid
func (action *SignInByEmail) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Input.Email == "" {
		result.AddFieldFailure("email", "Email is required.")
		return result
	}

	if len(action.Input.Email) > 200 {
		result.AddFieldFailure("email", "Email must have less than 200 characters.")
	}

	messages := validate.Email(action.Input.Email)
	result.AddFieldFailure("email", messages...)

	return result
}

// CompleteProfile happens when users completes their profile during first time sign in
type CompleteProfile struct {
	Input *models.CompleteProfile
}

// Returns the struct to bind the request to
func (action *CompleteProfile) BindTarget() interface{} {
	action.Input = new(models.CompleteProfile)
	return action.Input
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CompleteProfile) IsAuthorized(ctx context.Context, user *models.User) bool {
	return true
}

// Validate if current model is valid
func (action *CompleteProfile) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Input.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	if len(action.Input.Name) > 50 {
		result.AddFieldFailure("name", "Name must have less than 50 characters.")
	}

	if action.Input.Key == "" {
		result.AddFieldFailure("key", "Key is required.")
	} else {
		findBySignIn := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindSignIn, Key: action.Input.Key}
		err1 := bus.Dispatch(ctx, findBySignIn)

		findByUserInvitation := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindUserInvitation, Key: action.Input.Key}
		err2 := bus.Dispatch(ctx, findByUserInvitation)

		if err1 == nil {
			action.Input.Email = findBySignIn.Result.Email
		} else if err2 == nil {
			action.Input.Email = findByUserInvitation.Result.Email
		} else {
			result.AddFieldFailure("key", "Key is invalid.")
		}
	}

	return result
}
