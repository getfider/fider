package enum

type BillingStatus int

var (
	//BillingTrial is used for tenants in trial
	BillingTrial BillingStatus = 1
	//BillingActive is used for tenants with an active subscription
	BillingActive BillingStatus = 2
	//BillingCancelled is used for tenants that had an active subscription, but have cancelled it
	BillingCancelled BillingStatus = 3
	//BillingFreeForever is used for tenants that are on a forever free plan
	BillingFreeForever BillingStatus = 4
	//BillingOpenCollective is used for tenants that have an active open collective subsription
	BillingOpenCollective BillingStatus = 5
)

var billingStatusIDs = map[BillingStatus]string{
	BillingTrial:          "trial",
	BillingActive:         "active",
	BillingCancelled:      "cancelled",
	BillingFreeForever:    "free_forever",
	BillingOpenCollective: "open_collective",
}

// String returns the string version of the billing status
func (status BillingStatus) String() string {
	return billingStatusIDs[status]
}
