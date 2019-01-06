// Package recipient provides the /recipients APIs
package recipient

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /recipients APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Since API version 2017-04-06, new recipients can't be created anymore.
// For that reason, there isn't a New() method for the Recipient resource.

// Get returns the details of a recipient.
func Get(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return getC().Get(id, params)
}

// Get returns the details of a recipient.
func (c Client) Get(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	path := stripe.FormatURLPath("/v1/recipients/%s", id)
	recipient := &stripe.Recipient{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, recipient)
	return recipient, err
}

// Update updates a recipient's properties.
func Update(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return getC().Update(id, params)
}

// Update updates a recipient's properties.
func (c Client) Update(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	path := stripe.FormatURLPath("/v1/recipients/%s", id)
	recipient := &stripe.Recipient{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, recipient)
	return recipient, err
}

// Del removes a recipient.
func Del(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return getC().Del(id, params)
}

// Del removes a recipient.
func (c Client) Del(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	path := stripe.FormatURLPath("/v1/recipients/%s", id)
	recipient := &stripe.Recipient{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, recipient)
	return recipient, err
}

// List returns a list of recipients.
func List(params *stripe.RecipientListParams) *Iter {
	return getC().List(params)
}

// List returns a list of recipients.
func (c Client) List(listParams *stripe.RecipientListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.RecipientList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/recipients", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for recipients.
type Iter struct {
	*stripe.Iter
}

// Recipient returns the recipient which the iterator is currently pointing to.
func (i *Iter) Recipient() *stripe.Recipient {
	return i.Current().(*stripe.Recipient)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
