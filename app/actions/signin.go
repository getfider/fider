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
	Model *models.SignInByEmail
}

func NewSignInByEmail() *SignInByEmail {
	return &SignInByEmail{
		Model: &models.SignInByEmail{
			VerificationKey: models.GenerateSecretKey(),
		},
	}
}

// Returns the struct to bind the request to
func (action *SignInByEmail) BindTarget() interface{} {
	return action.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *SignInByEmail) IsAuthorized(ctx context.Context, user *models.User) bool {
	return true
}

// Validate if current model is valid
func (action *SignInByEmail) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Model.Email == "" {
		result.AddFieldFailure("email", "Email is required.")
		return result
	}

	if len(action.Model.Email) > 200 {
		result.AddFieldFailure("email", "Email must have less than 200 characters.")
	}

	messages := validate.Email(action.Model.Email)
	result.AddFieldFailure("email", messages...)

	return result
}

// CompleteProfile happens when users completes their profile during first time sign in
type CompleteProfile struct {
	Model *models.CompleteProfile
}

// Returns the struct to bind the request to
func (action *CompleteProfile) BindTarget() interface{} {
	action.Model = new(models.CompleteProfile)
	return action.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CompleteProfile) IsAuthorized(ctx context.Context, user *models.User) bool {
	return true
}

// Validate if current model is valid
func (action *CompleteProfile) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	if len(action.Model.Name) > 50 {
		result.AddFieldFailure("name", "Name must have less than 50 characters.")
	}

	if action.Model.Key == "" {
		result.AddFieldFailure("key", "Key is required.")
	} else {
		findBySignIn := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindSignIn, Key: action.Model.Key}
		err1 := bus.Dispatch(ctx, findBySignIn)

		findByUserInvitation := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindUserInvitation, Key: action.Model.Key}
		err2 := bus.Dispatch(ctx, findByUserInvitation)

		if err1 == nil {
			action.Model.Email = findBySignIn.Result.Email
		} else if err2 == nil {
			action.Model.Email = findByUserInvitation.Result.Email
		} else {
			result.AddFieldFailure("key", "Key is invalid.")
		}
	}

	return result
}
