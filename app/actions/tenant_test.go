package actions_test

import (
	"context"
	"testing"

	"github.com/Spicy-Bush/fider-tarkov-community/app"

	"github.com/Spicy-Bush/fider-tarkov-community/app/actions"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/query"

	. "github.com/Spicy-Bush/fider-tarkov-community/app/pkg/assert"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
)

func TestCreateTenant_ShouldHaveVerificationKey(t *testing.T) {
	RegisterT(t)

	action := actions.NewCreateTenant()
	Expect(action.VerificationKey).IsNotEmpty()
}

func TestCreateTenant_EmptyToken(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Token:          "",
		LegalAgreement: true,
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "token", "tenantName", "subdomain")
}

func TestCreateTenant_EmptyTenantName(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Token:          jonSnowToken,
		TenantName:     "",
		LegalAgreement: true,
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "tenantName", "subdomain")
}

func TestCreateTenant_EmptyEmail(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Name:           "Jon Snow",
		Email:          "",
		LegalAgreement: true,
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "email", "tenantName", "subdomain")
}

func TestCreateTenant_InvalidEmail(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Name:           "Jon Snow",
		Email:          "jonsnow",
		LegalAgreement: true,
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
		Name:           "Jon",
		Email:          "jon.snow@got.com",
		TenantName:     "My Company",
		Subdomain:      "company",
		LegalAgreement: false,
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "legalAgreement")
}

func TestCreateTenant_EmptyName(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Name:           "",
		Email:          "jon.snow@got.com",
		LegalAgreement: true,
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "name", "tenantName", "subdomain")
}

func TestCreateTenant_EmptySubdomain(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{
		Token:          jonSnowToken,
		TenantName:     "My Company",
		LegalAgreement: true,
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "subdomain")
}

func TestUpdateTenantSettings_Unauthorized(t *testing.T) {
	RegisterT(t)

	admin := &entity.User{ID: 1, Role: enum.RoleAdministrator}
	collaborator := &entity.User{ID: 2, Role: enum.RoleCollaborator}

	action := actions.NewUpdateTenantSettings()

	Expect(action.IsAuthorized(context.Background(), admin)).IsTrue()
	Expect(action.IsAuthorized(context.Background(), collaborator)).IsFalse()
	Expect(action.IsAuthorized(context.Background(), nil)).IsFalse()
}

func TestUpdateTenantSettings_EmptyTitle(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Locale: "en"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_InvalidCNAME(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Title: "Ok", CNAME: "bla", Locale: "en"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "cname")
}

func TestUpdateTenantSettings_LargeTitle(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Title: "123456789012345678901234567890123456789012345678901234567890123", Locale: "en"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_LargeInvitation(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Title: "Ok", Invitation: "123456789012345678901234567890123456789012345678901234567890123", Locale: "en"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "invitation")
}

func TestUpdateTenantSettings_InvalidLocale(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Title: "Some Name", Locale: "xx"}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "locale")
}

func TestUpdateTenantSettings_ExistingTenant_WithLogo(t *testing.T) {
	RegisterT(t)

	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{
		ID:          1,
		LogoBlobKey: "hello-world.png",
	})

	action := actions.NewUpdateTenantSettings()
	action.Title = "OK"
	action.Invitation = "Share your ideas!"
	action.Locale = "en"
	result := action.Validate(ctx, nil)
	ExpectSuccess(result)
	Expect(action.Logo.BlobKey).Equals("hello-world.png")
}
