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
