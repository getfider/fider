package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestCreateTenant_ShouldHaveVerificationKey(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{}
	action.Initialize()

	Expect(action.Model.VerificationKey).IsNotEmpty()
}

func TestCreateTenant_EmptyToken(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Model: &models.CreateTenant{
			Token:          "",
			LegalAgreement: true,
		},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "token", "tenantName", "subdomain")
}

func TestCreateTenant_EmptyTenantName(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Model: &models.CreateTenant{
			Token:          jonSnowToken,
			TenantName:     "",
			LegalAgreement: true,
		},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "tenantName", "subdomain")
}

func TestCreateTenant_EmptyEmail(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Model: &models.CreateTenant{
			Name:           "Jon Snow",
			Email:          "",
			LegalAgreement: true,
		},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "email", "tenantName", "subdomain")
}

func TestCreateTenant_InvalidEmail(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Model: &models.CreateTenant{
			Name:           "Jon Snow",
			Email:          "jonsnow",
			LegalAgreement: true,
		},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "email", "tenantName", "subdomain")
}

func TestCreateTenant_NoAgreement(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.IsSubdomainAvailable) error {
		q.Result = true
		return nil
	})

	action := actions.CreateTenant{
		Model: &models.CreateTenant{
			Name:           "Jon",
			Email:          "jon.snow@got.com",
			TenantName:     "My Company",
			Subdomain:      "company",
			LegalAgreement: false,
		},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "legalAgreement")
}

func TestCreateTenant_EmptyName(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Model: &models.CreateTenant{
			Name:           "",
			Email:          "jon.snow@got.com",
			LegalAgreement: true,
		},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "name", "tenantName", "subdomain")
}

func TestCreateTenant_EmptySubdomain(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Model: &models.CreateTenant{
			Token:          jonSnowToken,
			TenantName:     "My Company",
			LegalAgreement: true,
		},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "subdomain")
}

func TestUpdateTenantSettings_Unauthorized(t *testing.T) {
	RegisterT(t)

	admin := &models.User{ID: 1, Role: enum.RoleAdministrator}
	collaborator := &models.User{ID: 2, Role: enum.RoleCollaborator}

	action := actions.UpdateTenantSettings{}
	action.Initialize()

	Expect(action.IsAuthorized(context.Background(), admin)).IsTrue()
	Expect(action.IsAuthorized(context.Background(), collaborator)).IsFalse()
	Expect(action.IsAuthorized(context.Background(), nil)).IsFalse()
}

func TestUpdateTenantSettings_EmptyTitle(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{}
	action.Initialize()
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_InvalidCNAME(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "Ok", CNAME: "bla"}}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "cname")
}

func TestUpdateTenantSettings_LargeTitle(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "123456789012345678901234567890123456789012345678901234567890123"}}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_LargeInvitation(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "Ok", Invitation: "123456789012345678901234567890123456789012345678901234567890123"}}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "invitation")
}

func TestUpdateTenantSettings_ExistingTenant_WithLogo(t *testing.T) {
	RegisterT(t)

	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &models.Tenant{
		ID:          1,
		LogoBlobKey: "hello-world.png",
	})

	action := actions.UpdateTenantSettings{}
	action.Initialize()
	action.Model.Title = "OK"
	action.Model.Invitation = "Share your ideas!"
	result := action.Validate(ctx, nil)
	ExpectSuccess(result)
	Expect(action.Model.Logo.BlobKey).Equals("hello-world.png")
}
