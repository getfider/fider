package handlers

import (
	"net/http"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// GeneralSettingsPage is the general settings page
func GeneralSettingsPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/GeneralSettings.page",
			Title: "General · Site Settings",
		})
	}
}

// AdvancedSettingsPage is the advanced settings page
func AdvancedSettingsPage() web.HandlerFunc {
	return func(c *web.Context) error {
		billingState := &query.GetStripeBillingState{}
		if err := bus.Dispatch(c, billingState); err != nil {
			return c.Failure(err)
		}

		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/AdvancedSettings.page",
			Title: "Advanced · Site Settings",
			Data: web.Map{
				"customCSS":              c.Tenant().CustomCSS,
				"allowedSchemes":         c.Tenant().AllowedSchemes,
				"licenseKey":             billingState.Result.LicenseKey,
				"hasCommercialFeatures": c.Tenant().HasCommercialFeatures,
			},
		})
	}
}

// UpdateSettings update current tenant' settings
func UpdateSettings() web.HandlerFunc {
	return func(c *web.Context) error {
		action := actions.NewUpdateTenantSettings()
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c,
			&cmd.UploadImage{
				Image:  action.Logo,
				Folder: "logos",
			},
			&cmd.UpdateTenantSettings{
				Logo:           action.Logo,
				Title:          action.Title,
				Invitation:     action.Invitation,
				WelcomeMessage: action.WelcomeMessage,
				WelcomeHeader:  action.WelcomeHeader,
				CNAME:          action.CNAME,
				Locale:         action.Locale,
			},
		); err != nil {
			return c.Failure(err)
		}

		// Handle userlist.
		if env.Config.UserList.Enabled {
			c.Enqueue(tasks.UserListUpdateCompany(&dto.UserListUpdateCompany{
				TenantID: c.Tenant().ID,
				Name:     action.Title,
			}))
		}

		return c.Ok(web.Map{})
	}
}

// UpdateAdvancedSettings update current tenant' advanced settings
func UpdateAdvancedSettings() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.UpdateTenantAdvancedSettings)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, &cmd.UpdateTenantAdvancedSettings{
			CustomCSS:      action.CustomCSS,
			AllowedSchemes: action.AllowedSchemes,
		}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UpdatePrivacySettings update current tenant's privacy settings
func UpdatePrivacySettings() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.UpdateTenantPrivacySettings)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		updateSettings := &cmd.UpdateTenantPrivacySettings{
			IsPrivate:           action.IsPrivate,
			IsFeedEnabled:       action.IsFeedEnabled,
			IsModerationEnabled: action.IsModerationEnabled,
		}
		if err := bus.Dispatch(c, updateSettings); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UpdateEmailAuthAllowed update current tenant's allow email auth settings
func UpdateEmailAuthAllowed() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.UpdateTenantEmailAuthAllowed)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		updateSettings := &cmd.UpdateTenantEmailAuthAllowedSettings{
			IsEmailAuthAllowed: action.IsEmailAuthAllowed,
		}
		if err := bus.Dispatch(c, updateSettings); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ManageMembers is the page used by administrators to change member's role
func ManageMembers() web.HandlerFunc {
	return func(c *web.Context) error {
		// Only load first page for initial page load - subsequent pagination handled by API
		page, _ := c.QueryParamAsInt("page")
		if page <= 0 {
			page = 1
		}

		searchUsers := &query.SearchUsers{
			Query: c.QueryParam("query"),
			Roles: c.QueryParamAsArray("roles"),
			Page:  page,
			Limit: 10,
		}

		if err := bus.Dispatch(c, searchUsers); err != nil {
			return c.Failure(err)
		}

		// Create an array of UserWithEmail structs from the searchUsers.Result
		allUsersWithEmail := make([]entity.UserWithEmail, len(searchUsers.Result))
		for i, user := range searchUsers.Result {
			allUsersWithEmail[i] = entity.UserWithEmail{
				User: user,
			}
		}

		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/ManageMembers.page",
			Title: "Manage Members · Site Settings",
			Data: web.Map{
				"users":      allUsersWithEmail,
				"totalPages": (searchUsers.TotalCount + 10 - 1) / 10,
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

		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/ManageAuthentication.page",
			Title: "Authentication · Site Settings",
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
		action := actions.NewCreateEditOAuthConfig()
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c,
			&cmd.UploadImage{
				Image:  action.Logo,
				Folder: "logos",
			},
			&cmd.SaveCustomOAuthConfig{
				ID:                action.ID,
				Logo:              action.Logo,
				Provider:          action.Provider,
				Status:            action.Status,
				DisplayName:       action.DisplayName,
				ClientID:          action.ClientID,
				ClientSecret:      action.ClientSecret,
				AuthorizeURL:      action.AuthorizeURL,
				TokenURL:          action.TokenURL,
				Scope:             action.Scope,
				ProfileURL:        action.ProfileURL,
				IsTrusted:         action.IsTrusted,
				JSONUserIDPath:    action.JSONUserIDPath,
				JSONUserNamePath:  action.JSONUserNamePath,
				JSONUserEmailPath: action.JSONUserEmailPath,
			},
		); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// SetSystemProviderStatus is used to enable/disable built-in OAuth providers for a tenant
func SetSystemProviderStatus() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.SetSystemProviderStatus)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, &cmd.SetTenantProviderStatus{
			Provider:  action.Provider,
			IsEnabled: action.IsEnabled,
		}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
