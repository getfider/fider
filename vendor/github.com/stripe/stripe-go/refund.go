package stripe

import "encoding/json"

// RefundReason is, if set, the reason the refund is being made
type RefundReason string

// List of values that RefundReason can take.
const (
	RefundReasonDuplicate           RefundReason = "duplicate"
	RefundReasonFraudulent          RefundReason = "fraudulent"
	RefundReasonRequestedByCustomer RefundReason = "requested_by_customer"
)

// RefundStatus is the status of the refund.
type RefundStatus string

// List of values that RefundStatus can take.
const (
	RefundStatusCanceled  RefundStatus = "canceled"
	RefundStatusFailed    RefundStatus = "failed"
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusSucceeded RefundStatus = "succeeded"
)

// RefundParams is the set of parameters that can be used when refunding a charge.
// For more details see https://stripe.com/docs/api#refund.
type RefundParams struct {
	Params               `form:"*"`
	Amount               *int64  `form:"amount"`
	Charge               *string `form:"charge"`
	Reason               *string `form:"reason"`
	RefundApplicationFee *bool   `form:"refund_application_fee"`
	ReverseTransfer      *bool   `form:"reverse_transfer"`
}

// RefundListParams is the set of parameters that can be used when listing refunds.
// For more details see https://stripe.com/docs/api#list_refunds.
type RefundListParams struct {
	ListParams `form:"*"`
}

// Refund is the resource representing a Stripe refund.
// For more details see https://stripe.com/docs/api#refunds.
type Refund struct {
	Amount                    int64               `json:"amount"`
	BalanceTransaction        *BalanceTransaction `json:"balance_transaction"`
	Charge                    *Charge             `json:"charge"`
	Created                   int64               `json:"created"`
	Currency                  Currency            `json:"currency"`
	FailureBalanceTransaction *BalanceTransaction `json:"failure_balance_transaction"`
	ID                        string              `json:"id"`
	Metadata                  map[string]string   `json:"metadata"`
	Object                    string              `json:"object"`
	Reason                    RefundReason        `json:"reason"`
	ReceiptNumber             string              `json:"receipt_number"`
	SourceTransferReversal    *Reversal           `json:"source_transfer_reversal"`
	Status                    RefundStatus        `json:"status"`
	TransferReversal          *Reversal           `json:"transfer_reversal"`
}

// RefundList is a list object for refunds.
type RefundList struct {
	ListMeta
	Data []*Refund `json:"data"`
}

// UnmarshalJSON handles deserialization of a Refund.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (r *Refund) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		r.ID = id
		return nil
	}

	type refund Refund
	var v refund
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = Refund(v)
	return nil
}
