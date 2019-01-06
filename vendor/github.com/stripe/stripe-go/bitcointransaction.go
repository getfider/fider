package stripe

import "encoding/json"

// BitcoinTransactionListParams is the set of parameters that can be used when listing BitcoinTransactions.
type BitcoinTransactionListParams struct {
	ListParams `form:"*"`
	Customer   *string `form:"customer"`
	Receiver   *string `form:"-"` // Sent in with the URL
}

// BitcoinTransactionList is a list object for BitcoinTransactions.
// It is a child object of BitcoinRecievers
// For more details see https://stripe.com/docs/api/#retrieve_bitcoin_receiver
type BitcoinTransactionList struct {
	ListMeta
	Data []*BitcoinTransaction `json:"data"`
}

// BitcoinTransaction is the resource representing a Stripe bitcoin transaction.
// For more details see https://stripe.com/docs/api/#bitcoin_receivers
type BitcoinTransaction struct {
	Amount        int64    `json:"amount"`
	BitcoinAmount int64    `json:"bitcoin_amount"`
	Created       int64    `json:"created"`
	Currency      Currency `json:"currency"`
	Customer      string   `json:"customer"`
	ID            string   `json:"id"`
	Receiver      string   `json:"receiver"`
}

// UnmarshalJSON handles deserialization of a BitcoinTransaction.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (bt *BitcoinTransaction) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		bt.ID = id
		return nil
	}

	type bitcoinTransaction BitcoinTransaction
	var v bitcoinTransaction
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*bt = BitcoinTransaction(v)
	return nil
}
