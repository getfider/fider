package handlers

import (
	"github.com/getfider/fider/app/models/im"
	"github.com/getfider/fider/app/pkg/web"
)

// UpdateSettings update current tenant' settings
func UpdateSettings() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(im.UpdateTenantSettings)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tenants.UpdateSettings(c.Tenant().ID, input.Title, input.Invitation, input.WelcomeMessage)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
