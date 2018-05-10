package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// GeneralSettingsPage is the general settings page
func GeneralSettingsPage() web.HandlerFunc {
	return func(c web.Context) error {
		publicIP, err := env.GetPublicIP()
		if err != nil {
			c.Logger().Error(err)
		}

		return c.Page(web.Props{
			Title: "General · Site Settings",
			Data: web.Map{
				"publicIP": publicIP,
			},
		})
	}
}

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

// UpdatePrivacy update current tenant's privacy settings
func UpdatePrivacy() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.UpdateTenantPrivacy)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tenants.UpdatePrivacy(input.Model)
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

		return c.Page(web.Props{
			Title: "Manage Members · Site Settings",
			Data: web.Map{
				"users": users,
			},
		})
	}
}
