package actions

import (
	"fmt"

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
func (input *UpdateUserSettings) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil
}

// Validate is current model is valid
func (input *UpdateUserSettings) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	if len(input.Model.Name) > 50 {
		result.AddFieldFailure("name", "Name must have less than 50 characters.")
	}

	if input.Model.Settings != nil {
		for k, v := range input.Model.Settings {
			ok := false
			for _, e := range models.AllNotificationEvents {
				if e.UserSettingsKeyName == k {
					ok = true
					if !e.Validate(v) {
						result.AddFieldFailure("settings", fmt.Sprintf("Settings %s has an invalid value %s.", k, v))
					}
				}
			}
			if ok == false {
				result.AddFieldFailure("settings", fmt.Sprintf("Unknown settings named %s.", k))
			}
		}
	}

	return result
}
