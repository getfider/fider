package actions

import (
	"context"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/validate"
)

// Actionable is any action that the user can perform using the web app
type Actionable interface {
	IsAuthorized(ctx context.Context, user *entity.User) bool
	Validate(ctx context.Context, user *entity.User) *validate.Result
}

// PreExecuteAction can add custom pre processing logic for any action
// OnPreExecute is executed before IsAuthorized and Validate
type PreExecuteAction interface {
	OnPreExecute(ctx context.Context) error
}
