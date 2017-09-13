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

	emailResult := validate.Email(input.Model.Email)
	if !emailResult.Ok {
		result.AddFieldFailure("email", emailResult.Messages...)
	}

	input.Model.Email = strings.ToLower(input.Model.Email)
	return result
}
