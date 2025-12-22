package actions

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/validate"
)

// AddVoteOnBehalf is used by administrators to add votes on behalf of users
type AddVoteOnBehalf struct {
	Number int    `route:"number"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *AddVoteOnBehalf) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *AddVoteOnBehalf) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	} else if len(action.Name) > 100 {
		result.AddFieldFailure("name", "Name must have less than 100 characters.")
	}

	if action.Email == "" {
		result.AddFieldFailure("email", "Email is required.")
	} else {
		messages := validate.Email(ctx, action.Email)
		if len(messages) > 0 {
			result.AddFieldFailure("email", messages...)
		}
	}

	return result
}

// RemoveVoteOnBehalf is used by administrators to remove votes on behalf of users
type RemoveVoteOnBehalf struct {
	Number int `route:"number"`
	UserID int `route:"userID"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *RemoveVoteOnBehalf) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *RemoveVoteOnBehalf) Validate(ctx context.Context, user *entity.User) *validate.Result {
	return validate.Success()
}
