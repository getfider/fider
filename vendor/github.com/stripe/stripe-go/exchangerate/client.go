// Package exchangerate provides the /exchange_rates APIs
package exchangerate

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /exchange_rates and exchangerates-related APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns the exchante rate for a given currency.
func Get(currency string, params *stripe.ExchangeRateParams) (*stripe.ExchangeRate, error) {
	return getC().Get(currency, params)
}

// Get returns the exchante rate for a given currency.
func (c Client) Get(currency string, params *stripe.ExchangeRateParams) (*stripe.ExchangeRate, error) {
	path := stripe.FormatURLPath("/v1/exchange_rates/%s", currency)
	exchangeRate := &stripe.ExchangeRate{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, exchangeRate)

	return exchangeRate, err
}

// List lists available exchange rates.
func List(params *stripe.ExchangeRateListParams) *Iter {
	return getC().List(params)
}

// List lists available exchange rates.
func (c Client) List(listParams *stripe.ExchangeRateListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.ExchangeRateList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/exchange_rates", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for exchange rates.
type Iter struct {
	*stripe.Iter
}

// ExchangeRate returns the exchange rate which the iterator is currently pointing to.
func (i *Iter) ExchangeRate() *stripe.ExchangeRate {
	return i.Current().(*stripe.ExchangeRate)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
