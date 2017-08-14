package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestCreateTenant_EmptyToken(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: ""}}
	result := action.Validate(services)
	ExpectFailed(result, "token")
}

func TestCreateTenant_EmptyName(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: jonSnowToken, Name: ""}}
	result := action.Validate(services)
	ExpectFailed(result, "name")
}

func TestCreateTenant_EmptySubdomain(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: jonSnowToken, Name: "My Company"}}
	result := action.Validate(services)
	ExpectFailed(result, "subdomain")
}

func TestCreateTenant_UpperCaseSubdomain(t *testing.T) {
	RegisterTestingT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: jonSnowToken, Name: "My Company", Subdomain: "MyCompany"}}
	result := action.Validate(services)
	ExpectSuccess(result)
	Expect(action.Model.Subdomain).To(Equal("mycompany"))
}

func TestUpdateTenantSettings_EmptyTitle(t *testing.T) {
	RegisterTestingT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{}}
	result := action.Validate(services)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_LargeTitle(t *testing.T) {
	RegisterTestingT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "123456789012345678901234567890123456789012345678901234567890123"}}
	result := action.Validate(services)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_LargeInvitation(t *testing.T) {
	RegisterTestingT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "Ok", Invitation: "123456789012345678901234567890123456789012345678901234567890123"}}
	result := action.Validate(services)
	ExpectFailed(result, "invitation")
}
