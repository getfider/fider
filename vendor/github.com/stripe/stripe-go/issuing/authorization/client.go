// Package authorization provides API functions related to issuing authorizations.
//
// For more details, see: https://stripe.com/docs/api/go#issuing_authorizations
package authorization

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /issuing/authorizations APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Approve approves an issuing authorization.
func Approve(id string, params *stripe.IssuingAuthorizationParams) (*stripe.IssuingAuthorization, error) {
	return getC().Approve(id, params)
}

// Approve updates an issuing authorization.
func (c Client) Approve(id string, params *stripe.IssuingAuthorizationParams) (*stripe.IssuingAuthorization, error) {
	path := stripe.FormatURLPath("/v1/issuing/authorizations/%s/approve", id)
	authorization := &stripe.IssuingAuthorization{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, authorization)
	return authorization, err
}

// Decline decline an issuing authorization.
func Decline(id string, params *stripe.IssuingAuthorizationParams) (*stripe.IssuingAuthorization, error) {
	return getC().Decline(id, params)
}

// Decline updates an issuing authorization.
func (c Client) Decline(id string, params *stripe.IssuingAuthorizationParams) (*stripe.IssuingAuthorization, error) {
	path := stripe.FormatURLPath("/v1/issuing/authorizations/%s/decline", id)
	authorization := &stripe.IssuingAuthorization{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, authorization)
	return authorization, err
}

// Get returns the details of an issuing authorization.
func Get(id string, params *stripe.IssuingAuthorizationParams) (*stripe.IssuingAuthorization, error) {
	return getC().Get(id, params)
}

// Get returns the details of an issuing authorization.
func (c Client) Get(id string, params *stripe.IssuingAuthorizationParams) (*stripe.IssuingAuthorization, error) {
	path := stripe.FormatURLPath("/v1/issuing/authorizations/%s", id)
	authorization := &stripe.IssuingAuthorization{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, authorization)
	return authorization, err
}

// Update updates an issuing authorization.
func Update(id string, params *stripe.IssuingAuthorizationParams) (*stripe.IssuingAuthorization, error) {
	return getC().Update(id, params)
}

// Update updates an issuing authorization.
func (c Client) Update(id string, params *stripe.IssuingAuthorizationParams) (*stripe.IssuingAuthorization, error) {
	path := stripe.FormatURLPath("/v1/issuing/authorizations/%s", id)
	authorization := &stripe.IssuingAuthorization{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, authorization)
	return authorization, err
}

// List returns a list of issuing authorizations.
func List(params *stripe.IssuingAuthorizationListParams) *Iter {
	return getC().List(params)
}

// List returns a list of issuing authorizations.
func (c Client) List(listParams *stripe.IssuingAuthorizationListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.IssuingAuthorizationList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/issuing/authorizations", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for issuing authorizations.
type Iter struct {
	*stripe.Iter
}

// IssuingAuthorization returns the issuing authorization which the iterator is currently pointing to.
func (i *Iter) IssuingAuthorization() *stripe.IssuingAuthorization {
	return i.Current().(*stripe.IssuingAuthorization)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
