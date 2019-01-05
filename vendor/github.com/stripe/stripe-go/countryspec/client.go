// Package countryspec provides the /country_specs APIs
package countryspec

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /country_specs and countryspec-related APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns a Country Spec for a given country code.
func Get(country string, params *stripe.CountrySpecParams) (*stripe.CountrySpec, error) {
	return getC().Get(country, params)
}

// Get returns a Country Spec for a given country code.
func (c Client) Get(country string, params *stripe.CountrySpecParams) (*stripe.CountrySpec, error) {
	path := stripe.FormatURLPath("/v1/country_specs/%s", country)
	countrySpec := &stripe.CountrySpec{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, countrySpec)
	return countrySpec, err
}

// List lists available Country Specs.
func List(params *stripe.CountrySpecListParams) *Iter {
	return getC().List(params)
}

// List lists available Country Specs.
func (c Client) List(listParams *stripe.CountrySpecListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.CountrySpecList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/country_specs", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for Country Specs.
type Iter struct {
	*stripe.Iter
}

// CountrySpec returns the Country Spec which the iterator is currently pointing to.
func (i *Iter) CountrySpec() *stripe.CountrySpec {
	return i.Current().(*stripe.CountrySpec)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
