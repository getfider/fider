package actions

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

//ChangeUserRole is the input model change role of an user
type ChangeUserRole struct {
	CurrentTenant *models.Tenant
	Model         *models.ChangeUserRole
}

// Initialize the model
func (input *ChangeUserRole) Initialize() interface{} {
	input.Model = new(models.ChangeUserRole)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *ChangeUserRole) IsAuthorized(user *models.User) bool {
	if user == nil {
		return false
	}
	input.CurrentTenant = user.Tenant
	return user.IsAdministrator() && user.ID != input.Model.UserID
}

// Validate is current model is valid
func (input *ChangeUserRole) Validate(services *app.Services) *validate.Result {
	result := validate.Success()
	if input.Model.Role < models.RoleVisitor || input.Model.Role > models.RoleAdministrator {
		result.AddFieldFailure("role", "Invalid role")
	}
	user, err := services.Users.GetByID(input.Model.UserID)
	if err != nil {
		if err == app.ErrNotFound {
			result.AddFieldFailure("user_id", "User not found")
		} else {
			return validate.Error(err)
		}
	} else if user.Tenant.ID != input.CurrentTenant.ID {
		result.AddFieldFailure("user_id", "User not found")
	}
	return result
}
