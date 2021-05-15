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
	Model *models.CreateTenant
}

func NewCreateTenant() *CreateTenant {
	return &CreateTenant{
		Model: &models.CreateTenant{
			VerificationKey: models.GenerateSecretKey(),
		},
	}
}

// Returns the struct to bind the request to
func (action *CreateTenant) BindTarget() interface{} {
	return action.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateTenant) IsAuthorized(ctx context.Context, user *models.User) bool {
	return true
}

// Validate if current model is valid
func (action *CreateTenant) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	var err error
	if action.Model.Name == "" && action.Model.Email == "" {
		if action.Model.Token == "" {
			result.AddFieldFailure("token", "Please identify yourself before proceeding.")
		} else {
			if action.Model.UserClaims, err = jwt.DecodeOAuthClaims(action.Model.Token); err != nil {
				return validate.Error(err)
			}
		}
	} else {
		if action.Model.Email == "" {
			result.AddFieldFailure("email", "Email is required.")
		} else {
			messages := validate.Email(action.Model.Email)
			result.AddFieldFailure("email", messages...)
		}

		if action.Model.Name == "" {
			result.AddFieldFailure("name", "Name is required.")
		}
		if len(action.Model.Name) > 60 {
			result.AddFieldFailure("name", "Name must have less than 60 characters.")
		}
	}

	if env.IsSingleHostMode() {
		action.Model.Subdomain = "default"
	}

	if action.Model.TenantName == "" {
		result.AddFieldFailure("tenantName", "Name is required.")
	}

	messages, err := validate.Subdomain(ctx, action.Model.Subdomain)
	if err != nil {
		return validate.Error(err)
	}

	result.AddFieldFailure("subdomain", messages...)

	if env.HasLegal() && !action.Model.LegalAgreement {
		result.AddFieldFailure("legalAgreement", "You must agree before proceeding.")
	}

	return result
}

//UpdateTenantSettings is the input model used to update tenant settings
type UpdateTenantSettings struct {
	Model *models.UpdateTenantSettings
}

func NewUpdateTenantSettings() *UpdateTenantSettings {
	return &UpdateTenantSettings{
		Model: &models.UpdateTenantSettings{
			Logo: &models.ImageUpload{},
		},
	}
}

// Returns the struct to bind the request to
func (action *UpdateTenantSettings) BindTarget() interface{} {
	return action.Model
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
		action.Model.Logo.BlobKey = tenant.LogoBlobKey
	}

	messages, err := validate.ImageUpload(action.Model.Logo, validate.ImageUploadOpts{
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

	if action.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	}

	if len(action.Model.Title) > 60 {
		result.AddFieldFailure("title", "Title must have less than 60 characters.")
	}

	if len(action.Model.Invitation) > 60 {
		result.AddFieldFailure("invitation", "Invitation must have less than 60 characters.")
	}

	if action.Model.CNAME != "" {
		messages := validate.CNAME(ctx, action.Model.CNAME)
		result.AddFieldFailure("cname", messages...)
	}

	return result
}

//UpdateTenantAdvancedSettings is the input model used to update tenant advanced settings
type UpdateTenantAdvancedSettings struct {
	Model *models.UpdateTenantAdvancedSettings
}

// Returns the struct to bind the request to
func (action *UpdateTenantAdvancedSettings) BindTarget() interface{} {
	action.Model = new(models.UpdateTenantAdvancedSettings)
	return action.Model
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
	Model *models.UpdateTenantPrivacy
}

// Returns the struct to bind the request to
func (action *UpdateTenantPrivacy) BindTarget() interface{} {
	action.Model = new(models.UpdateTenantPrivacy)
	return action.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *UpdateTenantPrivacy) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.Role == enum.RoleAdministrator
}

// Validate if current model is valid
func (action *UpdateTenantPrivacy) Validate(ctx context.Context, user *models.User) *validate.Result {
	return validate.Success()
}
