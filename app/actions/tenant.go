package actions

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/validate"
)

//CreateTenant is the input model used to create a tenant
type CreateTenant struct {
	Input *models.CreateTenant
}

func NewCreateTenant() *CreateTenant {
	return &CreateTenant{
		Input: &models.CreateTenant{
			VerificationKey: models.GenerateSecretKey(),
		},
	}
}

// Returns the struct to bind the request to
func (action *CreateTenant) BindTarget() interface{} {
	return action.Input
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateTenant) IsAuthorized(ctx context.Context, user *models.User) bool {
	return true
}

// Validate if current model is valid
func (action *CreateTenant) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	var err error
	if action.Input.Name == "" && action.Input.Email == "" {
		if action.Input.Token == "" {
			result.AddFieldFailure("token", "Please identify yourself before proceeding.")
		} else {
			if action.Input.UserClaims, err = jwt.DecodeOAuthClaims(action.Input.Token); err != nil {
				return validate.Error(err)
			}
		}
	} else {
		if action.Input.Email == "" {
			result.AddFieldFailure("email", "Email is required.")
		} else {
			messages := validate.Email(action.Input.Email)
			result.AddFieldFailure("email", messages...)
		}

		if action.Input.Name == "" {
			result.AddFieldFailure("name", "Name is required.")
		}
		if len(action.Input.Name) > 60 {
			result.AddFieldFailure("name", "Name must have less than 60 characters.")
		}
	}

	if env.IsSingleHostMode() {
		action.Input.Subdomain = "default"
	}

	if action.Input.TenantName == "" {
		result.AddFieldFailure("tenantName", "Name is required.")
	}

	messages, err := validate.Subdomain(ctx, action.Input.Subdomain)
	if err != nil {
		return validate.Error(err)
	}

	result.AddFieldFailure("subdomain", messages...)

	if env.HasLegal() && !action.Input.LegalAgreement {
		result.AddFieldFailure("legalAgreement", "You must agree before proceeding.")
	}

	return result
}

//UpdateTenantSettings is the input model used to update tenant settings
type UpdateTenantSettings struct {
	Input *models.UpdateTenantSettings
}

func NewUpdateTenantSettings() *UpdateTenantSettings {
	return &UpdateTenantSettings{
		Input: &models.UpdateTenantSettings{
			Logo: &models.ImageUpload{},
		},
	}
}

// Returns the struct to bind the request to
func (action *UpdateTenantSettings) BindTarget() interface{} {
	return action.Input
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *UpdateTenantSettings) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.Role == enum.RoleAdministrator
}

// Validate if current model is valid
func (action *UpdateTenantSettings) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	tenant, hasTenant := ctx.Value(app.TenantCtxKey).(*models.Tenant)
	if hasTenant {
		action.Input.Logo.BlobKey = tenant.LogoBlobKey
	}

	messages, err := validate.ImageUpload(action.Input.Logo, validate.ImageUploadOpts{
		IsRequired:   false,
		MinHeight:    200,
		MinWidth:     200,
		MaxKilobytes: 100,
		ExactRatio:   true,
	})
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("logo", messages...)

	if action.Input.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	}

	if len(action.Input.Title) > 60 {
		result.AddFieldFailure("title", "Title must have less than 60 characters.")
	}

	if len(action.Input.Invitation) > 60 {
		result.AddFieldFailure("invitation", "Invitation must have less than 60 characters.")
	}

	if action.Input.CNAME != "" {
		messages := validate.CNAME(ctx, action.Input.CNAME)
		result.AddFieldFailure("cname", messages...)
	}

	return result
}

//UpdateTenantAdvancedSettings is the input model used to update tenant advanced settings
type UpdateTenantAdvancedSettings struct {
	Input *models.UpdateTenantAdvancedSettings
}

// Returns the struct to bind the request to
func (action *UpdateTenantAdvancedSettings) BindTarget() interface{} {
	action.Input = new(models.UpdateTenantAdvancedSettings)
	return action.Input
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *UpdateTenantAdvancedSettings) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.Role == enum.RoleAdministrator
}

// Validate if current model is valid
func (action *UpdateTenantAdvancedSettings) Validate(ctx context.Context, user *models.User) *validate.Result {
	return validate.Success()
}

//UpdateTenantPrivacy is the input model used to update tenant privacy settings
type UpdateTenantPrivacy struct {
	Input *models.UpdateTenantPrivacy
}

// Returns the struct to bind the request to
func (action *UpdateTenantPrivacy) BindTarget() interface{} {
	action.Input = new(models.UpdateTenantPrivacy)
	return action.Input
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *UpdateTenantPrivacy) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.Role == enum.RoleAdministrator
}

// Validate if current model is valid
func (action *UpdateTenantPrivacy) Validate(ctx context.Context, user *models.User) *validate.Result {
	return validate.Success()
}
