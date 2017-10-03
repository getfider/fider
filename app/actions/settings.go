package actions

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

// UpdateUserSettings happens when users updates their settings
type UpdateUserSettings struct {
	Model *models.UpdateUserSettings
}

// Initialize the model
func (input *UpdateUserSettings) Initialize() interface{} {
	input.Model = new(models.UpdateUserSettings)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UpdateUserSettings) IsAuthorized(user *models.User) bool {
	return user != nil
}

// Validate is current model is valid
func (input *UpdateUserSettings) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	if len(input.Model.Name) > 50 {
		result.AddFieldFailure("name", "Name must be less than 50 characters.")
	}

	return result
}
