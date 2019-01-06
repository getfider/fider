package dispute

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke dispute-related APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns the details of a dispute.
func Get(id string, params *stripe.DisputeParams) (*stripe.Dispute, error) {
	return getC().Get(id, params)
}

// Get returns the details of a dispute.
func (c Client) Get(id string, params *stripe.DisputeParams) (*stripe.Dispute, error) {
	path := stripe.FormatURLPath("/v1/disputes/%s", id)
	dispute := &stripe.Dispute{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, dispute)
	return dispute, err
}

// List returns a list of disputes.
func List(params *stripe.DisputeListParams) *Iter {
	return getC().List(params)
}

// List returns a list of disputes.
func (c Client) List(listParams *stripe.DisputeListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.DisputeList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/disputes", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Update updates a dispute.
func Update(id string, params *stripe.DisputeParams) (*stripe.Dispute, error) {
	return getC().Update(id, params)
}

// Update updates a dispute.
func (c Client) Update(id string, params *stripe.DisputeParams) (*stripe.Dispute, error) {
	path := stripe.FormatURLPath("/v1/disputes/%s", id)
	dispute := &stripe.Dispute{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, dispute)
	return dispute, err
}

// Close dismisses a dispute in the customer's favor.
func Close(id string, params *stripe.DisputeParams) (*stripe.Dispute, error) {
	return getC().Close(id, params)
}

// Close dismisses a dispute in the customer's favor.
func (c Client) Close(id string, params *stripe.DisputeParams) (*stripe.Dispute, error) {
	path := stripe.FormatURLPath("/v1/disputes/%s/close", id)
	dispute := &stripe.Dispute{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, dispute)
	return dispute, err
}

// Iter is an iterator for disputes.
type Iter struct {
	*stripe.Iter
}

// Dispute returns the dispute which the iterator is currently pointing to.
func (i *Iter) Dispute() *stripe.Dispute {
	return i.Current().(*stripe.Dispute)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
