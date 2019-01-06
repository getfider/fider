package sku

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /skus APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new SKU.
func New(params *stripe.SKUParams) (*stripe.SKU, error) {
	return getC().New(params)
}

// New creates a new SKU.
func (c Client) New(params *stripe.SKUParams) (*stripe.SKU, error) {
	s := &stripe.SKU{}
	err := c.B.Call(http.MethodPost, "/v1/skus", c.Key, params, s)
	return s, err
}

// Update updates a SKU's properties.
func Update(id string, params *stripe.SKUParams) (*stripe.SKU, error) {
	return getC().Update(id, params)
}

// Update updates a SKU's properties.
func (c Client) Update(id string, params *stripe.SKUParams) (*stripe.SKU, error) {
	path := stripe.FormatURLPath("/v1/skus/%s", id)
	s := &stripe.SKU{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, s)
	return s, err
}

// Get returns the details of a SKU.
func Get(id string, params *stripe.SKUParams) (*stripe.SKU, error) {
	return getC().Get(id, params)
}

// Get returns the details of a SKU.
func (c Client) Get(id string, params *stripe.SKUParams) (*stripe.SKU, error) {
	path := stripe.FormatURLPath("/v1/skus/%s", id)
	s := &stripe.SKU{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, s)
	return s, err
}

// List returns a list of SKUs.
func List(params *stripe.SKUListParams) *Iter {
	return getC().List(params)
}

// List returns a list of SKUs.
func (c Client) List(listParams *stripe.SKUListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.SKUList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/skus", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Del removes a SKU.
func Del(id string, params *stripe.SKUParams) (*stripe.SKU, error) {
	return getC().Del(id, params)
}

// Del removes a SKU.
func (c Client) Del(id string, params *stripe.SKUParams) (*stripe.SKU, error) {
	path := stripe.FormatURLPath("/v1/skus/%s", id)
	s := &stripe.SKU{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, s)

	return s, err
}

// Iter is an iterator for SKUs.
type Iter struct {
	*stripe.Iter
}

// SKU returns the SKU which the iterator is currently pointing to.
func (i *Iter) SKU() *stripe.SKU {
	return i.Current().(*stripe.SKU)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
