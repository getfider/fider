package query

import "github.com/getfider/fider/app/models/entities"

type CountUsers struct {
	Result int
}

type UserSubscribedTo struct {
	PostID int

	Result bool
}

type GetUserByAPIKey struct {
	APIKey string

	Result *entities.User
}

type GetCurrentUserSettings struct {
	Result map[string]string
}

type GetUserByID struct {
	UserID int

	Result *entities.User
}

type GetUserByEmail struct {
	Email string

	Result *entities.User
}

type GetUserByProvider struct {
	Provider string
	UID      string

	Result *entities.User
}

type GetAllUsers struct {
	Result []*entities.User
}
