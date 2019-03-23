package dto

// BillingPlan is the model for billing plan from Stripe
type BillingPlan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Currency    string `json:"currency"`
	MaxUsers    int    `json:"maxUsers"`
	Price       int64  `json:"price"`
	Interval    string `json:"interval"`
}
