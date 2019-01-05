// Package dispute provides API functions related to issuing disputes.
//
// For more details, see: https://stripe.com/docs/api/go#issuing_disputes
package dispute

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /issuing/disputes APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new issuing dispute.
func New(params *stripe.IssuingDisputeParams) (*stripe.IssuingDispute, error) {
	return getC().New(params)
}

// New creates a new issuing dispute.
func (c Client) New(params *stripe.IssuingDisputeParams) (*stripe.IssuingDispute, error) {
	dispute := &stripe.IssuingDispute{}
	err := c.B.Call(http.MethodPost, "/v1/issuing/disputes", c.Key, params, dispute)
	return dispute, err
}

// Get returns the details of an issuing dispute.
func Get(id string, params *stripe.IssuingDisputeParams) (*stripe.IssuingDispute, error) {
	return getC().Get(id, params)
}

// Get returns the details of an issuing dispute.
func (c Client) Get(id string, params *stripe.IssuingDisputeParams) (*stripe.IssuingDispute, error) {
	path := stripe.FormatURLPath("/v1/issuing/disputes/%s", id)
	dispute := &stripe.IssuingDispute{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, dispute)
	return dispute, err
}

// Update updates an issuing dispute.
func Update(id string, params *stripe.IssuingDisputeParams) (*stripe.IssuingDispute, error) {
	return getC().Update(id, params)
}

// Update updates an issuing dispute.
func (c Client) Update(id string, params *stripe.IssuingDisputeParams) (*stripe.IssuingDispute, error) {
	path := stripe.FormatURLPath("/v1/issuing/disputes/%s", id)
	dispute := &stripe.IssuingDispute{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, dispute)
	return dispute, err
}

// List returns a list of issuing disputes.
func List(params *stripe.IssuingDisputeListParams) *Iter {
	return getC().List(params)
}

// List returns a list of issuing disputes.
func (c Client) List(listParams *stripe.IssuingDisputeListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.IssuingDisputeList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/issuing/disputes", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for issuing disputes.
type Iter struct {
	*stripe.Iter
}

// IssuingDispute returns the issuing dispute which the iterator is currently pointing to.
func (i *Iter) IssuingDispute() *stripe.IssuingDispute {
	return i.Current().(*stripe.IssuingDispute)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
