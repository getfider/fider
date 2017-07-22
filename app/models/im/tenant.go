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
	var err error
	if i.Token == "" {
		return validate.Failed([]string{"Please identify yourself before proceeding."})
	}

	if i.UserClaims, err = jwt.DecodeOAuthClaims(i.Token); err != nil {
		return validate.Error(err)
	}

	if env.IsSingleHostMode() {
		i.Subdomain = "default"
	}

	if i.Name == "" {
		return validate.Failed([]string{"Name is required."})
	}

	return validate.Subdomain(services.Tenants, i.Subdomain)
}
