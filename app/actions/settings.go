package actions

import (
	"context"
	"fmt"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/validate"
)

// UpdateUserSettings happens when users updates their settings
type UpdateUserSettings struct {
	Input *models.UpdateUserSettings
}

func NewUpdateUserSettings() *UpdateUserSettings {
	return &UpdateUserSettings{
		Input: &models.UpdateUserSettings{
			Avatar: &models.ImageUpload{},
		},
	}
}

// Returns the struct to bind the request to
func (action *UpdateUserSettings) BindTarget() interface{} {
	return action.Input
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *UpdateUserSettings) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil
}

// Validate if current model is valid
func (action *UpdateUserSettings) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Input.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	if action.Input.AvatarType < 1 || action.Input.AvatarType > 3 {
		result.AddFieldFailure("avatarType", "Invalid avatar type.")
	}

	if len(action.Input.Name) > 50 {
		result.AddFieldFailure("name", "Name must have less than 50 characters.")
	}

	action.Input.Avatar.BlobKey = user.AvatarBlobKey
	messages, err := validate.ImageUpload(action.Input.Avatar, validate.ImageUploadOpts{
		IsRequired:   action.Input.AvatarType == enum.AvatarTypeCustom,
		MinHeight:    50,
		MinWidth:     50,
		ExactRatio:   true,
		MaxKilobytes: 100,
	})
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("avatar", messages...)

	if action.Input.Settings != nil {
		for k, v := range action.Input.Settings {
			ok := false
			for _, e := range enum.AllNotificationEvents {
				if e.UserSettingsKeyName == k {
					ok = true
					if !e.Validate(v) {
						result.AddFieldFailure("settings", fmt.Sprintf("Settings %s has an invalid value %s.", k, v))
					}
				}
			}
			if !ok {
				result.AddFieldFailure("settings", fmt.Sprintf("Unknown settings named %s.", k))
			}
		}
	}

	return result
}
