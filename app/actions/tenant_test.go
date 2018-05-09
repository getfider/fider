package actions_test

import (
	"io/ioutil"
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
)

func TestCreateTenant_ShouldHaveVerificationKey(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{}
	action.Initialize()

	Expect(action.Model.VerificationKey).IsNotEmpty()
}

func TestCreateTenant_EmptyToken(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: ""}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "token", "tenantName", "subdomain")
}

func TestCreateTenant_EmptyTenantName(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: jonSnowToken, TenantName: ""}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "tenantName", "subdomain")
}

func TestCreateTenant_EmptyEmail(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Name: "Jon Snow", Email: ""}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "email", "tenantName", "subdomain")
}

func TestCreateTenant_InvalidEmail(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Name: "Jon Snow", Email: "jonsnow"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "email", "tenantName", "subdomain")
}

func TestCreateTenant_EmptyName(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Name: "", Email: "jon.snow@got.com"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "name", "tenantName", "subdomain")
}

func TestCreateTenant_EmptySubdomain(t *testing.T) {
	RegisterT(t)

	action := actions.CreateTenant{Model: &models.CreateTenant{Token: jonSnowToken, TenantName: "My Company"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "subdomain")
}
func TestUpdateTenantSettings_Unauthorized(t *testing.T) {
	RegisterT(t)

	admin := &models.User{ID: 1, Role: models.RoleAdministrator}
	collaborator := &models.User{ID: 2, Role: models.RoleCollaborator}

	action := actions.UpdateTenantSettings{}
	action.Initialize()

	Expect(action.IsAuthorized(admin, nil)).IsTrue()
	Expect(action.IsAuthorized(collaborator, nil)).IsFalse()
	Expect(action.IsAuthorized(nil, nil)).IsFalse()
}

func TestUpdateTenantSettings_EmptyTitle(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_InvalidCNAME(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "Ok", CNAME: "bla"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "cname")
}

func TestUpdateTenantSettings_LargeTitle(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "123456789012345678901234567890123456789012345678901234567890123"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "title")
}

func TestUpdateTenantSettings_LargeInvitation(t *testing.T) {
	RegisterT(t)

	action := actions.UpdateTenantSettings{Model: &models.UpdateTenantSettings{Title: "Ok", Invitation: "123456789012345678901234567890123456789012345678901234567890123"}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "invitation")
}

func TestUpdateTenantSettings_LargeLogo(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		fileName string
		valid    bool
	}{
		{"/app/pkg/img/testdata/logo1.png", true},
		{"/app/pkg/img/testdata/logo2.jpg", false},
		{"/app/pkg/img/testdata/logo3.gif", false},
		{"/app/pkg/img/testdata/logo4.png", false},
		{"/app/pkg/img/testdata/logo5.png", true},
		{"/README.md", false},
		{"/favicon.ico", false},
	}

	for _, testCase := range testCases {
		logo, _ := ioutil.ReadFile(env.Path(testCase.fileName))

		action := actions.UpdateTenantSettings{
			Model: &models.UpdateTenantSettings{
				Title: "Hello World",
				Logo: &models.UpdateTenantSettingsLogo{
					Upload: &models.UpdateTenantSettingsLogoUpload{
						Content: logo,
					},
				},
			},
		}
		result := action.Validate(nil, services)
		if testCase.valid {
			ExpectSuccess(result)
		} else {
			ExpectFailed(result, "logo")
		}
	}
}
