package actions

import (
	"context"

	"github.com/Spicy-Bush/fider-tarkov-community/app"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/cmd"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/errors"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/validate"
)

type UpdateProfanityWords struct {
	ProfanityWords string `json:"profanityWords"`
}

func NewUpdateProfanityWords() *UpdateProfanityWords {
	return &UpdateProfanityWords{}
}

func (action *UpdateProfanityWords) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

func (action *UpdateProfanityWords) Validate(ctx context.Context, user *entity.User) *validate.Result {
	return validate.Success()
}

func (action *UpdateProfanityWords) Run(ctx context.Context) error {
	tenant, ok := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	if !ok {
		return errors.New("tenant not found in context")
	}
	return bus.Dispatch(ctx, &cmd.UpdateTenantAdvancedSettings{
		CustomCSS:      tenant.CustomCSS, // preserve existing custom CSS
		ProfanityWords: action.ProfanityWords,
	})
}
