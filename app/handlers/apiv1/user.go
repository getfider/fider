package apiv1

import (
	"github.com/getfider/fider/app/pkg/web"
)

// ListUsers returns all registered users
func ListUsers() web.HandlerFunc {
	return func(c web.Context) error {
		users, err := c.Services().Users.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(users)
	}
}
