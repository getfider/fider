package billing

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/goenning/vat"
	"github.com/stripe/stripe-go"
)

var plansMutex sync.RWMutex
var allPlans []*dto.BillingPlan

func listPlans(ctx context.Context, q *query.ListBillingPlans) error {
	if allPlans != nil {
		q.Result = filterPlansByCountryCode(allPlans, q.CountryCode)
		return nil
	}

	plansMutex.Lock()
	defer plansMutex.Unlock()

	if allPlans == nil {
		allPlans = make([]*dto.BillingPlan, 0)
		it := stripeClient.Plans.List(&stripe.PlanListParams{
			Active: stripe.Bool(true),
		})
		for it.Next() {
			plan := it.Plan()
			name, ok := plan.Metadata["friendly_name"]
			if !ok {
				name = plan.Nickname
			}
			maxUsers, _ := strconv.Atoi(plan.Metadata["max_users"])
			allPlans = append(allPlans, &dto.BillingPlan{
				ID:          plan.ID,
				Name:        name,
				Description: plan.Metadata["description"],
				MaxUsers:    maxUsers,
				Price:       plan.Amount,
				Currency:    strings.ToUpper(string(plan.Currency)),
				Interval:    string(plan.Interval),
			})
		}
		if err := it.Err(); err != nil {
			return err
		}
		sort.Slice(allPlans, func(i, j int) bool {
			return allPlans[i].Price < allPlans[j].Price
		})
	}

	q.Result = filterPlansByCountryCode(allPlans, q.CountryCode)
	return nil
}

func getPlanByID(ctx context.Context, q *query.GetBillingPlanByID) error {
	listPlansQuery := &query.ListBillingPlans{CountryCode: q.CountryCode}
	err := listPlans(ctx, listPlansQuery)
	if err != nil {
		return err
	}

	for _, plan := range listPlansQuery.Result {
		if plan.ID == q.PlanID {
			q.Result = plan
			return nil
		}
	}
	return errors.New("failed to get plan by id '%s'", q.PlanID)
}

func filterPlansByCountryCode(plans []*dto.BillingPlan, countryCode string) []*dto.BillingPlan {
	currency := "USD"
	if vat.IsEU(countryCode) {
		currency = "EUR"
	}

	filteredPlans := make([]*dto.BillingPlan, 0)
	for _, p := range plans {
		if p.Currency == currency {
			filteredPlans = append(filteredPlans, p)
		}
	}
	return filteredPlans
}
