package actions

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/env"

	"github.com/getfider/fider/app/pkg/validate"
)

// GenerateCheckoutLink is used to generate a Paddle-hosted checkout link for the service subscription
type GenerateCheckoutLink struct {
	PlanID string `json:"planId"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *GenerateCheckoutLink) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user.IsAdministrator()
}

// Validate if current model is valid
func (action *GenerateCheckoutLink) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if !env.IsBillingEnabled() {
		result.AddFieldFailure("plan_id", "Billing is not enabled.")
	} else if action.PlanID != env.Config.Paddle.MonthlyPlanID && action.PlanID != env.Config.Paddle.YearlyPlanID {
		result.AddFieldFailure("plan_id", "Invalid Plan ID.")
	}

	return result
}
