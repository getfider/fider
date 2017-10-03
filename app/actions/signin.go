package actions

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/uuid"
	"github.com/getfider/fider/app/pkg/validate"
)

// SignInByEmail happens when user request to sign in by email
type SignInByEmail struct {
	Model *models.SignInByEmail
}

// Initialize the model
func (input *SignInByEmail) Initialize() interface{} {
	input.Model = new(models.SignInByEmail)
	input.Model.VerificationKey = strings.Replace(uuid.NewV4().String(), "-", "", 4)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *SignInByEmail) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (input *SignInByEmail) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Email == "" {
		result.AddFieldFailure("email", "E-mail is required.")
		return result
	}

	if len(input.Model.Email) > 200 {
		result.AddFieldFailure("email", "E-mail must be less than 200 characters.")
	}

	emailResult := validate.Email(input.Model.Email)
	if !emailResult.Ok {
		result.AddFieldFailure("email", emailResult.Messages...)
	}

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
func (input *CompleteProfile) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (input *CompleteProfile) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	if len(input.Model.Name) > 50 {
		result.AddFieldFailure("name", "Name must be less than 50 characters.")
	}

	if input.Model.Key == "" {
		result.AddFieldFailure("key", "Key is required.")
	} else {
		request, err := services.Tenants.FindVerificationByKey(input.Model.Key)
		if err != nil {
			if err == app.ErrNotFound {
				result.AddFieldFailure("key", "Key is invalid.")
			} else {
				return validate.Error(err)
			}
		} else {
			input.Model.Email = request.Email
		}
	}

	return result
}
