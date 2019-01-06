// Package scheduledqueryrun provides API functions related to scheduled query runs.
//
// For more details, see: https://stripe.com/docs/api#scheduled_queries
package scheduledqueryrun

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /sigma/scheduled_query_runs APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns the details of an scheduled query run.
func Get(id string, params *stripe.SigmaScheduledQueryRunParams) (*stripe.SigmaScheduledQueryRun, error) {
	return getC().Get(id, params)
}

// Get returns the details of an scheduled query run.
func (c Client) Get(id string, params *stripe.SigmaScheduledQueryRunParams) (*stripe.SigmaScheduledQueryRun, error) {
	path := stripe.FormatURLPath("/v1/sigma/scheduled_query_runs/%s", id)
	run := &stripe.SigmaScheduledQueryRun{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, run)
	return run, err
}

// List returns a list of scheduled query runs.
func List(params *stripe.SigmaScheduledQueryRunListParams) *Iter {
	return getC().List(params)
}

// List returns a list of scheduled query runs.
func (c Client) List(listParams *stripe.SigmaScheduledQueryRunListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.SigmaScheduledQueryRunList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/sigma/scheduled_query_runs", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for scheduled query runs.
type Iter struct {
	*stripe.Iter
}

// SigmaScheduledQueryRun returns the scheduled query run which the iterator is currently pointing to.
func (i *Iter) SigmaScheduledQueryRun() *stripe.SigmaScheduledQueryRun {
	return i.Current().(*stripe.SigmaScheduledQueryRun)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
