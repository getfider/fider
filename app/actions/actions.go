package actions

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/validate"
)

// Actionable is any action that the user can perform using the web app
type Actionable interface {
	IsAuthorized(ctx context.Context, user *entity.User) bool
	Validate(ctx context.Context, user *entity.User) *validate.Result
}
