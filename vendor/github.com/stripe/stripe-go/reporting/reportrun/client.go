// Package reportrun provides API functions related to report runs.
//
// For more details, see: https://stripe.com/docs/api/go#reporting_report_run
package reportrun

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /reporting/report_runs APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new report run.
func New(params *stripe.ReportRunParams) (*stripe.ReportRun, error) {
	return getC().New(params)
}

// New creates a new report run.
func (c Client) New(params *stripe.ReportRunParams) (*stripe.ReportRun, error) {
	reportrun := &stripe.ReportRun{}
	err := c.B.Call(http.MethodPost, "/v1/reporting/report_runs", c.Key, params, reportrun)
	return reportrun, err
}

// Get returns the details of a report run.
func Get(id string, params *stripe.ReportRunParams) (*stripe.ReportRun, error) {
	return getC().Get(id, params)
}

// Get returns the details of a report run.
func (c Client) Get(id string, params *stripe.ReportRunParams) (*stripe.ReportRun, error) {
	path := stripe.FormatURLPath("/v1/reporting/report_runs/%s", id)
	reportrun := &stripe.ReportRun{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, reportrun)
	return reportrun, err
}

// List returns a list of report runs.
func List(params *stripe.ReportRunListParams) *Iter {
	return getC().List(params)
}

// List returns a list of report runs.
func (c Client) List(listParams *stripe.ReportRunListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.ReportRunList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/reporting/report_runs", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for report runs.
type Iter struct {
	*stripe.Iter
}

// ReportRun returns the report run which the iterator is currently pointing to.
func (i *Iter) ReportRun() *stripe.ReportRun {
	return i.Current().(*stripe.ReportRun)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
