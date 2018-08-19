package actions

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
)

// APIAuthorize is used during API Authorize process
type APIAuthorize struct {
	Model *models.APIAuthorize
	User  *models.User
}

// Initialize the model
func (input *APIAuthorize) Initialize() interface{} {
	input.Model = new(models.APIAuthorize)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *APIAuthorize) IsAuthorized(user *models.User, services *app.Services) bool {
	return true
}

// Validate is current model is valid
func (input *APIAuthorize) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.APIKey == "" {
		result.AddFieldFailure("apiKey", "API Key is required.")
	} else {
		user, err := services.Users.GetByAPIKey(input.Model.APIKey)
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				result.AddFieldFailure("apiKey", "API Key is invalid.")
			} else {
				return validate.Error(err)
			}
		} else {
			input.User = user
		}
	}

	return result
}
