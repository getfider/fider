package im

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/validate"
)

//CreateTenant is the input model used to create a tenant
type CreateTenant struct {
	Token      string `json:"token"`
	Name       string `json:"name"`
	Subdomain  string `json:"subdomain"`
	UserClaims *models.OAuthClaims
}

// IsAuthorized returns true if current user is authorized to perform this action
func (i *CreateTenant) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (i *CreateTenant) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	var err error
	if i.Token == "" {
		result.AddFieldFailure("token", "Please identify yourself before proceeding.")
	} else {
		if i.UserClaims, err = jwt.DecodeOAuthClaims(i.Token); err != nil {
			return validate.Error(err)
		}
	}

	if env.IsSingleHostMode() {
		i.Subdomain = "default"
	}

	if i.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	subdomainResult := validate.Subdomain(services.Tenants, i.Subdomain)
	if !subdomainResult.Ok {
		result.AddFieldFailure("subdomain", subdomainResult.Messages...)
	}

	return result
}
