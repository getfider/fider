package stripe

import (
	"encoding/json"
)

// BitcoinReceiverListParams is the set of parameters that can be used when listing BitcoinReceivers.
// For more details see https://stripe.com/docs/api/#list_bitcoin_receivers.
type BitcoinReceiverListParams struct {
	ListParams      `form:"*"`
	Active          *bool `form:"active"`
	Filled          *bool `form:"filled"`
	UncapturedFunds *bool `form:"uncaptured_funds"`
}

// BitcoinReceiver is the resource representing a Stripe bitcoin receiver.
// For more details see https://stripe.com/docs/api/#bitcoin_receivers
type BitcoinReceiver struct {
	Active                bool                    `json:"active"`
	Amount                int64                   `json:"amount"`
	AmountReceived        int64                   `json:"amount_received"`
	BitcoinAmount         int64                   `json:"bitcoin_amount"`
	BitcoinAmountReceived int64                   `json:"bitcoin_amount_received"`
	BitcoinURI            string                  `json:"bitcoin_uri"`
	Created               int64                   `json:"created"`
	Currency              Currency                `json:"currency"`
	Customer              string                  `json:"customer"`
	Description           string                  `json:"description"`
	Email                 string                  `json:"email"`
	Filled                bool                    `json:"filled"`
	ID                    string                  `json:"id"`
	InboundAddress        string                  `json:"inbound_address"`
	Metadata              map[string]string       `json:"metadata"`
	Payment               string                  `json:"payment"`
	RefundAddress         string                  `json:"refund_address"`
	RejectTransactions    bool                    `json:"reject_transactions"`
	Transactions          *BitcoinTransactionList `json:"transactions"`
}

// BitcoinReceiverList is a list of bitcoin receivers as retrieved from a list endpoint.
type BitcoinReceiverList struct {
	ListMeta
	Data []*BitcoinReceiver `json:"data"`
}

// UnmarshalJSON handles deserialization of a BitcoinReceiver.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (r *BitcoinReceiver) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		r.ID = id
		return nil
	}

	type bitcoinReceiver BitcoinReceiver
	var v bitcoinReceiver
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = BitcoinReceiver(v)
	return nil
}
