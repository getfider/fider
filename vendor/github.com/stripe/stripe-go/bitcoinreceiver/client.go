// Package bitcoinreceiver provides the /bitcoin/receivers APIs.
//
// Note that this entire package is deprecated. Please use the new sources API
// instead.
package bitcoinreceiver

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /bitcoin/receivers APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns the details of a bitcoin receiver.
func Get(id string) (*stripe.BitcoinReceiver, error) {
	return getC().Get(id)
}

// Get returns the details of a bitcoin receiver.
func (c Client) Get(id string) (*stripe.BitcoinReceiver, error) {
	path := stripe.FormatURLPath("/v1/bitcoin/receivers/%s", id)
	bitcoinReceiver := &stripe.BitcoinReceiver{}
	err := c.B.Call(http.MethodGet, path, c.Key, nil, bitcoinReceiver)
	return bitcoinReceiver, err
}

// List returns a list of bitcoin receivers.
func List(params *stripe.BitcoinReceiverListParams) *Iter {
	return getC().List(params)
}

// List returns a list of bitcoin receivers.
func (c Client) List(listParams *stripe.BitcoinReceiverListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.BitcoinReceiverList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/bitcoin/receivers", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for bitcoin receivers.
type Iter struct {
	*stripe.Iter
}

// BitcoinReceiver returns the bitcoin receiver which the iterator is currently pointing to.
func (i *Iter) BitcoinReceiver() *stripe.BitcoinReceiver {
	return i.Current().(*stripe.BitcoinReceiver)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
