package actions

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/uuid"
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

//ChangeUserEmail is the action used to change current user's e-mail
type ChangeUserEmail struct {
	Model *models.ChangeUserEmail
	User  *models.User
}

// Initialize the model
func (input *ChangeUserEmail) Initialize() interface{} {
	input.Model = new(models.ChangeUserEmail)
	input.Model.VerificationKey = strings.Replace(uuid.NewV4().String(), "-", "", 4)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *ChangeUserEmail) IsAuthorized(user *models.User) bool {
	input.User = user
	return user != nil
}

// Validate is current model is valid
func (input *ChangeUserEmail) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Email == "" {
		result.AddFieldFailure("email", "E-mail is required.")
		return result
	}

	if len(input.Model.Email) > 200 {
		result.AddFieldFailure("email", "E-mail must be less than 200 characters.")
		return result
	}

	if input.User.Email == input.Model.Email {
		result.AddFieldFailure("email", "Choose a different e-mail.")
		return result
	}

	emailResult := validate.Email(input.Model.Email)
	if !emailResult.Ok {
		result.AddFieldFailure("email", emailResult.Messages...)
		return result
	}

	existing, err := services.Users.GetByEmail(input.User.Tenant.ID, input.Model.Email)
	if err != nil && err != app.ErrNotFound {
		return validate.Error(err)
	}
	if err == nil && existing.ID != input.User.ID {
		result.AddFieldFailure("email", "This e-mail is already in use by someone else")
		return result
	}

	return result
}
