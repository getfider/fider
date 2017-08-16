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
		err := c.Services().Tenants.UpdateSettings(c.Tenant().ID, input.Model.Title, input.Model.Invitation, input.Model.WelcomeMessage)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
