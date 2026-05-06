package actions

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/validate"
)

// BulkDeleteDisposable is the input model for bulk-deleting disposable users.
type BulkDeleteDisposable struct {
	UserIDs []int `json:"userIds"`
}

// IsAuthorized returns true if the current user is an administrator.
func (a *BulkDeleteDisposable) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate checks that the request is well-formed.
func (a *BulkDeleteDisposable) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()
	if len(a.UserIDs) > 500 {
		result.AddFieldFailure("userIds", "Too many users in a single request (max 500).")
	}
	return result
}
