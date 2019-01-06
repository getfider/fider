package stripe

// ThreeDSecureStatus represents the possible statuses of a ThreeDSecure object.
type ThreeDSecureStatus string

// ThreeDSecureParams is the set of parameters that can be used when creating a 3DS object.
type ThreeDSecureParams struct {
	Params    `form:"*"`
	Amount    *int64  `form:"amount"`
	Card      *string `form:"card"`
	Currency  *string `form:"currency"`
	Customer  *string `form:"customer"`
	ReturnURL *string `form:"return_url"`
}

// ThreeDSecure is the resource representing a Stripe 3DS object
// For more details see https://stripe.com/docs/api#three_d_secure.
type ThreeDSecure struct {
	Amount        int64              `json:"amount"`
	Authenticated bool               `json:"authenticated"`
	Card          *Card              `json:"card"`
	Created       int64              `json:"created"`
	Currency      Currency           `json:"currency"`
	ID            string             `json:"id"`
	Livemode      bool               `json:"livemode"`
	RedirectURL   string             `json:"redirect_url"`
	Status        ThreeDSecureStatus `json:"status"`
	Supported     string             `json:"supported"`
}
