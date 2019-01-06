package order

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /orders APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new order.
func New(params *stripe.OrderParams) (*stripe.Order, error) {
	return getC().New(params)
}

// New creates a new order.
func (c Client) New(params *stripe.OrderParams) (*stripe.Order, error) {
	p := &stripe.Order{}
	err := c.B.Call(http.MethodPost, "/v1/orders", c.Key, params, p)
	return p, err
}

// Update updates an order's properties.
func Update(id string, params *stripe.OrderUpdateParams) (*stripe.Order, error) {
	return getC().Update(id, params)
}

// Update updates an order's properties.
func (c Client) Update(id string, params *stripe.OrderUpdateParams) (*stripe.Order, error) {
	path := stripe.FormatURLPath("/v1/orders/%s", id)
	o := &stripe.Order{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, o)
	return o, err
}

// Pay pays an order.
func Pay(id string, params *stripe.OrderPayParams) (*stripe.Order, error) {
	return getC().Pay(id, params)
}

// Pay pays an order.
func (c Client) Pay(id string, params *stripe.OrderPayParams) (*stripe.Order, error) {
	path := stripe.FormatURLPath("/v1/orders/%s/pay", id)
	o := &stripe.Order{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, o)
	return o, err
}

// Get returns the details of an order.
func Get(id string, params *stripe.OrderParams) (*stripe.Order, error) {
	return getC().Get(id, params)
}

// Get returns the details of an order.
func (c Client) Get(id string, params *stripe.OrderParams) (*stripe.Order, error) {
	path := stripe.FormatURLPath("/v1/orders/%s", id)
	order := &stripe.Order{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, order)
	return order, err
}

// List returns a list of orders.
func List(params *stripe.OrderListParams) *Iter {
	return getC().List(params)
}

// List returns a list of orders.
func (c Client) List(listParams *stripe.OrderListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.OrderList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/orders", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Return returns all or part of an order.
func Return(id string, params *stripe.OrderReturnParams) (*stripe.OrderReturn, error) {
	return getC().Return(id, params)
}

// Return returns all or part of an order.
func (c Client) Return(id string, params *stripe.OrderReturnParams) (*stripe.OrderReturn, error) {
	path := stripe.FormatURLPath("/v1/orders/%s/returns", id)
	ret := &stripe.OrderReturn{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, ret)
	return ret, err
}

// Iter is an iterator for orders.
type Iter struct {
	*stripe.Iter
}

// Order returns the order which the iterator is currently pointing to.
func (i *Iter) Order() *stripe.Order {
	return i.Current().(*stripe.Order)
}
func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
