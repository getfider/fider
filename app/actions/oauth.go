package actions

import (
	"context"
	"strings"

	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/validate"
)

// CreateEditOAuthConfig is used to create/edit OAuth config
type CreateEditOAuthConfig struct {
	Input *models.CreateEditOAuthConfig
}

func NewCreateEditOAuthConfig() *CreateEditOAuthConfig {
	return &CreateEditOAuthConfig{
		Input: &models.CreateEditOAuthConfig{
			Logo: &models.ImageUpload{},
		},
	}
}

// Returns the struct to bind the request to
func (action *CreateEditOAuthConfig) BindTarget() interface{} {
	return action.Input
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateEditOAuthConfig) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *CreateEditOAuthConfig) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Input.Provider != "" {
		getConfig := &query.GetCustomOAuthConfigByProvider{Provider: action.Input.Provider}
		err := bus.Dispatch(ctx, getConfig)
		if err != nil {
			return validate.Error(err)
		}

		action.Input.ID = getConfig.Result.ID
		action.Input.Logo.BlobKey = getConfig.Result.LogoBlobKey
		if action.Input.ClientSecret == "" {
			action.Input.ClientSecret = getConfig.Result.ClientSecret
		}
	} else {
		action.Input.Provider = "_" + strings.ToLower(rand.String(10))
	}

	messages, err := validate.ImageUpload(action.Input.Logo, validate.ImageUploadOpts{
		IsRequired:   false,
		MinHeight:    24,
		MinWidth:     24,
		ExactRatio:   true,
		MaxKilobytes: 50,
	})
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("logo", messages...)

	if action.Input.Status != enum.OAuthConfigEnabled &&
		action.Input.Status != enum.OAuthConfigDisabled {
		result.AddFieldFailure("status", "Invalid status.")
	}

	if action.Input.DisplayName == "" {
		result.AddFieldFailure("displayName", "Display Name is required.")
	} else if len(action.Input.DisplayName) > 50 {
		result.AddFieldFailure("displayName", "Display Name must have less than 50 characters.")
	}

	if action.Input.ClientID == "" {
		result.AddFieldFailure("clientID", "Client ID is required.")
	} else if len(action.Input.ClientID) > 100 {
		result.AddFieldFailure("clientID", "Client ID must have less than 100 characters.")
	}

	if action.Input.ClientSecret == "" {
		result.AddFieldFailure("clientSecret", "Client Secret is required.")
	} else if len(action.Input.ClientSecret) > 500 {
		result.AddFieldFailure("clientSecret", "Client Secret must have less than 500 characters.")
	}

	if action.Input.Scope == "" {
		result.AddFieldFailure("scope", "Scope is required.")
	} else if len(action.Input.Scope) > 100 {
		result.AddFieldFailure("scope", "Scope must have less than 100 characters.")
	}

	if action.Input.AuthorizeURL == "" {
		result.AddFieldFailure("authorizeURL", "Authorize URL is required.")
	} else if messages := validate.URL(action.Input.AuthorizeURL); len(messages) > 0 {
		result.AddFieldFailure("authorizeURL", messages...)
	}

	if action.Input.TokenURL == "" {
		result.AddFieldFailure("tokenURL", "Token URL is required.")
	} else if messages := validate.URL(action.Input.TokenURL); len(messages) > 0 {
		result.AddFieldFailure("tokenURL", messages...)
	}

	if action.Input.ProfileURL != "" {
		if messages := validate.URL(action.Input.ProfileURL); len(messages) > 0 {
			result.AddFieldFailure("profileURL", messages...)
		}
	}

	if action.Input.JSONUserIDPath == "" {
		result.AddFieldFailure("jsonUserIDPath", "JSON User ID Path is required.")
	} else if len(action.Input.JSONUserIDPath) > 100 {
		result.AddFieldFailure("jsonUserIDPath", "JSON User ID Path must have less than 100 characters.")
	}

	if len(action.Input.JSONUserNamePath) > 100 {
		result.AddFieldFailure("jsonUserNamePath", "JSON User Name Path must have less than 100 characters.")
	}

	if len(action.Input.JSONUserEmailPath) > 100 {
		result.AddFieldFailure("jsonUserEmailPath", "JSON User Email Path must have less than 100 characters.")
	}

	return result
}
