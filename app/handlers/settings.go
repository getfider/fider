package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/web"
)

// UpdateUserSettings updates current user settings
func UpdateUserSettings() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.UpdateUserSettings)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Users.Update(c.User().ID, input.Model)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ChangeUserRole changes given user role
func ChangeUserRole() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.ChangeUserRole)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Users.ChangeRole(input.Model.UserID, input.Model.Role)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
