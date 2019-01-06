// Package fee provides the /application_fees APIs
package fee

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke application_fees APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns the details of an application fee.
func Get(id string, params *stripe.ApplicationFeeParams) (*stripe.ApplicationFee, error) {
	return getC().Get(id, params)
}

// Get returns the details of an application fee.
func (c Client) Get(id string, params *stripe.ApplicationFeeParams) (*stripe.ApplicationFee, error) {
	path := stripe.FormatURLPath("/v1/application_fees/%s", id)
	fee := &stripe.ApplicationFee{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, fee)
	return fee, err
}

// List returns a list of application fees.
func List(params *stripe.ApplicationFeeListParams) *Iter {
	return getC().List(params)
}

// List returns a list of application fees.
func (c Client) List(listParams *stripe.ApplicationFeeListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.ApplicationFeeList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/application_fees", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for application fees.
type Iter struct {
	*stripe.Iter
}

// ApplicationFee returns the application fee which the iterator is currently pointing to.
func (i *Iter) ApplicationFee() *stripe.ApplicationFee {
	return i.Current().(*stripe.ApplicationFee)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
