// Package payout provides the /payouts APIs
package payout

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /payouts APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new payout.
func New(params *stripe.PayoutParams) (*stripe.Payout, error) {
	return getC().New(params)
}

// New creates a new payout.
func (c Client) New(params *stripe.PayoutParams) (*stripe.Payout, error) {
	payout := &stripe.Payout{}
	err := c.B.Call(http.MethodPost, "/v1/payouts", c.Key, params, payout)
	return payout, err
}

// Get returns the details of a payout.
func Get(id string, params *stripe.PayoutParams) (*stripe.Payout, error) {
	return getC().Get(id, params)
}

// Get returns the details of a payout.
func (c Client) Get(id string, params *stripe.PayoutParams) (*stripe.Payout, error) {
	path := stripe.FormatURLPath("/v1/payouts/%s", id)
	payout := &stripe.Payout{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, payout)
	return payout, err
}

// Update updates a payout's properties.
func Update(id string, params *stripe.PayoutParams) (*stripe.Payout, error) {
	return getC().Update(id, params)
}

// Update updates a payout's properties.
func (c Client) Update(id string, params *stripe.PayoutParams) (*stripe.Payout, error) {
	path := stripe.FormatURLPath("/v1/payouts/%s", id)
	payout := &stripe.Payout{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, payout)
	return payout, err
}

// Cancel cancels a pending payout.
func Cancel(id string, params *stripe.PayoutParams) (*stripe.Payout, error) {
	return getC().Cancel(id, params)
}

// Cancel cancels a pending payout.
func (c Client) Cancel(id string, params *stripe.PayoutParams) (*stripe.Payout, error) {
	path := stripe.FormatURLPath("/v1/payouts/%s/cancel", id)
	payout := &stripe.Payout{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, payout)
	return payout, err
}

// List returns a list of payouts.
func List(params *stripe.PayoutListParams) *Iter {
	return getC().List(params)
}

// List returns a list of payouts.
func (c Client) List(listParams *stripe.PayoutListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.PayoutList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/payouts", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for payouts.
type Iter struct {
	*stripe.Iter
}

// Payout returns the payout which the iterator is currently pointing to.
func (i *Iter) Payout() *stripe.Payout {
	return i.Current().(*stripe.Payout)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
