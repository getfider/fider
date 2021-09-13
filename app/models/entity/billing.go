package entity

import (
	"time"

	"github.com/getfider/fider/app/models/enum"
)

type BillingState struct {
	Status             enum.BillingStatus `json:"status"`
	PlanID             string             `json:"planID"`
	SubscriptionID     string             `json:"subscriptionID"`
	TrialEndsAt        *time.Time         `json:"trialEndsAt"`
	SubscriptionEndsAt *time.Time         `json:"subscriptionEndsAt"`
}
