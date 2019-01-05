package stripe

import "encoding/json"

// SourceTransactionListParams is the set of parameters that can be used when listing SourceTransactions.
type SourceTransactionListParams struct {
	ListParams `form:"*"`
	Source     *string `form:"-"` // Sent in with the URL
}

// SourceTransactionList is a list object for SourceTransactions.
type SourceTransactionList struct {
	ListMeta
	Data []*SourceTransaction `json:"data"`
}

// SourceTransaction is the resource representing a Stripe source transaction.
type SourceTransaction struct {
	Amount       int64    `json:"amount"`
	Created      int64    `json:"created"`
	Currency     Currency `json:"currency"`
	CustomerData string   `json:"customer_data"`
	ID           string   `json:"id"`
	Livemode     bool     `json:"livemode"`
	Source       string   `json:"source"`
	Type         string   `json:"type"`
	TypeData     map[string]interface{}
}

// UnmarshalJSON handles deserialization of a SourceTransaction. This custom
// unmarshaling is needed to extract the type specific data (accessible under
// `TypeData`) but stored in JSON under a hash named after the `type` of the
// source transaction.
func (t *SourceTransaction) UnmarshalJSON(data []byte) error {
	type sourceTransaction SourceTransaction
	var v sourceTransaction
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*t = SourceTransaction(v)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	if d, ok := raw[t.Type]; ok {
		if m, ok := d.(map[string]interface{}); ok {
			t.TypeData = m
		}
	}

	return nil
}
