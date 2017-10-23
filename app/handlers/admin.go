package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/web"
)

// UpdateSettings update current tenant' settings
func UpdateSettings() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.UpdateTenantSettings)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}
		err := c.Services().Tenants.UpdateSettings(input.Model)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ManageMembers is the page used by administrators to change member's role
func ManageMembers() web.HandlerFunc {
	return func(c web.Context) error {
		users, err := c.Services().Users.GetAll()
		if err != nil {
			return c.Failure(err)
		}
		return c.Page(web.Map{
			"users": users,
		})
	}
}
