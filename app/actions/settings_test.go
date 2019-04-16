package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/enum"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestInvalidUserNames(t *testing.T) {
	RegisterT(t)

	for _, name := range []string{
		"",
		"123456789012345678901234567890123456789012345678901", // 51 chars
	} {

		action := actions.UpdateUserSettings{}
		action.Initialize()
		action.Model.Name = name
		action.Model.AvatarType = enum.AvatarTypeGravatar
		result := action.Validate(context.Background(), &models.User{})
		ExpectFailed(result, "name")
	}
}

func TestValidUserNames(t *testing.T) {
	RegisterT(t)

	for _, name := range []string{
		"Jon Snow",
		"Arya",
	} {
		action := actions.UpdateUserSettings{}
		action.Initialize()
		action.Model.Name = name
		action.Model.AvatarType = enum.AvatarTypeGravatar
		result := action.Validate(context.Background(), &models.User{})
		ExpectSuccess(result)
	}
}

func TestInvalidSettings(t *testing.T) {
	RegisterT(t)

	for _, settings := range []map[string]string{
		map[string]string{
			"bad_name": "3",
		},
		map[string]string{
			enum.NotificationEventNewComment.UserSettingsKeyName: "4",
		},
	} {
		action := actions.UpdateUserSettings{}
		action.Initialize()
		action.Model.Name = "John Snow"
		action.Model.Settings = settings
		result := action.Validate(context.Background(), &models.User{})
		ExpectFailed(result, "settings", "avatarType")
	}
}

func TestValidSettings(t *testing.T) {
	RegisterT(t)

	for _, settings := range []map[string]string{
		nil,
		map[string]string{
			enum.NotificationEventNewPost.UserSettingsKeyName:      enum.NotificationEventNewPost.DefaultSettingValue,
			enum.NotificationEventNewComment.UserSettingsKeyName:   enum.NotificationEventNewComment.DefaultSettingValue,
			enum.NotificationEventChangeStatus.UserSettingsKeyName: enum.NotificationEventChangeStatus.DefaultSettingValue,
		},
		map[string]string{
			enum.NotificationEventNewComment.UserSettingsKeyName: enum.NotificationEventNewComment.DefaultSettingValue,
		},
	} {
		action := actions.UpdateUserSettings{}
		action.Initialize()
		action.Model.Name = "John Snow"
		action.Model.Settings = settings
		action.Model.AvatarType = enum.AvatarTypeGravatar

		result := action.Validate(context.Background(), &models.User{
			AvatarBlobKey: "jon.png",
		})

		ExpectSuccess(result)
		Expect(action.Model.Avatar.BlobKey).Equals("jon.png")
	}
}
