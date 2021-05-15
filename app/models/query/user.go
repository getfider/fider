package query

import "github.com/getfider/fider/app/models/entity"

type CountUsers struct {
	Result int
}

type UserSubscribedTo struct {
	PostID int

	Result bool
}

type GetUserByAPIKey struct {
	APIKey string

	Result *entity.User
}

type GetCurrentUserSettings struct {
	Result map[string]string
}

type GetUserByID struct {
	UserID int

	Result *entity.User
}

type GetUserByEmail struct {
	Email string

	Result *entity.User
}

type GetUserByProvider struct {
	Provider string
	UID      string

	Result *entity.User
}

type GetAllUsers struct {
	Result []*entity.User
}
