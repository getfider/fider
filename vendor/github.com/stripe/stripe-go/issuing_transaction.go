package stripe

import "encoding/json"

// IssuingTransactionType is the type of an issuing transaction.
type IssuingTransactionType string

// List of values that IssuingTransactionType can take.
const (
	IssuingTransactionTypeCapture        IssuingTransactionType = "capture"
	IssuingTransactionTypeCashWithdrawal IssuingTransactionType = "cash_withdrawal"
	IssuingTransactionTypeRefund         IssuingTransactionType = "refund"
	IssuingTransactionTypeRefundReversal IssuingTransactionType = "refund_reversal"
)

// IssuingTransactionParams is the set of parameters that can be used when creating or updating an issuing transaction.
type IssuingTransactionParams struct {
	Params `form:"*"`
}

// IssuingTransactionListParams is the set of parameters that can be used when listing issuing transactions.
type IssuingTransactionListParams struct {
	ListParams   `form:"*"`
	Card         *string           `form:"card"`
	Cardholder   *string           `form:"cardholder"`
	Created      *int64            `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
}

// IssuingTransaction is the resource representing a Stripe issuing transaction.
type IssuingTransaction struct {
	Authorization      *IssuingAuthorization  `json:"authorization"`
	BalanceTransaction *BalanceTransaction    `json:"balance_transaction"`
	Card               *IssuingCard           `json:"card"`
	Cardholder         *IssuingCardholder     `json:"cardholder"`
	Created            int64                  `json:"created"`
	Currency           Currency               `json:"currency"`
	Dispute            *IssuingDispute        `json:"dispute"`
	ID                 string                 `json:"id"`
	Livemode           bool                   `json:"livemode"`
	MerchantData       *IssuingMerchantData   `json:"merchant_data"`
	Metadata           map[string]string      `json:"metadata"`
	Object             string                 `json:"object"`
	Type               IssuingTransactionType `json:"type"`
}

// IssuingTransactionList is a list of issuing transactions as retrieved from a list endpoint.
type IssuingTransactionList struct {
	ListMeta
	Data []*IssuingTransaction `json:"data"`
}

// UnmarshalJSON handles deserialization of an IssuingTransaction.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (i *IssuingTransaction) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		i.ID = id
		return nil
	}

	type issuingTransaction IssuingTransaction
	var v issuingTransaction
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = IssuingTransaction(v)
	return nil
}
