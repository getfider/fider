package actions

import (
	"context"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

// Actionable is any action that the user can perform using the web app
type Actionable interface {
	IsAuthorized(ctx context.Context, user *models.User) bool
	Validate(ctx context.Context, user *models.User) *validate.Result
}

// BindTargetAction defines an action where the bindable model is not itself
type BindTargetAction interface {
	BindTarget() interface{}
}
