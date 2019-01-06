// Package applepaydomain provides the /apple_pay/domains APIs
package applepaydomain

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /apple_pay/domains and Apple Pay domain-related APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new Apple Pay domain.
func New(params *stripe.ApplePayDomainParams) (*stripe.ApplePayDomain, error) {
	return getC().New(params)
}

// New creates a new Apple Pay domain.
func (c Client) New(params *stripe.ApplePayDomainParams) (*stripe.ApplePayDomain, error) {
	domain := &stripe.ApplePayDomain{}
	err := c.B.Call(http.MethodPost, "/v1/apple_pay/domains", c.Key, params, domain)
	return domain, err
}

// Get retrieves an Apple Pay domain.
func Get(id string, params *stripe.ApplePayDomainParams) (*stripe.ApplePayDomain, error) {
	return getC().Get(id, params)
}

// Get retrieves an Apple Pay domain.
func (c Client) Get(id string, params *stripe.ApplePayDomainParams) (*stripe.ApplePayDomain, error) {
	path := stripe.FormatURLPath("/v1/apple_pay/domains/%s", id)
	domain := &stripe.ApplePayDomain{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, domain)
	return domain, err
}

// Del removes an Apple Pay domain.
func Del(id string, params *stripe.ApplePayDomainParams) (*stripe.ApplePayDomain, error) {
	return getC().Del(id, params)
}

// Del removes an Apple Pay domain.
func (c Client) Del(id string, params *stripe.ApplePayDomainParams) (*stripe.ApplePayDomain, error) {
	path := stripe.FormatURLPath("/v1/apple_pay/domains/%s", id)
	domain := &stripe.ApplePayDomain{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, domain)
	return domain, err
}

// List lists available Apple Pay domains.
func List(params *stripe.ApplePayDomainListParams) *Iter {
	return getC().List(params)
}

// List lists available Apple Pay domains.
func (c Client) List(listParams *stripe.ApplePayDomainListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.ApplePayDomainList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/apple_pay/domains", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for Apple Pay domains.
type Iter struct {
	*stripe.Iter
}

// ApplePayDomain returns the Apple Pay domain which the iterator is currently pointing to.
func (i *Iter) ApplePayDomain() *stripe.ApplePayDomain {
	return i.Current().(*stripe.ApplePayDomain)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
