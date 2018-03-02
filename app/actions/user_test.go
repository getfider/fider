package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/onsi/gomega"
)

func TestChangeUserRole_Unauthorized(t *testing.T) {
	RegisterTestingT(t)

	for _, user := range []*models.User{
		&models.User{ID: 1, Role: models.RoleVisitor},
		&models.User{ID: 1, Role: models.RoleCollaborator},
		&models.User{ID: 2, Role: models.RoleAdministrator},
	} {
		action := actions.ChangeUserRole{Model: &models.ChangeUserRole{UserID: 2}}
		Expect(action.IsAuthorized(user)).To(BeFalse())
	}
}

func TestChangeUserRole_Authorized(t *testing.T) {
	RegisterTestingT(t)

	user := &models.User{ID: 2, Role: models.RoleAdministrator}
	action := actions.ChangeUserRole{Model: &models.ChangeUserRole{UserID: 1}}
	Expect(action.IsAuthorized(user)).To(BeTrue())
}

func TestChangeUserRole_InvalidRole(t *testing.T) {
	RegisterTestingT(t)

	tenant := &models.Tenant{ID: 1}
	services.SetCurrentTenant(tenant)

	targetUser := &models.User{
		Tenant: tenant,
	}
	services.Users.Register(targetUser)

	currentUser := &models.User{
		Tenant: tenant,
		Role:   models.RoleAdministrator,
	}
	services.Users.Register(currentUser)

	action := actions.ChangeUserRole{Model: &models.ChangeUserRole{UserID: targetUser.ID, Role: 4}}
	action.IsAuthorized(currentUser)
	result := action.Validate(currentUser, services)
	ExpectFailed(result, "role")
}

func TestChangeUserRole_InvalidUser(t *testing.T) {
	RegisterTestingT(t)

	currentUser := &models.User{
		Tenant: &models.Tenant{ID: 1},
		Role:   models.RoleAdministrator,
	}

	action := actions.ChangeUserRole{Model: &models.ChangeUserRole{UserID: 999, Role: models.RoleAdministrator}}
	action.IsAuthorized(currentUser)
	result := action.Validate(currentUser, services)
	ExpectFailed(result, "user_id")
}

func TestChangeUserRole_InvalidUser_Tenant(t *testing.T) {
	RegisterTestingT(t)

	targetUser := &models.User{
		Tenant: &models.Tenant{ID: 1},
	}
	services.Users.Register(targetUser)

	currentUser := &models.User{
		Tenant: &models.Tenant{ID: 2},
		Role:   models.RoleAdministrator,
	}
	services.Users.Register(currentUser)

	action := actions.ChangeUserRole{Model: &models.ChangeUserRole{UserID: targetUser.ID, Role: models.RoleAdministrator}}
	action.IsAuthorized(currentUser)
	result := action.Validate(currentUser, services)
	ExpectFailed(result, "user_id")
}
