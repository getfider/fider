// Package reader provides API functions related to terminal readers
package reader

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /terminal/readers APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new terminal reader.
func New(params *stripe.TerminalReaderParams) (*stripe.TerminalReader, error) {
	return getC().New(params)
}

// New creates a new terminal reader.
func (c Client) New(params *stripe.TerminalReaderParams) (*stripe.TerminalReader, error) {
	reader := &stripe.TerminalReader{}
	err := c.B.Call(http.MethodPost, "/v1/terminal/readers", c.Key, params, reader)
	return reader, err
}

// Get returns the details of a terminal reader.
func Get(id string, params *stripe.TerminalReaderGetParams) (*stripe.TerminalReader, error) {
	return getC().Get(id, params)
}

// Get returns the details of a terminal reader.
func (c Client) Get(id string, params *stripe.TerminalReaderGetParams) (*stripe.TerminalReader, error) {
	path := stripe.FormatURLPath("/v1/terminal/readers/%s", id)
	reader := &stripe.TerminalReader{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, reader)
	return reader, err
}

// Update updates a terminal reader.
func Update(id string, params *stripe.TerminalReaderParams) (*stripe.TerminalReader, error) {
	return getC().Update(id, params)
}

// Update updates a terminal reader.
func (c Client) Update(id string, params *stripe.TerminalReaderParams) (*stripe.TerminalReader, error) {
	path := stripe.FormatURLPath("/v1/terminal/readers/%s", id)
	reader := &stripe.TerminalReader{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, reader)
	return reader, err
}

// List returns a list of terminal readers.
func List(params *stripe.TerminalReaderListParams) *Iter {
	return getC().List(params)
}

// List returns a list of terminal readers.
func (c Client) List(listParams *stripe.TerminalReaderListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.TerminalReaderList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/terminal/readers", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for terminal readers.
type Iter struct {
	*stripe.Iter
}

// TerminalReader returns the terminal reader which the iterator is currently pointing to.
func (i *Iter) TerminalReader() *stripe.TerminalReader {
	return i.Current().(*stripe.TerminalReader)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
