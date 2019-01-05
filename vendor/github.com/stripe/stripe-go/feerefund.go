package stripe

import (
	"encoding/json"
)

// FeeRefundParams is the set of parameters that can be used when refunding an application fee.
// For more details see https://stripe.com/docs/api#fee_refund.
type FeeRefundParams struct {
	Params         `form:"*"`
	Amount         *int64  `form:"amount"`
	ApplicationFee *string `form:"-"` // Included in the URL
}

// FeeRefundListParams is the set of parameters that can be used when listing application fee refunds.
// For more details see https://stripe.com/docs/api#list_fee_refunds.
type FeeRefundListParams struct {
	ListParams     `form:"*"`
	ApplicationFee *string `form:"-"` // Included in the URL
}

// FeeRefund is the resource representing a Stripe application fee refund.
// For more details see https://stripe.com/docs/api#fee_refunds.
type FeeRefund struct {
	Amount             int64               `json:"amount"`
	BalanceTransaction *BalanceTransaction `json:"balance_transaction"`
	Created            int64               `json:"created"`
	Currency           Currency            `json:"currency"`
	Fee                *ApplicationFee     `json:"fee"`
	ID                 string              `json:"id"`
	Metadata           map[string]string   `json:"metadata"`
}

// FeeRefundList is a list object for application fee refunds.
type FeeRefundList struct {
	ListMeta
	Data []*FeeRefund `json:"data"`
}

// UnmarshalJSON handles deserialization of a FeeRefund.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (r *FeeRefund) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		r.ID = id
		return nil
	}

	type feeRefund FeeRefund
	var v feeRefund
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = FeeRefund(v)
	return nil
}
