package issuerfraudrecord

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to interact with the /issuer_fraud_records API.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns the details of an issuer fraud record.
func Get(id string, params *stripe.IssuerFraudRecordParams) (*stripe.IssuerFraudRecord, error) {
	return getC().Get(id, params)
}

// Get returns the details of an issuer fraud record.
func (c Client) Get(id string, params *stripe.IssuerFraudRecordParams) (*stripe.IssuerFraudRecord, error) {
	path := stripe.FormatURLPath("/v1/issuer_fraud_records/%s", id)
	ifr := &stripe.IssuerFraudRecord{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, ifr)
	return ifr, err
}

// List returns a list of issuer fraud records.
func List(params *stripe.IssuerFraudRecordListParams) *Iter {
	return getC().List(params)
}

// List returns a list of issuer fraud records.
func (c Client) List(listParams *stripe.IssuerFraudRecordListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.IssuerFraudRecordList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/issuer_fraud_records", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Values))
		for i, v := range list.Values {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for issuer fraud records.
type Iter struct {
	*stripe.Iter
}

// IssuerFraudRecord returns the issuer fraud record which the iterator is currently pointing to.
func (i *Iter) IssuerFraudRecord() *stripe.IssuerFraudRecord {
	return i.Current().(*stripe.IssuerFraudRecord)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
