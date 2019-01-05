package stripe

// DiscountParams is the set of parameters that can be used when deleting a discount.
type DiscountParams struct {
	Params `form:"*"`
}

// Discount is the resource representing a Stripe discount.
// For more details see https://stripe.com/docs/api#discounts.
type Discount struct {
	Coupon       *Coupon `json:"coupon"`
	Customer     string  `json:"customer"`
	Deleted      bool    `json:"deleted"`
	End          int64   `json:"end"`
	Start        int64   `json:"start"`
	Subscription string  `json:"subscription"`
}
