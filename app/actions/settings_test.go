package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestInvalidUserNames(t *testing.T) {
	RegisterT(t)

	for _, name := range []string{
		"",
		"123456789012345678901234567890123456789012345678901", // 51 chars
	} {

		action := actions.NewUpdateUserSettings()
		action.Name = name
		action.AvatarType = enum.AvatarTypeGravatar
		result := action.Validate(context.Background(), &entity.User{})
		ExpectFailed(result, "name")
	}
}

func TestValidUserNames(t *testing.T) {
	RegisterT(t)

	for _, name := range []string{
		"Jon Snow",
		"Arya",
	} {
		action := actions.NewUpdateUserSettings()
		action.Name = name
		action.AvatarType = enum.AvatarTypeGravatar
		result := action.Validate(context.Background(), &entity.User{})
		ExpectSuccess(result)
	}
}

func TestInvalidSettings(t *testing.T) {
	RegisterT(t)

	for _, settings := range []map[string]string{
		{
			"bad_name": "3",
		},
		{
			enum.NotificationEventNewComment.UserSettingsKeyName: "4",
		},
	} {
		action := actions.NewUpdateUserSettings()
		action.Name = "John Snow"
		action.Settings = settings
		result := action.Validate(context.Background(), &entity.User{})
		ExpectFailed(result, "settings", "avatarType")
	}
}

func TestValidSettings(t *testing.T) {
	RegisterT(t)

	for _, settings := range []map[string]string{
		nil,
		{
			enum.NotificationEventNewPost.UserSettingsKeyName:      enum.NotificationEventNewPost.DefaultSettingValue,
			enum.NotificationEventNewComment.UserSettingsKeyName:   enum.NotificationEventNewComment.DefaultSettingValue,
			enum.NotificationEventChangeStatus.UserSettingsKeyName: enum.NotificationEventChangeStatus.DefaultSettingValue,
		},
		{
			enum.NotificationEventNewComment.UserSettingsKeyName: enum.NotificationEventNewComment.DefaultSettingValue,
		},
	} {
		action := actions.NewUpdateUserSettings()
		action.Name = "John Snow"
		action.Settings = settings
		action.AvatarType = enum.AvatarTypeGravatar

		result := action.Validate(context.Background(), &entity.User{
			AvatarBlobKey: "jon.png",
		})

		ExpectSuccess(result)
		Expect(action.Avatar.BlobKey).Equals("jon.png")
	}
}
