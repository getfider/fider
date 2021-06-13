package actions

import (
	"context"
	"fmt"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/i18n"
	"github.com/getfider/fider/app/pkg/validate"
)

// UpdateUserSettings happens when users updates their settings
type UpdateUserSettings struct {
	Name       string            `json:"name"`
	AvatarType enum.AvatarType   `json:"avatarType"`
	Avatar     *dto.ImageUpload  `json:"avatar"`
	Settings   map[string]string `json:"settings"`
}

func NewUpdateUserSettings() *UpdateUserSettings {
	return &UpdateUserSettings{
		Avatar: &dto.ImageUpload{},
	}
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *UpdateUserSettings) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil
}

// Validate if current model is valid
func (action *UpdateUserSettings) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	if action.AvatarType < 1 || action.AvatarType > 3 {
		result.AddFieldFailure("avatarType", "Invalid avatar type.")
	}

	if len(action.Name) > 50 {
		result.AddFieldFailure("name", "Name must have less than 50 characters.")
	}

	action.Avatar.BlobKey = user.AvatarBlobKey
	messages, err := validate.ImageUpload(action.Avatar, validate.ImageUploadOpts{
		IsRequired:   action.AvatarType == enum.AvatarTypeCustom,
		MinHeight:    50,
		MinWidth:     50,
		ExactRatio:   true,
		MaxKilobytes: 100,
	})
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("avatar", messages...)

	if action.Settings != nil {
		for k, v := range action.Settings {
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
				result.AddFieldFailure("settings", i18n.T(ctx, "validation.settings.unknown", i18n.Params{"name": k}))
			}
		}
	}

	return result
}
