// Package sourcetransaction provides the /source/transactions APIs.
package sourcetransaction

import (
	"errors"
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /sources/:source_id/transactions APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// List returns a list of source transactions.
func List(params *stripe.SourceTransactionListParams) *Iter {
	return getC().List(params)
}

// List returns a list of source transactions.
func (c Client) List(listParams *stripe.SourceTransactionListParams) *Iter {
	var outerErr error
	var path string

	if listParams == nil || listParams.Source == nil {
		outerErr = errors.New("Invalid source transaction params: Source needs to be set")
	} else {
		path = stripe.FormatURLPath("/v1/sources/%s/source_transactions",
			stripe.StringValue(listParams.Source))
	}

	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.SourceTransactionList{}

		if outerErr != nil {
			return nil, list.ListMeta, outerErr
		}

		err := c.B.CallRaw(http.MethodGet, path, c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for source transactions.
type Iter struct {
	*stripe.Iter
}

// SourceTransaction returns the source transaction which the iterator is currently pointing to.
func (i *Iter) SourceTransaction() *stripe.SourceTransaction {
	return i.Current().(*stripe.SourceTransaction)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
