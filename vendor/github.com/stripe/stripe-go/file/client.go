// Package file provides the file related APIs
package file

import (
	"fmt"
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke file APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new file.
func New(params *stripe.FileParams) (*stripe.File, error) {
	return getC().New(params)
}

// New creates a new file.
func (c Client) New(params *stripe.FileParams) (*stripe.File, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil, and params.Purpose and params.File must be set")
	}

	bodyBuffer, boundary, err := params.GetBody()
	if err != nil {
		return nil, err
	}

	file := &stripe.File{}
	err = c.B.CallMultipart(http.MethodPost, "/v1/files", c.Key, boundary, bodyBuffer, &params.Params, file)

	return file, err
}

// Get returns the details of a file.
func Get(id string, params *stripe.FileParams) (*stripe.File, error) {
	return getC().Get(id, params)

}

// Get returns the details of a file.
func (c Client) Get(id string, params *stripe.FileParams) (*stripe.File, error) {
	path := stripe.FormatURLPath("/v1/files/%s", id)
	file := &stripe.File{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, file)
	return file, err
}

// List returns a list of files.
func List(params *stripe.FileListParams) *Iter {
	return getC().List(params)
}

// List returns a list of files.
func (c Client) List(listParams *stripe.FileListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.FileList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/files", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for files.
type Iter struct {
	*stripe.Iter
}

// File returns the file which the iterator is currently pointing to.
func (i *Iter) File() *stripe.File {
	return i.Current().(*stripe.File)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.UploadsBackend), stripe.Key}
}
