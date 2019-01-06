// Package charge provides API functions related to charges.
//
// For more details, see: https://stripe.com/docs/api/go#charges.
package charge

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke APIs related to charges.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new charge.
func New(params *stripe.ChargeParams) (*stripe.Charge, error) {
	return getC().New(params)
}

// New creates a new charge.
func (c Client) New(params *stripe.ChargeParams) (*stripe.Charge, error) {
	charge := &stripe.Charge{}
	err := c.B.Call(http.MethodPost, "/v1/charges", c.Key, params, charge)
	return charge, err
}

// Get retrieves a charge.
func Get(id string, params *stripe.ChargeParams) (*stripe.Charge, error) {
	return getC().Get(id, params)
}

// Get retrieves a charge.
func (c Client) Get(id string, params *stripe.ChargeParams) (*stripe.Charge, error) {
	path := stripe.FormatURLPath("/v1/charges/%s", id)
	charge := &stripe.Charge{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, charge)
	return charge, err
}

// Update updates a charge.
func Update(id string, params *stripe.ChargeParams) (*stripe.Charge, error) {
	return getC().Update(id, params)
}

// Update updates a charge.
func (c Client) Update(id string, params *stripe.ChargeParams) (*stripe.Charge, error) {
	path := stripe.FormatURLPath("/v1/charges/%s", id)
	charge := &stripe.Charge{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, charge)
	return charge, err
}

// Capture captures a charge that's not yet captured.
func Capture(id string, params *stripe.CaptureParams) (*stripe.Charge, error) {
	return getC().Capture(id, params)
}

// Capture captures a charge that's not yet captured.
func (c Client) Capture(id string, params *stripe.CaptureParams) (*stripe.Charge, error) {
	path := stripe.FormatURLPath("/v1/charges/%s/capture", id)
	charge := &stripe.Charge{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, charge)
	return charge, err
}

// List returns an iterator that iterates all charges.
func List(params *stripe.ChargeListParams) *Iter {
	return getC().List(params)
}

// List returns an iterator that iterates all charges.
func (c Client) List(listParams *stripe.ChargeListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.ChargeList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/charges", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for charges.
type Iter struct {
	*stripe.Iter
}

// Charge returns the charge which the iterator is currently pointing to.
func (i *Iter) Charge() *stripe.Charge {
	return i.Current().(*stripe.Charge)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
