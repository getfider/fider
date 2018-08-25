package actions

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/validate"
)

// CreateEditOAuthConfig is used to create/edit OAuth config
type CreateEditOAuthConfig struct {
	Model *models.CreateEditOAuthConfig
}

// Initialize the model
func (input *CreateEditOAuthConfig) Initialize() interface{} {
	input.Model = new(models.CreateEditOAuthConfig)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateEditOAuthConfig) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsAdministrator()
}

// Validate is current model is valid
func (input *CreateEditOAuthConfig) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Provider != "" {
		config, err := services.Tenants.GetOAuthConfigByProvider(input.Model.Provider)
		if err != nil {
			return validate.Error(err)
		}

		input.Model.ID = config.ID
		if input.Model.ClientSecret == "" {
			input.Model.ClientSecret = config.ClientSecret
		}
	} else {
		input.Model.Provider = "_" + strings.ToLower(rand.String(10))
	}

	messages, err := validate.ImageUpload(input.Model.Logo, 24, 24, 50)
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("logo", messages...)

	if input.Model.Status != models.OAuthConfigEnabled &&
		input.Model.Status != models.OAuthConfigDisabled {
		result.AddFieldFailure("status", "Invalid status.")
	}

	if input.Model.DisplayName == "" {
		result.AddFieldFailure("displayName", "Display Name is required.")
	} else if len(input.Model.DisplayName) > 50 {
		result.AddFieldFailure("displayName", "Display Name must have less than 50 characters.")
	}

	if input.Model.ClientID == "" {
		result.AddFieldFailure("clientID", "Client ID is required.")
	} else if len(input.Model.ClientID) > 100 {
		result.AddFieldFailure("clientID", "Client ID must have less than 100 characters.")
	}

	if input.Model.ClientSecret == "" {
		result.AddFieldFailure("clientSecret", "Client Secret is required.")
	} else if len(input.Model.ClientSecret) > 500 {
		result.AddFieldFailure("clientSecret", "Client Secret must have less than 500 characters.")
	}

	if input.Model.Scope == "" {
		result.AddFieldFailure("scope", "Scope is required.")
	} else if len(input.Model.Scope) > 100 {
		result.AddFieldFailure("scope", "Scope must have less than 100 characters.")
	}

	if input.Model.AuthorizeURL == "" {
		result.AddFieldFailure("authorizeURL", "Authorize URL is required.")
	} else if messages := validate.URL(input.Model.AuthorizeURL); len(messages) > 0 {
		result.AddFieldFailure("authorizeURL", messages...)
	}

	if input.Model.TokenURL == "" {
		result.AddFieldFailure("tokenURL", "Token URL is required.")
	} else if messages := validate.URL(input.Model.TokenURL); len(messages) > 0 {
		result.AddFieldFailure("tokenURL", messages...)
	}

	if input.Model.ProfileURL != "" {
		if messages := validate.URL(input.Model.ProfileURL); len(messages) > 0 {
			result.AddFieldFailure("profileURL", messages...)
		}
	}

	if input.Model.JSONUserIDPath == "" {
		result.AddFieldFailure("jsonUserIDPath", "JSON User ID Path is required.")
	} else if len(input.Model.JSONUserIDPath) > 100 {
		result.AddFieldFailure("jsonUserIDPath", "JSON User ID Path must have less than 100 characters.")
	}

	if len(input.Model.JSONUserNamePath) > 100 {
		result.AddFieldFailure("jsonUserNamePath", "JSON User Name Path must have less than 100 characters.")
	}

	if len(input.Model.JSONUserEmailPath) > 100 {
		result.AddFieldFailure("jsonUserEmailPath", "JSON User Email Path must have less than 100 characters.")
	}

	return result
}
