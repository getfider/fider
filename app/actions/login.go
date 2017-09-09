package actions

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/uuid"
	"github.com/getfider/fider/app/pkg/validate"
)

// LoginByEmail happens when user request to log in by email
type LoginByEmail struct {
	Model *models.LoginByEmail
}

// Initialize the model
func (input *LoginByEmail) Initialize() interface{} {
	input.Model = new(models.LoginByEmail)
	input.Model.VerificationKey = uuid.NewV4().String()
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *LoginByEmail) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (input *LoginByEmail) Validate(services *app.Services) *validate.Result {
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
