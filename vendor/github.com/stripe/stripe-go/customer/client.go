// Package customer provides the /customers APIs
package customer

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /customers APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new customer.
func New(params *stripe.CustomerParams) (*stripe.Customer, error) {
	return getC().New(params)
}

// New creates a new customer.
func (c Client) New(params *stripe.CustomerParams) (*stripe.Customer, error) {
	cust := &stripe.Customer{}
	err := c.B.Call(http.MethodPost, "/v1/customers", c.Key, params, cust)
	return cust, err
}

// Get returns the details of a customer.
func Get(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	return getC().Get(id, params)
}

// Get returns the details of a customer.
func (c Client) Get(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	path := stripe.FormatURLPath("/v1/customers/%s", id)
	cust := &stripe.Customer{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, cust)
	return cust, err
}

// Update updates a customer's properties.
func Update(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	return getC().Update(id, params)
}

// Update updates a customer's properties.
func (c Client) Update(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	path := stripe.FormatURLPath("/v1/customers/%s", id)
	cust := &stripe.Customer{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, cust)
	return cust, err
}

// Del removes a customer.
func Del(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	return getC().Del(id, params)
}

// Del removes a customer.
func (c Client) Del(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	path := stripe.FormatURLPath("/v1/customers/%s", id)
	cust := &stripe.Customer{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, cust)
	return cust, err
}

// List returns a list of customers.
func List(params *stripe.CustomerListParams) *Iter {
	return getC().List(params)
}

// List returns a list of customers.
func (c Client) List(listParams *stripe.CustomerListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.CustomerList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/customers", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for customers.
type Iter struct {
	*stripe.Iter
}

// Customer returns the customer which the iterator is currently pointing to.
func (i *Iter) Customer() *stripe.Customer {
	return i.Current().(*stripe.Customer)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
