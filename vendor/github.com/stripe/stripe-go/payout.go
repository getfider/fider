package stripe

import "encoding/json"

// PayoutDestinationType consts represent valid payout destinations.
type PayoutDestinationType string

// List of values that PayoutDestinationType can take.
const (
	PayoutDestinationTypeBankAccount PayoutDestinationType = "bank_account"
	PayoutDestinationTypeCard        PayoutDestinationType = "card"
)

// PayoutFailureCode is the list of allowed values for the payout's failure code.
type PayoutFailureCode string

// List of values that PayoutFailureCode can take.
const (
	PayoutFailureCodeAccountClosed         PayoutFailureCode = "account_closed"
	PayoutFailureCodeAccountFrozen         PayoutFailureCode = "account_frozen"
	PayoutFailureCodeBankAccountRestricted PayoutFailureCode = "bank_account_restricted"
	PayoutFailureCodeBankOwnershipChanged  PayoutFailureCode = "bank_ownership_changed"
	PayoutFailureCodeCouldNotProcess       PayoutFailureCode = "could_not_process"
	PayoutFailureCodeDebitNotAuthorized    PayoutFailureCode = "debit_not_authorized"
	PayoutFailureCodeInsufficientFunds     PayoutFailureCode = "insufficient_funds"
	PayoutFailureCodeInvalidAccountNumber  PayoutFailureCode = "invalid_account_number"
	PayoutFailureCodeInvalidCurrency       PayoutFailureCode = "invalid_currency"
	PayoutFailureCodeNoAccount             PayoutFailureCode = "no_account"
)

// PayoutSourceType is the list of allowed values for the payout's source_type field.
type PayoutSourceType string

// List of values that PayoutSourceType can take.
const (
	PayoutSourceTypeAlipayAccount   PayoutSourceType = "alipay_account"
	PayoutSourceTypeBankAccount     PayoutSourceType = "bank_account"
	PayoutSourceTypeBitcoinReceiver PayoutSourceType = "bitcoin_receiver"
	PayoutSourceTypeCard            PayoutSourceType = "card"
)

// PayoutStatus is the list of allowed values for the payout's status.
type PayoutStatus string

// List of values that PayoutStatus can take.
const (
	PayoutStatusCanceled  PayoutStatus = "canceled"
	PayoutStatusFailed    PayoutStatus = "failed"
	PayoutStatusInTransit PayoutStatus = "in_transit"
	PayoutStatusPaid      PayoutStatus = "paid"
	PayoutStatusPending   PayoutStatus = "pending"
)

// PayoutType is the list of allowed values for the payout's type.
type PayoutType string

// List of values that PayoutType can take.
const (
	PayoutTypeBank PayoutType = "bank_account"
	PayoutTypeCard PayoutType = "card"
)

// PayoutMethodType represents the type of payout
type PayoutMethodType string

// List of values that PayoutMethodType can take.
const (
	PayoutMethodInstant  PayoutMethodType = "instant"
	PayoutMethodStandard PayoutMethodType = "standard"
)

// PayoutDestination describes the destination of a Payout.
// The Type should indicate which object is fleshed out
// For more details see https://stripe.com/docs/api/go#payout_object
type PayoutDestination struct {
	BankAccount *BankAccount          `json:"-"`
	Card        *Card                 `json:"-"`
	ID          string                `json:"id"`
	Type        PayoutDestinationType `json:"object"`
}

// PayoutParams is the set of parameters that can be used when creating or updating a payout.
// For more details see https://stripe.com/docs/api#create_payout and https://stripe.com/docs/api#update_payout.
type PayoutParams struct {
	Params              `form:"*"`
	Amount              *int64  `form:"amount"`
	Currency            *string `form:"currency"`
	Destination         *string `form:"destination"`
	Method              *string `form:"method"`
	SourceType          *string `form:"source_type"`
	StatementDescriptor *string `form:"statement_descriptor"`
}

// PayoutListParams is the set of parameters that can be used when listing payouts.
// For more details see https://stripe.com/docs/api#list_payouts.
type PayoutListParams struct {
	ListParams       `form:"*"`
	ArrivalDate      *int64            `form:"arrival_date"`
	ArrivalDateRange *RangeQueryParams `form:"arrival_date"`
	Created          *int64            `form:"created"`
	CreatedRange     *RangeQueryParams `form:"created"`
	Destination      *string           `form:"destination"`
	Status           *string           `form:"status"`
}

// Payout is the resource representing a Stripe payout.
// For more details see https://stripe.com/docs/api#payouts.
type Payout struct {
	Amount                    int64               `json:"amount"`
	ArrivalDate               int64               `json:"arrival_date"`
	Automatic                 bool                `json:"automatic"`
	BalanceTransaction        *BalanceTransaction `json:"balance_transaction"`
	BankAccount               *BankAccount        `json:"bank_account"`
	Card                      *Card               `json:"card"`
	Created                   int64               `json:"created"`
	Currency                  Currency            `json:"currency"`
	Destination               *PayoutDestination  `json:"destination"`
	FailureBalanceTransaction *BalanceTransaction `json:"failure_balance_transaction"`
	FailureCode               PayoutFailureCode   `json:"failure_code"`
	FailureMessage            string              `json:"failure_message"`
	ID                        string              `json:"id"`
	Livemode                  bool                `json:"livemode"`
	Metadata                  map[string]string   `json:"metadata"`
	Method                    PayoutMethodType    `json:"method"`
	SourceType                PayoutSourceType    `json:"source_type"`
	StatementDescriptor       string              `json:"statement_descriptor"`
	Status                    PayoutStatus        `json:"status"`
	Type                      PayoutType          `json:"type"`
}

// PayoutList is a list of payouts as retrieved from a list endpoint.
type PayoutList struct {
	ListMeta
	Data []*Payout `json:"data"`
}

// UnmarshalJSON handles deserialization of a Payout.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (p *Payout) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		p.ID = id
		return nil
	}

	type payout Payout
	var v payout
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = Payout(v)
	return nil
}

// UnmarshalJSON handles deserialization of a PayoutDestination.
// This custom unmarshaling is needed because the specific
// type of destination it refers to is specified in the JSON
func (d *PayoutDestination) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		d.ID = id
		return nil
	}

	type payoutDestination PayoutDestination
	var v payoutDestination
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var err error
	*d = PayoutDestination(v)

	switch d.Type {
	case PayoutDestinationTypeBankAccount:
		err = json.Unmarshal(data, &d.BankAccount)
	case PayoutDestinationTypeCard:
		err = json.Unmarshal(data, &d.Card)
	}

	return err
}

// MarshalJSON handles serialization of a PayoutDestination.
// This custom marshaling is needed because we can only send a string
// ID as a destination, even though it can be expanded to a full
// object when retrieving
func (d *PayoutDestination) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.ID)
}
