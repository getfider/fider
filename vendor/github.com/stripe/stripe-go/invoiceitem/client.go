// Package invoiceitem provides the /invoiceitems APIs
package invoiceitem

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /invoiceitems APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new invoice item.
func New(params *stripe.InvoiceItemParams) (*stripe.InvoiceItem, error) {
	return getC().New(params)
}

// New creates a new invoice item.
func (c Client) New(params *stripe.InvoiceItemParams) (*stripe.InvoiceItem, error) {
	invoiceItem := &stripe.InvoiceItem{}
	err := c.B.Call(http.MethodPost, "/v1/invoiceitems", c.Key, params, invoiceItem)
	return invoiceItem, err
}

// Get returns the details of an invoice item.
func Get(id string, params *stripe.InvoiceItemParams) (*stripe.InvoiceItem, error) {
	return getC().Get(id, params)
}

// Get returns the details of an invoice item.
func (c Client) Get(id string, params *stripe.InvoiceItemParams) (*stripe.InvoiceItem, error) {
	path := stripe.FormatURLPath("/v1/invoiceitems/%s", id)
	invoiceItem := &stripe.InvoiceItem{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, invoiceItem)
	return invoiceItem, err
}

// Update updates an invoice item.
func Update(id string, params *stripe.InvoiceItemParams) (*stripe.InvoiceItem, error) {
	return getC().Update(id, params)
}

// Update updates an invoice item.
func (c Client) Update(id string, params *stripe.InvoiceItemParams) (*stripe.InvoiceItem, error) {
	path := stripe.FormatURLPath("/v1/invoiceitems/%s", id)
	invoiceItem := &stripe.InvoiceItem{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, invoiceItem)
	return invoiceItem, err
}

// Del removes an invoice item.
func Del(id string, params *stripe.InvoiceItemParams) (*stripe.InvoiceItem, error) {
	return getC().Del(id, params)
}

// Del removes an invoice item.
func (c Client) Del(id string, params *stripe.InvoiceItemParams) (*stripe.InvoiceItem, error) {
	path := stripe.FormatURLPath("/v1/invoiceitems/%s", id)
	ii := &stripe.InvoiceItem{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, ii)
	return ii, err
}

// List returns a list of invoice items.
func List(params *stripe.InvoiceItemListParams) *Iter {
	return getC().List(params)
}

// List returns a list of invoice items.
func (c Client) List(listParams *stripe.InvoiceItemListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.InvoiceItemList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/invoiceitems", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for invoice items.
type Iter struct {
	*stripe.Iter
}

// InvoiceItem returns the invoice item which the iterator is currently pointing to.
func (i *Iter) InvoiceItem() *stripe.InvoiceItem {
	return i.Current().(*stripe.InvoiceItem)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
