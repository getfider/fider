// Package paymentintent provides API functions related to payment intents.
//
// For more details, see: https://stripe.com/docs/api/go#payment_intents.
package paymentintent

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke APIs related to payment intents.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a payment intent.
func New(params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	return getC().New(params)
}

// New creates a payment intent.
func (c Client) New(params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	intent := &stripe.PaymentIntent{}
	err := c.B.Call(http.MethodPost, "/v1/payment_intents", c.Key, params, intent)
	return intent, err
}

// Get retrieves a payment intent.
func Get(id string, params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	return getC().Get(id, params)
}

// Get retrieves a payment intent.
func (c Client) Get(id string, params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	path := stripe.FormatURLPath("/v1/payment_intents/%s", id)
	intent := &stripe.PaymentIntent{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, intent)
	return intent, err
}

// Update updates a payment intent.
func Update(id string, params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	return getC().Update(id, params)
}

// Update updates a payment intent.
func (c Client) Update(id string, params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	path := stripe.FormatURLPath("/v1/payment_intents/%s", id)
	intent := &stripe.PaymentIntent{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, intent)
	return intent, err
}

// Cancel cancels a payment intent.
func Cancel(id string, params *stripe.PaymentIntentCancelParams) (*stripe.PaymentIntent, error) {
	return getC().Cancel(id, params)
}

// Cancel cancels a payment intent.
func (c Client) Cancel(id string, params *stripe.PaymentIntentCancelParams) (*stripe.PaymentIntent, error) {
	path := stripe.FormatURLPath("/v1/payment_intents/%s/cancel", id)
	intent := &stripe.PaymentIntent{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, intent)
	return intent, err
}

// Capture captures a payment intent.
func Capture(id string, params *stripe.PaymentIntentCaptureParams) (*stripe.PaymentIntent, error) {
	return getC().Capture(id, params)
}

// Capture captures a payment intent.
func (c Client) Capture(id string, params *stripe.PaymentIntentCaptureParams) (*stripe.PaymentIntent, error) {
	path := stripe.FormatURLPath("/v1/payment_intents/%s/capture", id)
	intent := &stripe.PaymentIntent{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, intent)
	return intent, err
}

// Confirm confirms a payment intent.
func Confirm(id string, params *stripe.PaymentIntentConfirmParams) (*stripe.PaymentIntent, error) {
	return getC().Confirm(id, params)
}

// Confirm confirms a payment intent.
func (c Client) Confirm(id string, params *stripe.PaymentIntentConfirmParams) (*stripe.PaymentIntent, error) {
	path := stripe.FormatURLPath("/v1/payment_intents/%s/confirm", id)
	intent := &stripe.PaymentIntent{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, intent)
	return intent, err
}

// List returns a list of payment intents.
func List(params *stripe.PaymentIntentListParams) *Iter {
	return getC().List(params)
}

// List returns a list of payment intents.
func (c Client) List(listParams *stripe.PaymentIntentListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.PaymentIntentList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/payment_intents", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for payment intents.
type Iter struct {
	*stripe.Iter
}

// PaymentIntent returns the payment intent which the iterator is currently pointing to.
func (i *Iter) PaymentIntent() *stripe.PaymentIntent {
	return i.Current().(*stripe.PaymentIntent)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
