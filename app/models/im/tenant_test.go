package im_test

import (
	"testing"

	"github.com/getfider/fider/app/models/im"
	. "github.com/onsi/gomega"
)

func TestCreateTenant_EmptyToken(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: ""}
	result := input.Validate(services)
	ExpectFailed(result, "token")
}

func TestCreateTenant_EmptyName(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: jonSnowToken, Name: ""}
	result := input.Validate(services)
	ExpectFailed(result, "name")
}

func TestCreateTenant_EmptySubdomain(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: jonSnowToken, Name: "My Company"}
	result := input.Validate(services)
	ExpectFailed(result, "subdomain")
}

func TestCreateTenant_UpperCaseSubdomain(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: jonSnowToken, Name: "My Company", Subdomain: "MyCompany"}
	result := input.Validate(services)
	ExpectSuccess(result)
	Expect(input.Subdomain).To(Equal("mycompany"))
}

func TestUpdateTenantSettings_EmptyTitle(t *testing.T) {
	RegisterTestingT(t)

	input := im.UpdateTenantSettings{}
	result := input.Validate(services)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_LargeTitle(t *testing.T) {
	RegisterTestingT(t)

	input := im.UpdateTenantSettings{Title: "123456789012345678901234567890123456789012345678901234567890123"}
	result := input.Validate(services)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_LargeInvitation(t *testing.T) {
	RegisterTestingT(t)

	input := im.UpdateTenantSettings{Title: "Ok", Invitation: "123456789012345678901234567890123456789012345678901234567890123"}
	result := input.Validate(services)
	ExpectFailed(result, "invitation")
}
