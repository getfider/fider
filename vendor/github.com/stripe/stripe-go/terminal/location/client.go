// Package location provides API functions related to terminal locations
package location

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invokte /terminal/locations APIs
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new terminal location.
func New(params *stripe.TerminalLocationParams) (*stripe.TerminalLocation, error) {
	return getC().New(params)
}

// New creates a new terminal location.
func (c Client) New(params *stripe.TerminalLocationParams) (*stripe.TerminalLocation, error) {
	location := &stripe.TerminalLocation{}
	err := c.B.Call(http.MethodPost, "/v1/terminal/locations", c.Key, params, location)
	return location, err
}

// Get returns the details of a terminal location.
func Get(id string, params *stripe.TerminalLocationParams) (*stripe.TerminalLocation, error) {
	return getC().Get(id, params)
}

// Get returns the details of a terminal location.
func (c Client) Get(id string, params *stripe.TerminalLocationParams) (*stripe.TerminalLocation, error) {
	path := stripe.FormatURLPath("/v1/terminal/locations/%s", id)
	location := &stripe.TerminalLocation{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, location)
	return location, err
}

// Update updates a terminal location.
func Update(id string, params *stripe.TerminalLocationParams) (*stripe.TerminalLocation, error) {
	return getC().Update(id, params)
}

// Update updates a terminal location.
func (c Client) Update(id string, params *stripe.TerminalLocationParams) (*stripe.TerminalLocation, error) {
	path := stripe.FormatURLPath("/v1/terminal/locations/%s", id)
	location := &stripe.TerminalLocation{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, location)
	return location, err
}

// List returns a list of terminal location.
func List(params *stripe.TerminalLocationListParams) *Iter {
	return getC().List(params)
}

// List returns a list of terminal location.
func (c Client) List(listParams *stripe.TerminalLocationListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.TerminalLocationList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/terminal/locations", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for terminal locations.
type Iter struct {
	*stripe.Iter
}

// TerminalLocation returns the terminal location which the iterator is currently pointing to.
func (i *Iter) TerminalLocation() *stripe.TerminalLocation {
	return i.Current().(*stripe.TerminalLocation)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
