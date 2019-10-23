package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// GeneralSettingsPage is the general settings page
func GeneralSettingsPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(web.Props{
			Title:     "General · Site Settings",
			ChunkName: "GeneralSettings.page",
		})
	}
}

// AdvancedSettingsPage is the advanced settings page
func AdvancedSettingsPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(web.Props{
			Title:     "Advanced · Site Settings",
			ChunkName: "AdvancedSettings.page",
			Data: web.Map{
				"customCSS": c.Tenant().CustomCSS,
			},
		})
	}
}

// UpdateSettings update current tenant' settings
func UpdateSettings() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.UpdateTenantSettings)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c,
			&cmd.UploadImage{
				Image:  input.Model.Logo,
				Folder: "logos",
			},
			&cmd.UpdateTenantSettings{
				Settings: input.Model,
			},
		); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UpdateAdvancedSettings update current tenant' advanced settings
func UpdateAdvancedSettings() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.UpdateTenantAdvancedSettings)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, &cmd.UpdateTenantAdvancedSettings{Settings: input.Model}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UpdatePrivacy update current tenant's privacy settings
func UpdatePrivacy() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.UpdateTenantPrivacy)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		updateSettings := &cmd.UpdateTenantPrivacySettings{Settings: input.Model}
		if err := bus.Dispatch(c, updateSettings); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ManageMembers is the page used by administrators to change member's role
func ManageMembers() web.HandlerFunc {
	return func(c *web.Context) error {
		allUsers := &query.GetAllUsers{}
		if err := bus.Dispatch(c, allUsers); err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:     "Manage Members · Site Settings",
			ChunkName: "ManageMembers.page",
			Data: web.Map{
				"users": allUsers.Result,
			},
		})
	}
}

// ManageAuthentication is the page used by administrators to change site authentication settings
func ManageAuthentication() web.HandlerFunc {
	return func(c *web.Context) error {
		listProviders := &query.ListAllOAuthProviders{}
		if err := bus.Dispatch(c, listProviders); err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:     "Authentication · Site Settings",
			ChunkName: "ManageAuthentication.page",
			Data: web.Map{
				"providers": listProviders.Result,
			},
		})
	}
}

// GetOAuthConfig returns OAuth config based on given provider
func GetOAuthConfig() web.HandlerFunc {
	return func(c *web.Context) error {
		getConfig := &query.GetCustomOAuthConfigByProvider{
			Provider: c.Param("provider"),
		}
		if err := bus.Dispatch(c, getConfig); err != nil {
			return c.Failure(err)
		}

		return c.Ok(getConfig.Result)
	}
}

// SaveOAuthConfig is used to create/edit OAuth configurations
func SaveOAuthConfig() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.CreateEditOAuthConfig)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c,
			&cmd.UploadImage{
				Image:  input.Model.Logo,
				Folder: "logos",
			},
			&cmd.SaveCustomOAuthConfig{
				Config: input.Model,
			},
		); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
