package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestCreateTenant_ShouldHaveVerificationKey(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{}
	action.Initialize()

	Expect(action.Model.VerificationKey).NotTo(Equal(""))
}

func TestCreateTenant_EmptyToken(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: ""}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "token", "tenantName", "subdomain")
}

func TestCreateTenant_EmptyTenantName(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: jonSnowToken, TenantName: ""}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "tenantName", "subdomain")
}

func TestCreateTenant_EmptyEmail(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Name: "Jon Snow", Email: ""}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "email", "tenantName", "subdomain")
}

func TestCreateTenant_InvalidEmail(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Name: "Jon Snow", Email: "jonsnow"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "email", "tenantName", "subdomain")
}

func TestCreateTenant_EmptyName(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Name: "", Email: "jon.snow@got.com"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "name", "tenantName", "subdomain")
}

func TestCreateTenant_EmptySubdomain(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: jonSnowToken, TenantName: "My Company"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "subdomain")
}
func TestUpdateTenantSettings_Unauthorized(t *testing.T) {
	RegisterTestingT(t)

	admin := &models.User{ID: 1, Role: models.RoleAdministrator}
	collaborator := &models.User{ID: 2, Role: models.RoleCollaborator}

	action := actions.UpdateTenantSettings{}
	action.Initialize()

	Expect(action.IsAuthorized(admin, nil)).To(BeTrue())
	Expect(action.IsAuthorized(collaborator, nil)).To(BeFalse())
	Expect(action.IsAuthorized(nil, nil)).To(BeFalse())
}

func TestUpdateTenantSettings_EmptyTitle(t *testing.T) {
	RegisterTestingT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_InvalidCNAME(t *testing.T) {
	RegisterTestingT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "Ok", CNAME: "bla"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "cname")
}

func TestUpdateTenantSettings_LargeTitle(t *testing.T) {
	RegisterTestingT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "123456789012345678901234567890123456789012345678901234567890123"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_LargeInvitation(t *testing.T) {
	RegisterTestingT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "Ok", Invitation: "123456789012345678901234567890123456789012345678901234567890123"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "invitation")
}
