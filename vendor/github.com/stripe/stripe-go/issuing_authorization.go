package stripe

import "encoding/json"

// IssuingAuthorizationAuthorizationMethod is the list of possible values for the authorization method
// on an issuing authorization.
type IssuingAuthorizationAuthorizationMethod string

// List of values that IssuingAuthorizationAuthorizationMethod can take.
const (
	IssuingAuthorizationAuthorizationMethodChip        IssuingAuthorizationAuthorizationMethod = "chip"
	IssuingAuthorizationAuthorizationMethodContactless IssuingAuthorizationAuthorizationMethod = "contactless"
	IssuingAuthorizationAuthorizationMethodKeyedIn     IssuingAuthorizationAuthorizationMethod = "keyed_in"
	IssuingAuthorizationAuthorizationMethodOnline      IssuingAuthorizationAuthorizationMethod = "online"
	IssuingAuthorizationAuthorizationMethodSwipe       IssuingAuthorizationAuthorizationMethod = "swipe"
)

// IssuingAuthorizationRequestHistoryReason is the list of possible values for the request history
// reason on an issuing authorization.
type IssuingAuthorizationRequestHistoryReason string

// List of values that IssuingAuthorizationRequestHistoryReason can take.
const (
	IssuingAuthorizationRequestHistoryReasonAuthorizationControls IssuingAuthorizationRequestHistoryReason = "authorization_controls"
	IssuingAuthorizationRequestHistoryReasonCardActive            IssuingAuthorizationRequestHistoryReason = "card_active"
	IssuingAuthorizationRequestHistoryReasonCardInactive          IssuingAuthorizationRequestHistoryReason = "card_inactive"
	IssuingAuthorizationRequestHistoryReasonInsufficientFunds     IssuingAuthorizationRequestHistoryReason = "insufficient_funds"
	IssuingAuthorizationRequestHistoryReasonWebhookApproved       IssuingAuthorizationRequestHistoryReason = "webhook_approved"
	IssuingAuthorizationRequestHistoryReasonWebhookDeclined       IssuingAuthorizationRequestHistoryReason = "webhook_declined"
	IssuingAuthorizationRequestHistoryReasonWebhookTimeout        IssuingAuthorizationRequestHistoryReason = "webhook_timeout"
)

// IssuingAuthorizationStatus is the possible values for status for an issuing authorization.
type IssuingAuthorizationStatus string

// List of values that IssuingAuthorizationStatus can take.
const (
	IssuingAuthorizationStatusClosed   IssuingAuthorizationStatus = "closed"
	IssuingAuthorizationStatusPending  IssuingAuthorizationStatus = "pending"
	IssuingAuthorizationStatusReversed IssuingAuthorizationStatus = "reversed"
)

// IssuingAuthorizationVerificationDataCheck is the list of possible values for result of a check
// for verification data on an issuing authorization.
type IssuingAuthorizationVerificationDataCheck string

// List of values that IssuingAuthorizationVerificationDataCheck can take.
const (
	IssuingAuthorizationVerificationDataCheckMatch       IssuingAuthorizationVerificationDataCheck = "match"
	IssuingAuthorizationVerificationDataCheckMismatch    IssuingAuthorizationVerificationDataCheck = "mismatch"
	IssuingAuthorizationVerificationDataCheckNotProvided IssuingAuthorizationVerificationDataCheck = "not_provided"
)

// IssuingAuthorizationParams is the set of parameters that can be used when updating an issuing authorization.
type IssuingAuthorizationParams struct {
	Params `form:"*"`
}

// IssuingAuthorizationListParams is the set of parameters that can be used when listing issuing authorizations.
type IssuingAuthorizationListParams struct {
	ListParams   `form:"*"`
	Card         *string           `form:"card"`
	Cardholder   *string           `form:"cardholder"`
	Created      *int64            `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
	Status       *string           `form:"status"`
}

// IssuingAuthorizationAuthorizationControls is the resource representing authorization controls on an issuing authorization.
type IssuingAuthorizationAuthorizationControls struct {
	AllowedCategories []string `json:"allowed_categories"`
	BlockedCategories []string `json:"blocked_categories"`
	Currency          Currency `json:"currency"`
	MaxAmount         int64    `json:"max_amount"`
	MaxApprovals      int64    `json:"max_approvals"`
}

// IssuingAuthorizationRequestHistory is the resource representing a request history on an issuing authorization.
type IssuingAuthorizationRequestHistory struct {
	Approved           bool                                     `json:"approved"`
	AuthorizedAmount   int64                                    `json:"authorized_amount"`
	AuthorizedCurrency Currency                                 `json:"authorized_currency"`
	Created            int64                                    `json:"created"`
	HeldAmount         int64                                    `json:"held_amount"`
	HeldCurrency       Currency                                 `json:"held_currency"`
	Reason             IssuingAuthorizationRequestHistoryReason `json:"reason"`
}

// IssuingAuthorizationVerificationData is the resource representing verification data on an issuing authorization.
type IssuingAuthorizationVerificationData struct {
	AddressLine1Check IssuingAuthorizationVerificationDataCheck `json:"address_line1_check"`
	AddressZipCheck   IssuingAuthorizationVerificationDataCheck `json:"address_zip_check"`
	CVCCheck          IssuingAuthorizationVerificationDataCheck `json:"cvc_check"`
}

// IssuingAuthorization is the resource representing a Stripe issuing authorization.
type IssuingAuthorization struct {
	Approved                 bool                                    `json:"approved"`
	AuthorizationMethod      IssuingAuthorizationAuthorizationMethod `json:"authorization_method"`
	AuthorizedAmount         int64                                   `json:"authorized_amount"`
	AuthorizedCurrency       Currency                                `json:"authorized_currency"`
	BalanceTransactions      []*BalanceTransaction                   `json:"balance_transactions"`
	Card                     *IssuingCard                            `json:"card"`
	Cardholder               *IssuingCardholder                      `json:"cardholder"`
	Created                  int64                                   `json:"created"`
	HeldAmount               int64                                   `json:"held_amount"`
	HeldCurrency             Currency                                `json:"held_currency"`
	ID                       string                                  `json:"id"`
	IsHeldAmountControllable bool                                    `json:"is_held_amount_controllable"`
	Livemode                 bool                                    `json:"livemode"`
	MerchantData             *IssuingMerchantData                    `json:"merchant_data"`
	Metadata                 map[string]string                       `json:"metadata"`
	Object                   string                                  `json:"object"`
	PendingAuthorizedAmount  int64                                   `json:"pending_authorized_amount"`
	PendingHeldAmount        int64                                   `json:"pending_held_amount"`
	RequestHistory           []*IssuingAuthorizationRequestHistory   `json:"request_history"`
	Status                   IssuingAuthorizationStatus              `json:"status"`
	Transactions             []*IssuingTransaction                   `json:"transactions"`
	VerificationData         *IssuingAuthorizationVerificationData   `json:"verification_data"`
}

// IssuingMerchantData is the resource representing merchant data on Issuing APIs.
type IssuingMerchantData struct {
	Category   string `json:"category"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Name       string `json:"name"`
	NetworkID  string `json:"network_id"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
}

// IssuingAuthorizationList is a list of issuing authorizations as retrieved from a list endpoint.
type IssuingAuthorizationList struct {
	ListMeta
	Data []*IssuingAuthorization `json:"data"`
}

// UnmarshalJSON handles deserialization of an IssuingAuthorization.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (i *IssuingAuthorization) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		i.ID = id
		return nil
	}

	type issuingAuthorization IssuingAuthorization
	var v issuingAuthorization
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = IssuingAuthorization(v)
	return nil
}
