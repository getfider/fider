// Package invoice provides the /invoices APIs
package invoice

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is the client used to invoke /invoices APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new invoice.
func New(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().New(params)
}

// New creates a new invoice.
func (c Client) New(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodPost, "/v1/invoices", c.Key, params, invoice)
	return invoice, err
}

// Del deletes an invoice.
func Del(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().Del(id, params)
}

// Del deletes an invoice.
func (c Client) Del(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	path := stripe.FormatURLPath("/v1/invoices/%s", id)
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, invoice)
	return invoice, err
}

// Get returns the details of an invoice.
func Get(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().Get(id, params)
}

// Get returns the details of an invoice.
func (c Client) Get(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	path := stripe.FormatURLPath("/v1/invoices/%s", id)
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, invoice)
	return invoice, err
}

// Pay pays an invoice.
func Pay(id string, params *stripe.InvoicePayParams) (*stripe.Invoice, error) {
	return getC().Pay(id, params)
}

// Pay pays an invoice.
func (c Client) Pay(id string, params *stripe.InvoicePayParams) (*stripe.Invoice, error) {
	path := stripe.FormatURLPath("/v1/invoices/%s/pay", id)
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, invoice)
	return invoice, err
}

// Update updates an invoice.
func Update(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().Update(id, params)
}

// Update updates an invoice.
func (c Client) Update(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	path := stripe.FormatURLPath("/v1/invoices/%s", id)
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, invoice)
	return invoice, err
}

// GetNext returns an upcoming invoice.
func GetNext(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().GetNext(params)
}

// GetNext returns an upcoming invoice.
func (c Client) GetNext(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodGet, "/v1/invoices/upcoming", c.Key, params, invoice)
	return invoice, err
}

// List returns a list of invoices.
func List(params *stripe.InvoiceListParams) *Iter {
	return getC().List(params)
}

// List returns a list of invoices.
func (c Client) List(listParams *stripe.InvoiceListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.InvoiceList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/invoices", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// ListLines returns a list of line items on an invoice.
func ListLines(params *stripe.InvoiceLineListParams) *LineIter {
	return getC().ListLines(params)
}

// ListLines returns a list of line items on an invoice.
func (c Client) ListLines(listParams *stripe.InvoiceLineListParams) *LineIter {
	path := stripe.FormatURLPath("/v1/invoices/%s/lines", stripe.StringValue(listParams.ID))
	return &LineIter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.InvoiceLineList{}
		err := c.B.CallRaw(http.MethodGet, path, c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// FinalizeInvoice finalizes an invoice.
func FinalizeInvoice(id string, params *stripe.InvoiceFinalizeParams) (*stripe.Invoice, error) {
	return getC().FinalizeInvoice(id, params)
}

// FinalizeInvoice finalizes an invoice.
func (c Client) FinalizeInvoice(id string, params *stripe.InvoiceFinalizeParams) (*stripe.Invoice, error) {
	path := stripe.FormatURLPath("/v1/invoices/%s/finalize", id)
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, invoice)
	return invoice, err
}

// MarkUncollectible marks an invoice as uncollectible.
func MarkUncollectible(id string, params *stripe.InvoiceMarkUncollectibleParams) (*stripe.Invoice, error) {
	return getC().MarkUncollectible(id, params)
}

// MarkUncollectible marks an invoice as uncollectible.
func (c Client) MarkUncollectible(id string, params *stripe.InvoiceMarkUncollectibleParams) (*stripe.Invoice, error) {
	path := stripe.FormatURLPath("/v1/invoices/%s/mark_uncollectible", id)
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, invoice)
	return invoice, err
}

// SendInvoice sends an invoice.
func SendInvoice(id string, params *stripe.InvoiceSendParams) (*stripe.Invoice, error) {
	return getC().SendInvoice(id, params)
}

// SendInvoice sends an invoice.
func (c Client) SendInvoice(id string, params *stripe.InvoiceSendParams) (*stripe.Invoice, error) {
	path := stripe.FormatURLPath("/v1/invoices/%s/send", id)
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, invoice)
	return invoice, err
}

// VoidInvoice voids an invoice.
func VoidInvoice(id string, params *stripe.InvoiceVoidParams) (*stripe.Invoice, error) {
	return getC().VoidInvoice(id, params)
}

// VoidInvoice voids an invoice.
func (c Client) VoidInvoice(id string, params *stripe.InvoiceVoidParams) (*stripe.Invoice, error) {
	path := stripe.FormatURLPath("/v1/invoices/%s/void", id)
	invoice := &stripe.Invoice{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, invoice)
	return invoice, err
}

// Iter is an iterator for invoices.
type Iter struct {
	*stripe.Iter
}

// Invoice returns the invoice which the iterator is currently pointing to.
func (i *Iter) Invoice() *stripe.Invoice {
	return i.Current().(*stripe.Invoice)
}

// LineIter is an iterator for line items on an invoice.
type LineIter struct {
	*stripe.Iter
}

// InvoiceLine returns the line item which the iterator is currently pointing to.
func (i *LineIter) InvoiceLine() *stripe.InvoiceLine {
	return i.Current().(*stripe.InvoiceLine)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
