package cmd

import "github.com/getfider/fider/app/models"

type BlockUser struct {
	UserID int
}

type UnblockUser struct {
	UserID int
}

type RegenerateAPIKey struct {
	Result string
}

type DeleteCurrentUser struct {
}

type ChangeUserRole struct {
	UserID int
	Role   models.Role
}

type ChangeUserEmail struct {
	UserID int
	Email  string
}

type UpdateCurrentUserSettings struct {
	Settings map[string]string
}

type RegisterUser struct {
	User *models.User
}

type RegisterUserProvider struct {
	UserID       int
	ProviderName string
	ProviderUID  string
}

type UpdateCurrentUser struct {
	Name       string
	AvatarType models.AvatarType
	Avatar     *models.ImageUpload
}
