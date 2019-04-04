package query

import "github.com/getfider/fider/app/models"

type CountUsers struct {
	Result int
}

type GetUserByAPIKey struct {
	APIKey string

	Result *models.User
}
