package actions

import (
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

// Initialize the model
func (input *CreateTenant) Initialize() interface{} {
	input.Model = new(models.CreateTenant)
	input.Model.VerificationKey = models.GenerateSecretKey()
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateTenant) IsAuthorized(user *models.User, services *app.Services) bool {
	return true
}

// Validate is current model is valid
func (input *CreateTenant) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	var err error
	if input.Model.Name == "" && input.Model.Email == "" {
		if input.Model.Token == "" {
			result.AddFieldFailure("token", "Please identify yourself before proceeding.")
		} else {
			if input.Model.UserClaims, err = jwt.DecodeOAuthClaims(input.Model.Token); err != nil {
				return validate.Error(err)
			}
		}
	} else {
		if input.Model.Email == "" {
			result.AddFieldFailure("email", "Email is required.")
		} else {
			messages := validate.Email(input.Model.Email)
			result.AddFieldFailure("email", messages...)
		}

		if input.Model.Name == "" {
			result.AddFieldFailure("name", "Name is required.")
		}
		if len(input.Model.Name) > 60 {
			result.AddFieldFailure("name", "Name must have less than 60 characters.")
		}
	}

	if env.IsSingleHostMode() {
		input.Model.Subdomain = "default"
	}

	if input.Model.TenantName == "" {
		result.AddFieldFailure("tenantName", "Name is required.")
	}

	messages, err := validate.Subdomain(services.Tenants, input.Model.Subdomain)
	if err != nil {
		return validate.Error(err)
	}

	result.AddFieldFailure("subdomain", messages...)

	if env.HasLegal() && !input.Model.LegalAgreement {
		result.AddFieldFailure("legalAgreement", "You must agree before proceeding.")
	}

	return result
}

//UpdateTenantSettings is the input model used to update tenant settings
type UpdateTenantSettings struct {
	Model *models.UpdateTenantSettings
}

// Initialize the model
func (input *UpdateTenantSettings) Initialize() interface{} {
	input.Model = new(models.UpdateTenantSettings)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UpdateTenantSettings) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.Role == models.RoleAdministrator
}

// Validate is current model is valid
func (input *UpdateTenantSettings) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	messages, err := validate.ImageUpload(input.Model.Logo, 200, 200, 100)
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("logo", messages...)

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	}

	if len(input.Model.Title) > 60 {
		result.AddFieldFailure("title", "Title must have less than 60 characters.")
	}

	if len(input.Model.Invitation) > 60 {
		result.AddFieldFailure("invitation", "Invitation must have less than 60 characters.")
	}

	if input.Model.CNAME != "" {
		messages := validate.CNAME(services.Tenants, input.Model.CNAME)
		result.AddFieldFailure("cname", messages...)
	}

	return result
}

//UpdateTenantAdvancedSettings is the input model used to update tenant advanced settings
type UpdateTenantAdvancedSettings struct {
	Model *models.UpdateTenantAdvancedSettings
}

// Initialize the model
func (input *UpdateTenantAdvancedSettings) Initialize() interface{} {
	input.Model = new(models.UpdateTenantAdvancedSettings)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UpdateTenantAdvancedSettings) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.Role == models.RoleAdministrator
}

// Validate is current model is valid
func (input *UpdateTenantAdvancedSettings) Validate(user *models.User, services *app.Services) *validate.Result {
	return validate.Success()
}

//UpdateTenantPrivacy is the input model used to update tenant privacy settings
type UpdateTenantPrivacy struct {
	Model *models.UpdateTenantPrivacy
}

// Initialize the model
func (input *UpdateTenantPrivacy) Initialize() interface{} {
	input.Model = new(models.UpdateTenantPrivacy)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UpdateTenantPrivacy) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.Role == models.RoleAdministrator
}

// Validate is current model is valid
func (input *UpdateTenantPrivacy) Validate(user *models.User, services *app.Services) *validate.Result {
	return validate.Success()
}
