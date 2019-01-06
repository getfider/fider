// Package refund provides the /refunds APIs
package refund

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /refunds APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a refund.
func New(params *stripe.RefundParams) (*stripe.Refund, error) {
	return getC().New(params)
}

// New creates a refund.
func (c Client) New(params *stripe.RefundParams) (*stripe.Refund, error) {
	refund := &stripe.Refund{}
	err := c.B.Call(http.MethodPost, "/v1/refunds", c.Key, params, refund)
	return refund, err
}

// Get returns the details of a refund.
func Get(id string, params *stripe.RefundParams) (*stripe.Refund, error) {
	return getC().Get(id, params)
}

// Get returns the details of a refund.
func (c Client) Get(id string, params *stripe.RefundParams) (*stripe.Refund, error) {
	path := stripe.FormatURLPath("/v1/refunds/%s", id)
	refund := &stripe.Refund{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, refund)
	return refund, err
}

// Update updates a refund's properties.
func Update(id string, params *stripe.RefundParams) (*stripe.Refund, error) {
	return getC().Update(id, params)
}

// Update updates a refund's properties.
func (c Client) Update(id string, params *stripe.RefundParams) (*stripe.Refund, error) {
	path := stripe.FormatURLPath("/v1/refunds/%s", id)
	refund := &stripe.Refund{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, refund)
	return refund, err
}

// List returns a list of refunds.
func List(params *stripe.RefundListParams) *Iter {
	return getC().List(params)
}

// List returns a list of refunds.
func (c Client) List(listParams *stripe.RefundListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.RefundList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/refunds", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for refunds.
type Iter struct {
	*stripe.Iter
}

// Refund returns the refund which the iterator is currently pointing to.
func (i *Iter) Refund() *stripe.Refund {
	return i.Current().(*stripe.Refund)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
