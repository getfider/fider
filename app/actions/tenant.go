package actions

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/validate"
)

//CreateTenant is the input model used to create a tenant
type CreateTenant struct {
	Model *models.CreateTenant
}

// NewModel initializes the model
func (input *CreateTenant) NewModel() interface{} {
	input.Model = new(models.CreateTenant)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateTenant) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (input *CreateTenant) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	var err error
	if input.Model.Token == "" {
		result.AddFieldFailure("token", "Please identify yourself before proceeding.")
	} else {
		if input.Model.UserClaims, err = jwt.DecodeOAuthClaims(input.Model.Token); err != nil {
			return validate.Error(err)
		}
	}

	if env.IsSingleHostMode() {
		input.Model.Subdomain = "default"
	}

	if input.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	subdomainResult := validate.Subdomain(services.Tenants, input.Model.Subdomain)
	if !subdomainResult.Ok {
		result.AddFieldFailure("subdomain", subdomainResult.Messages...)
	}

	input.Model.Subdomain = strings.ToLower(input.Model.Subdomain)

	return result
}

//UpdateTenantSettings is the input model used to update tenant settings
type UpdateTenantSettings struct {
	Model *models.UpdateTenantSettings
}

// NewModel initializes the model
func (input *UpdateTenantSettings) NewModel() interface{} {
	input.Model = new(models.UpdateTenantSettings)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UpdateTenantSettings) IsAuthorized(user *models.User) bool {
	return user != nil && user.Role == models.RoleAdministrator
}

// Validate is current model is valid
func (input *UpdateTenantSettings) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	input.Model.Title = strings.Trim(input.Model.Title, " ")
	input.Model.Invitation = strings.Trim(input.Model.Invitation, " ")
	input.Model.WelcomeMessage = strings.Trim(input.Model.WelcomeMessage, " ")

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	}

	if len(input.Model.Title) > 60 {
		result.AddFieldFailure("title", "Title must have less than 60 characters.")
	}

	if len(input.Model.Invitation) > 60 {
		result.AddFieldFailure("invitation", "Invitation must have less than 60 characters.")
	}

	return result
}
