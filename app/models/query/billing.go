package query

import "github.com/getfider/fider/app/models/dto"

type ListBillingPlans struct {
	CountryCode string

	Plans []*dto.BillingPlan
}

type GetBillingPlanByID struct {
	PlanID      string
	CountryCode string

	Plan *dto.BillingPlan
}
