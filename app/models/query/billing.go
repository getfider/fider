package query

import "github.com/getfider/fider/app/models/dto"

type ListBillingPlans struct {
	CountryCode string

	Result []*dto.BillingPlan
}

type GetBillingPlanByID struct {
	PlanID      string
	CountryCode string

	Result *dto.BillingPlan
}

type GetUpcomingInvoice struct {
	Result *dto.UpcomingInvoice
}

type GetPaymentInfo struct {
	Result *dto.PaymentInfo
}

type GetAllCountries struct {
	Result []*dto.Country
}

type GetCountryByCode struct {
	Code string

	Result *dto.Country
}
