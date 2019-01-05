package stripe

import "encoding/json"

// IssuingCardShippingStatus is the list of possible values for the shipping status
// on an issuing card.
type IssuingCardShippingStatus string

// List of values that IssuingCardShippingStatus can take.
const (
	IssuingCardShippingTypeDelivered IssuingCardShippingStatus = "delivered"
	IssuingCardShippingTypeFailure   IssuingCardShippingStatus = "failure"
	IssuingCardShippingTypePending   IssuingCardShippingStatus = "pending"
	IssuingCardShippingTypeReturned  IssuingCardShippingStatus = "returned"
	IssuingCardShippingTypeShipped   IssuingCardShippingStatus = "shipped"
)

// IssuingCardShippingType is the list of possible values for the shipping type
// on an issuing card.
type IssuingCardShippingType string

// List of values that IssuingCardShippingType can take.
const (
	IssuingCardShippingTypeBulk       IssuingCardShippingType = "bulk"
	IssuingCardShippingTypeIndividual IssuingCardShippingType = "individual"
)

// IssuingCardStatus is the list of possible values for status on an issuing card.
type IssuingCardStatus string

// List of values that IssuingCardStatus can take.
const (
	IssuingCardStatusActive   IssuingCardStatus = "active"
	IssuingCardStatusCanceled IssuingCardStatus = "canceled"
	IssuingCardStatusInactive IssuingCardStatus = "inactive"
	IssuingCardStatusPending  IssuingCardStatus = "pending"
)

// IssuingCardType is the type of an issuing card.
type IssuingCardType string

// List of values that IssuingCardType can take.
const (
	IssuingCardTypePhysical IssuingCardType = "physical"
	IssuingCardTypeVirtual  IssuingCardType = "virtual"
)

// AuthorizationControlsParams is the set of parameters that can be used for the shipping parameter.
type AuthorizationControlsParams struct {
	AllowedCategories []*string `form:"allowed_categories"`
	BlockedCategories []*string `form:"blocked_categories"`
	MaxAmount         *int64    `form:"max_amount"`
	MaxApprovals      *int64    `form:"max_approvals"`
}

// IssuingCardShippingParams is the set of parameters that can be used for the shipping parameter.
type IssuingCardShippingParams struct {
	Address *AddressParams `form:"address"`
	Name    string         `form:"name"`
	Type    *string        `form:"type"`
}

// IssuingCardParams is the set of parameters that can be used when creating or updating an issuing card.
type IssuingCardParams struct {
	Params                `form:"*"`
	AuthorizationControls *AuthorizationControlsParams `form:"authorization_controls"`
	Billing               *IssuingBillingParams        `form:"billing"`
	Cardholder            *string                      `form:"cardholder"`
	Currency              *string                      `form:"currency"`
	Name                  *string                      `form:"name"`
	Status                *string                      `form:"status"`
	Shipping              *IssuingCardShippingParams   `form:"shipping"`
	Type                  *string                      `form:"type"`
}

// IssuingCardListParams is the set of parameters that can be used when listing issuing cards.
type IssuingCardListParams struct {
	ListParams   `form:"*"`
	Cardholder   *string           `form:"cardholder"`
	Created      *int64            `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
	ExpMonth     *int64            `form:"exp_month"`
	ExpYear      *int64            `form:"exp_year"`
	Last4        *string           `form:"last4"`
	Source       *string           `form:"source"`
	Status       *string           `form:"status"`
	Type         *string           `form:"type"`
}

// IssuingCardDetails is the resource representing issuing card details.
type IssuingCardDetails struct {
	Card     *IssuingCard `json:"card"`
	CVC      string       `json:"cvc"`
	ExpMonth *string      `form:"exp_month"`
	ExpYear  *string      `form:"exp_year"`
	Number   string       `json:"number"`
	Object   string       `json:"object"`
}

// IssuingCardAuthorizationControls is the resource representing authorization controls on an issuing card.
type IssuingCardAuthorizationControls struct {
	AllowedCategories []string `json:"allowed_categories"`
	BlockedCategories []string `json:"blocked_categories"`
	Currency          Currency `json:"currency"`
	MaxAmount         int64    `json:"max_amount"`
	MaxApprovals      int64    `json:"max_approvals"`
}

// IssuingCardShipping is the resource representing shipping on an issuing card.
type IssuingCardShipping struct {
	Address        *Address                  `json:"address"`
	Carrier        string                    `json:"carrier"`
	ETA            int64                     `json:"eta"`
	Name           string                    `json:"name"`
	Phone          string                    `json:"phone"`
	Status         IssuingCardShippingStatus `json:"status"`
	TrackingNumber string                    `json:"tracking_number"`
	Type           IssuingCardShippingType   `json:"type"`
}

// IssuingCard is the resource representing a Stripe issuing card.
type IssuingCard struct {
	AuthorizationControls *IssuingCardAuthorizationControls `json:"authorization_controls"`
	Billing               *IssuingBilling                   `json:"billing"`
	Brand                 string                            `json:"brand"`
	Cardholder            *IssuingCardholder                `json:"cardholder"`
	Created               int64                             `json:"created"`
	ExpMonth              int64                             `json:"exp_month"`
	ExpYear               int64                             `json:"exp_year"`
	Last4                 string                            `json:"last4"`
	ID                    string                            `json:"id"`
	Livemode              bool                              `json:"livemode"`
	Metadata              map[string]string                 `json:"metadata"`
	Name                  string                            `json:"name"`
	Object                string                            `json:"object"`
	Shipping              *IssuingCardShipping              `json:"shipping"`
	Status                IssuingCardStatus                 `json:"status"`
	Type                  IssuingCardType                   `json:"type"`
}

// IssuingCardList is a list of issuing cards as retrieved from a list endpoint.
type IssuingCardList struct {
	ListMeta
	Data []*IssuingCard `json:"data"`
}

// UnmarshalJSON handles deserialization of an IssuingCard.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (i *IssuingCard) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		i.ID = id
		return nil
	}

	type issuingCard IssuingCard
	var v issuingCard
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = IssuingCard(v)
	return nil
}
