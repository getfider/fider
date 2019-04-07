package actions

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/validate"
)

// SignInByEmail happens when user request to sign in by email
type SignInByEmail struct {
	Model *models.SignInByEmail
}

// Initialize the model
func (input *SignInByEmail) Initialize() interface{} {
	input.Model = new(models.SignInByEmail)
	input.Model.VerificationKey = models.GenerateSecretKey()
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *SignInByEmail) IsAuthorized(user *models.User, services *app.Services) bool {
	return true
}

// Validate if current model is valid
func (input *SignInByEmail) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Email == "" {
		result.AddFieldFailure("email", "Email is required.")
		return result
	}

	if len(input.Model.Email) > 200 {
		result.AddFieldFailure("email", "Email must have less than 200 characters.")
	}

	messages := validate.Email(input.Model.Email)
	result.AddFieldFailure("email", messages...)

	return result
}

// CompleteProfile happens when users completes their profile during first time sign in
type CompleteProfile struct {
	Model *models.CompleteProfile
}

// Initialize the model
func (input *CompleteProfile) Initialize() interface{} {
	input.Model = new(models.CompleteProfile)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CompleteProfile) IsAuthorized(user *models.User, services *app.Services) bool {
	return true
}

// Validate if current model is valid
func (input *CompleteProfile) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	if len(input.Model.Name) > 50 {
		result.AddFieldFailure("name", "Name must have less than 50 characters.")
	}

	if input.Model.Key == "" {
		result.AddFieldFailure("key", "Key is required.")
	} else {
		findBySignIn := &query.GetVerificationByKey{Kind: models.EmailVerificationKindSignIn, Key: input.Model.Key}
		err1 := bus.Dispatch(services.Context, findBySignIn)

		findByUserInvitation := &query.GetVerificationByKey{Kind: models.EmailVerificationKindUserInvitation, Key: input.Model.Key}
		err2 := bus.Dispatch(services.Context, findByUserInvitation)

		if err1 == nil {
			input.Model.Email = findBySignIn.Result.Email
		} else if err2 == nil {
			input.Model.Email = findByUserInvitation.Result.Email
		} else {
			result.AddFieldFailure("key", "Key is invalid.")
		}
	}

	return result
}
