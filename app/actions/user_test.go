package actions_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestChangeUserRole_Unauthorized(t *testing.T) {
	RegisterT(t)

	for _, user := range []*models.User{
		&models.User{ID: 1, Role: models.RoleVisitor},
		&models.User{ID: 1, Role: models.RoleCollaborator},
		&models.User{ID: 2, Role: models.RoleAdministrator},
	} {
		action := actions.ChangeUserRole{Model: &models.ChangeUserRole{UserID: 2}}
		Expect(action.IsAuthorized(user, nil)).IsFalse()
	}
}

func TestChangeUserRole_Authorized(t *testing.T) {
	RegisterT(t)

	user := &models.User{ID: 2, Role: models.RoleAdministrator}
	action := actions.ChangeUserRole{Model: &models.ChangeUserRole{UserID: 1}}
	Expect(action.IsAuthorized(user, nil)).IsTrue()
}

func TestChangeUserRole_InvalidRole(t *testing.T) {
	RegisterT(t)

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
	action.IsAuthorized(currentUser, nil)
	result := action.Validate(currentUser, services)
	Expect(result.Err).Equals(app.ErrNotFound)
}

func TestChangeUserRole_InvalidUser(t *testing.T) {
	RegisterT(t)

	currentUser := &models.User{
		Tenant: &models.Tenant{ID: 1},
		Role:   models.RoleAdministrator,
	}

	action := actions.ChangeUserRole{Model: &models.ChangeUserRole{UserID: 999, Role: models.RoleAdministrator}}
	action.IsAuthorized(currentUser, nil)
	result := action.Validate(currentUser, services)
	ExpectFailed(result, "userId")
}

func TestChangeUserRole_InvalidUser_Tenant(t *testing.T) {
	RegisterT(t)

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
	action.IsAuthorized(currentUser, nil)
	result := action.Validate(currentUser, services)
	ExpectFailed(result, "userId")
}

func TestChangeUserRole_CurrentUser(t *testing.T) {
	RegisterT(t)

	currentUser := &models.User{
		Tenant: &models.Tenant{ID: 2},
		Role:   models.RoleAdministrator,
	}
	services.Users.Register(currentUser)

	action := actions.ChangeUserRole{Model: &models.ChangeUserRole{UserID: currentUser.ID, Role: models.RoleVisitor}}
	action.IsAuthorized(currentUser, nil)
	result := action.Validate(currentUser, services)
	ExpectFailed(result, "userId")

	user, err := services.Users.GetByID(currentUser.ID)
	Expect(err).IsNil()
	Expect(user.Role).Equals(models.RoleAdministrator)
}
