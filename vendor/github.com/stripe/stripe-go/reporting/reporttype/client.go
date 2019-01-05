// Package reporttype provides API functions related to report types.
//
// For more details, see: https://stripe.com/docs/api/go#reporting_report_type
package reporttype

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /reporting/report_types APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns the details of a report type.
func Get(id string, params *stripe.ReportTypeParams) (*stripe.ReportType, error) {
	return getC().Get(id, params)
}

// Get returns the details of a report type.
func (c Client) Get(id string, params *stripe.ReportTypeParams) (*stripe.ReportType, error) {
	path := stripe.FormatURLPath("/v1/reporting/report_types/%s", id)
	reporttype := &stripe.ReportType{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, reporttype)
	return reporttype, err
}

// List returns a list of report types.
func List(params *stripe.ReportTypeListParams) *Iter {
	return getC().List(params)
}

// List returns a list of report types.
func (c Client) List(listParams *stripe.ReportTypeListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.ReportTypeList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/reporting/report_types", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for report types.
type Iter struct {
	*stripe.Iter
}

// ReportType returns the report type which the iterator is currently pointing to.
func (i *Iter) ReportType() *stripe.ReportType {
	return i.Current().(*stripe.ReportType)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
