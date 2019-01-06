// Package filelink provides API functions related to file links.
//
// For more details, see: https://stripe.com/docs/api/go#file_links.
package filelink

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke APIs related to file links.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new file link.
func New(params *stripe.FileLinkParams) (*stripe.FileLink, error) {
	return getC().New(params)
}

// New creates a new file link.
func (c Client) New(params *stripe.FileLinkParams) (*stripe.FileLink, error) {
	fileLink := &stripe.FileLink{}
	err := c.B.Call(http.MethodPost, "/v1/file_links", c.Key, params, fileLink)
	return fileLink, err
}

// Get retrieves a file link.
func Get(id string, params *stripe.FileLinkParams) (*stripe.FileLink, error) {
	return getC().Get(id, params)
}

// Get retrieves a file link.
func (c Client) Get(id string, params *stripe.FileLinkParams) (*stripe.FileLink, error) {
	path := stripe.FormatURLPath("/v1/file_links/%s", id)
	fileLink := &stripe.FileLink{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, fileLink)
	return fileLink, err
}

// Update updates a file link.
func Update(id string, params *stripe.FileLinkParams) (*stripe.FileLink, error) {
	return getC().Update(id, params)
}

// Update updates a file link.
func (c Client) Update(id string, params *stripe.FileLinkParams) (*stripe.FileLink, error) {
	path := stripe.FormatURLPath("/v1/file_links/%s", id)
	fileLink := &stripe.FileLink{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, fileLink)
	return fileLink, err
}

// List returns an iterator that iterates all file links.
func List(params *stripe.FileLinkListParams) *Iter {
	return getC().List(params)
}

// List returns an iterator that iterates all file links.
func (c Client) List(listParams *stripe.FileLinkListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.FileLinkList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/file_links", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for file links.
type Iter struct {
	*stripe.Iter
}

// FileLink returns the file link which the iterator is currently pointing to.
func (i *Iter) FileLink() *stripe.FileLink {
	return i.Current().(*stripe.FileLink)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
