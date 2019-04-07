package query

import "github.com/getfider/fider/app/models"

type CountUsers struct {
	Result int
}

type UserSubscribedTo struct {
	PostID int

	Result bool
}

type GetUserByAPIKey struct {
	APIKey string

	Result *models.User
}

type GetCurrentUserSettings struct {
	Result map[string]string
}
