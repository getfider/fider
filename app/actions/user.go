package actions

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
)

//CreateUser is the action to create a new user
type CreateUser struct {
	Model *models.CreateUser
}

// Returns the struct to bind the request to
func (action *CreateUser) BindTarget() interface{} {
	action.Model = new(models.CreateUser)
	return action.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateUser) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *CreateUser) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	} else if len(action.Model.Name) > 100 {
		result.AddFieldFailure("name", "Name must have less than 100 characters.")
	}

	if action.Model.Email == "" && action.Model.Reference == "" {
		result.AddFieldFailure("", "Either email or reference is required")
	} else {
		if action.Model.Email != "" {
			messages := validate.Email(action.Model.Email)
			if len(messages) > 0 {
				result.AddFieldFailure("email", messages...)
			}
		}

		if len(action.Model.Reference) > 100 {
			result.AddFieldFailure("reference", "Reference must have less than 100 characters.")
		}
	}

	return result
}

//ChangeUserRole is the input model change role of an user
type ChangeUserRole struct {
	Model *models.ChangeUserRole
}

// Returns the struct to bind the request to
func (action *ChangeUserRole) BindTarget() interface{} {
	action.Model = new(models.ChangeUserRole)
	return action.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *ChangeUserRole) IsAuthorized(ctx context.Context, user *models.User) bool {
	if user == nil {
		return false
	}
	return user.IsAdministrator() && user.ID != action.Model.UserID
}

// Validate if current model is valid
func (action *ChangeUserRole) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()
	if action.Model.Role < enum.RoleVisitor || action.Model.Role > enum.RoleAdministrator {
		return validate.Error(app.ErrNotFound)
	}

	if user.ID == action.Model.UserID {
		result.AddFieldFailure("userID", "It is not allowed to change your own Role.")
	}

	userByID := &query.GetUserByID{UserID: action.Model.UserID}
	err := bus.Dispatch(ctx, userByID)
	if err != nil {
		if errors.Cause(err) == app.ErrNotFound {
			result.AddFieldFailure("userID", "User not found.")
		} else {
			return validate.Error(err)
		}
	} else if userByID.Result.Tenant.ID != user.Tenant.ID {
		result.AddFieldFailure("userID", "User not found.")
	}
	return result
}

//ChangeUserEmail is the action used to change current user's email
type ChangeUserEmail struct {
	Model *models.ChangeUserEmail
}

func NewChangeUserEmail() *ChangeUserEmail {
	return &ChangeUserEmail{
		Model: &models.ChangeUserEmail{
			VerificationKey: models.GenerateSecretKey(),
		},
	}
}

// Returns the struct to bind the request to
func (action *ChangeUserEmail) BindTarget() interface{} {
	return action.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *ChangeUserEmail) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil
}

// Validate if current model is valid
func (action *ChangeUserEmail) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Model.Email == "" {
		result.AddFieldFailure("email", "Email is required.")
		return result
	}

	if len(action.Model.Email) > 200 {
		result.AddFieldFailure("email", "Email must have less than 200 characters.")
		return result
	}

	if user.Email == action.Model.Email {
		result.AddFieldFailure("email", "Choose a different email.")
		return result
	}

	messages := validate.Email(action.Model.Email)
	if len(messages) > 0 {
		result.AddFieldFailure("email", messages...)
		return result
	}

	userByEmail := &query.GetUserByEmail{Email: action.Model.Email}
	err := bus.Dispatch(ctx, userByEmail)
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return validate.Error(err)
	}
	if err == nil && userByEmail.Result.ID != user.ID {
		result.AddFieldFailure("email", "This email is already in use by someone else")
		return result
	}
	action.Model.Requestor = user
	return result
}
