// Package feerefund provides the /application_fees/refunds APIs
package feerefund

import (
	"fmt"
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /application_fees/refunds APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates an application fee refund.
func New(params *stripe.FeeRefundParams) (*stripe.FeeRefund, error) {
	return getC().New(params)
}

// New creates an application fee refund.
func (c Client) New(params *stripe.FeeRefundParams) (*stripe.FeeRefund, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ApplicationFee == nil {
		return nil, fmt.Errorf("params.ApplicationFee must be set")
	}

	path := stripe.FormatURLPath("/v1/application_fees/%s/refunds",
		stripe.StringValue(params.ApplicationFee))
	refund := &stripe.FeeRefund{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, refund)
	return refund, err
}

// Get returns the details of an application fee refund.
func Get(id string, params *stripe.FeeRefundParams) (*stripe.FeeRefund, error) {
	return getC().Get(id, params)
}

// Get returns the details of an application fee refund.
func (c Client) Get(id string, params *stripe.FeeRefundParams) (*stripe.FeeRefund, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ApplicationFee == nil {
		return nil, fmt.Errorf("params.ApplicationFee must be set")
	}

	path := stripe.FormatURLPath("/v1/application_fees/%s/refunds/%s",
		stripe.StringValue(params.ApplicationFee), id)
	refund := &stripe.FeeRefund{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, refund)
	return refund, err
}

// Update updates an application fee refund.
func Update(id string, params *stripe.FeeRefundParams) (*stripe.FeeRefund, error) {
	return getC().Update(id, params)
}

// Update updates an application fee refund.
func (c Client) Update(id string, params *stripe.FeeRefundParams) (*stripe.FeeRefund, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ApplicationFee == nil {
		return nil, fmt.Errorf("params.ApplicationFee must be set")
	}

	path := stripe.FormatURLPath("/v1/application_fees/%s/refunds/%s",
		stripe.StringValue(params.ApplicationFee), id)
	refund := &stripe.FeeRefund{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, refund)

	return refund, err
}

// List returns a list of application fee refunds.
func List(params *stripe.FeeRefundListParams) *Iter {
	return getC().List(params)
}

// List returns a list of application fee refunds.
func (c Client) List(listParams *stripe.FeeRefundListParams) *Iter {
	path := stripe.FormatURLPath("/v1/application_fees/%s/refunds",
		stripe.StringValue(listParams.ApplicationFee))

	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.FeeRefundList{}
		err := c.B.CallRaw(http.MethodGet, path, c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for application fee refunds.
type Iter struct {
	*stripe.Iter
}

// FeeRefund returns the application fee refund which the iterator is currently pointing to.
func (i *Iter) FeeRefund() *stripe.FeeRefund {
	return i.Current().(*stripe.FeeRefund)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
