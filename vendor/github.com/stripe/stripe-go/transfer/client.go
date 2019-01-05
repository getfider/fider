// Package transfer provides the /transfers APIs
package transfer

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /transfers APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new transfer.
func New(params *stripe.TransferParams) (*stripe.Transfer, error) {
	return getC().New(params)
}

// New creates a new transfer.
func (c Client) New(params *stripe.TransferParams) (*stripe.Transfer, error) {
	transfer := &stripe.Transfer{}
	err := c.B.Call(http.MethodPost, "/v1/transfers", c.Key, params, transfer)
	return transfer, err
}

// Get returns the details of a transfer.
func Get(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	return getC().Get(id, params)
}

// Get returns the details of a transfer.
func (c Client) Get(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	path := stripe.FormatURLPath("/v1/transfers/%s", id)
	transfer := &stripe.Transfer{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, transfer)
	return transfer, err
}

// Update updates a transfer's properties.
func Update(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	return getC().Update(id, params)
}

// Update updates a transfer's properties.
func (c Client) Update(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	path := stripe.FormatURLPath("/v1/transfers/%s", id)
	transfer := &stripe.Transfer{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, transfer)
	return transfer, err
}

// List returns a list of transfers.
func List(params *stripe.TransferListParams) *Iter {
	return getC().List(params)
}

// List returns a list of transfers.
func (c Client) List(listParams *stripe.TransferListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.TransferList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/transfers", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for transfers.
type Iter struct {
	*stripe.Iter
}

// Transfer returns the transfer which the iterator is currently pointing to.
func (i *Iter) Transfer() *stripe.Transfer {
	return i.Current().(*stripe.Transfer)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
