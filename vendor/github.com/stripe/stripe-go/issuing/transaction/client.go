// Package transaction provides API functions related to issuing transactions.
//
// For more details, see: https://stripe.com/docs/api/go#issuing_transactions
package transaction

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /issuing/transactions APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns the details of an issuing transaction.
func Get(id string, params *stripe.IssuingTransactionParams) (*stripe.IssuingTransaction, error) {
	return getC().Get(id, params)
}

// Get returns the details of an issuing transaction.
func (c Client) Get(id string, params *stripe.IssuingTransactionParams) (*stripe.IssuingTransaction, error) {
	path := stripe.FormatURLPath("/v1/issuing/transactions/%s", id)
	transaction := &stripe.IssuingTransaction{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, transaction)
	return transaction, err
}

// Update updates an issuing transaction.
func Update(id string, params *stripe.IssuingTransactionParams) (*stripe.IssuingTransaction, error) {
	return getC().Update(id, params)
}

// Update updates an issuing transaction.
func (c Client) Update(id string, params *stripe.IssuingTransactionParams) (*stripe.IssuingTransaction, error) {
	path := stripe.FormatURLPath("/v1/issuing/transactions/%s", id)
	transaction := &stripe.IssuingTransaction{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, transaction)
	return transaction, err
}

// List returns a list of issuing transactions.
func List(params *stripe.IssuingTransactionListParams) *Iter {
	return getC().List(params)
}

// List returns a list of issuing transactions.
func (c Client) List(listParams *stripe.IssuingTransactionListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.IssuingTransactionList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/issuing/transactions", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for issuing transactions.
type Iter struct {
	*stripe.Iter
}

// IssuingTransaction returns the issuing transaction which the iterator is currently pointing to.
func (i *Iter) IssuingTransaction() *stripe.IssuingTransaction {
	return i.Current().(*stripe.IssuingTransaction)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
