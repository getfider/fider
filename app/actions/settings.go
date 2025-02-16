package actions

import (
	"context"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/i18n"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/profanity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/validate"
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
		result.AddFieldFailure("name", propertyIsRequired(ctx, "name"))
	} else if matches, err := profanity.ContainsProfanity(ctx, action.Name); err == nil && len(matches) > 0 {
		result.AddFieldFailure("content", i18n.T(ctx, "validation.custom.containsprofanity"))
	}

	if action.AvatarType < 1 || action.AvatarType > 3 {
		result.AddFieldFailure("avatarType", propertyIsInvalid(ctx, "avatarType"))
	}

	if len(action.Name) > 50 {
		result.AddFieldFailure("name", propertyMaxStringLen(ctx, "name", 50))
	}

	action.Avatar.BlobKey = user.AvatarBlobKey
	messages, err := validate.ImageUpload(ctx, action.Avatar, validate.ImageUploadOpts{
		IsRequired:   action.AvatarType == enum.AvatarTypeCustom,
		MinHeight:    50,
		MinWidth:     50,
		ExactRatio:   true,
		MaxKilobytes: 5000,
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
						result.AddFieldFailure("settings", i18n.T(ctx, "validation.invalidvalue", i18n.Params{"name": k}, i18n.Params{"value": v}))
					}
				}
			}
			if !ok {
				result.AddFieldFailure("settings", i18n.T(ctx, "validation.custom.unknownsettings", i18n.Params{"name": k}))
			}
		}
	}

	return result
}
