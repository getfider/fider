package actions

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/uuid"
	"github.com/getfider/fider/app/pkg/validate"
)

//CreateTenant is the input model used to create a tenant
type CreateTenant struct {
	Model *models.CreateTenant
}

// Initialize the model
func (input *CreateTenant) Initialize() interface{} {
	input.Model = new(models.CreateTenant)
	input.Model.VerificationKey = strings.Replace(uuid.NewV4().String(), "-", "", 4)
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
			result.AddFieldFailure("email", "E-mail is required.")
		} else {
			if emailResult := validate.Email(input.Model.Email); !emailResult.Ok {
				result.AddFieldFailure("email", emailResult.Messages...)
			}
		}

		if len(input.Model.Email) > 200 {
			result.AddFieldFailure("email", "E-mail must be less than 200 characters.")
		}

		if input.Model.Name == "" {
			result.AddFieldFailure("name", "Name is required.")
		}
		if len(input.Model.Name) > 50 {
			result.AddFieldFailure("name", "Name must be less than 50 characters.")
		}
	}

	if env.IsSingleHostMode() {
		input.Model.Subdomain = "default"
	}

	if input.Model.TenantName == "" {
		result.AddFieldFailure("tenantName", "Name is required.")
	}

	subdomainResult := validate.Subdomain(services.Tenants, input.Model.Subdomain)
	if !subdomainResult.Ok {
		result.AddFieldFailure("subdomain", subdomainResult.Messages...)
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
func (input *UpdateTenantSettings) IsAuthorized(user *models.User) bool {
	return user != nil && user.Role == models.RoleAdministrator
}

// Validate is current model is valid
func (input *UpdateTenantSettings) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

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
		if hostnameResult := validate.CNAME(services.Tenants, input.Model.CNAME); !hostnameResult.Ok {
			result.AddFieldFailure("cname", hostnameResult.Messages...)
		}
	}

	return result
}
