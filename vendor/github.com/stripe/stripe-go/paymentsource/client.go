// Package paymentsource provides the /sources APIs
package paymentsource

import (
	"errors"
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /sources APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new source for a customer.
func New(params *stripe.CustomerSourceParams) (*stripe.PaymentSource, error) {
	return getC().New(params)
}

// New creates a new source for a customer.
func (s Client) New(params *stripe.CustomerSourceParams) (*stripe.PaymentSource, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	if params.Customer == nil {
		return nil, errors.New("Invalid source params: customer needs to be set")
	}

	path := stripe.FormatURLPath("/v1/customers/%s/sources", stripe.StringValue(params.Customer))
	source := &stripe.PaymentSource{}
	err := s.B.Call(http.MethodPost, path, s.Key, params, source)
	return source, err
}

// Get returns the details of a source.
func Get(id string, params *stripe.CustomerSourceParams) (*stripe.PaymentSource, error) {
	return getC().Get(id, params)
}

// Get returns the details of a source.
func (s Client) Get(id string, params *stripe.CustomerSourceParams) (*stripe.PaymentSource, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	if params.Customer == nil {
		return nil, errors.New("Invalid source params: customer needs to be set")
	}

	path := stripe.FormatURLPath("/v1/customers/%s/sources/%s", stripe.StringValue(params.Customer), id)
	source := &stripe.PaymentSource{}
	err := s.B.Call(http.MethodGet, path, s.Key, params, source)
	return source, err
}

// Update updates a source's properties.
func Update(id string, params *stripe.CustomerSourceParams) (*stripe.PaymentSource, error) {
	return getC().Update(id, params)
}

// Update updates a source's properties.
func (s Client) Update(id string, params *stripe.CustomerSourceParams) (*stripe.PaymentSource, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	if params.Customer == nil {
		return nil, errors.New("Invalid source params: customer needs to be set")
	}

	path := stripe.FormatURLPath("/v1/customers/%s/sources/%s", stripe.StringValue(params.Customer), id)
	source := &stripe.PaymentSource{}
	err := s.B.Call(http.MethodPost, path, s.Key, params, source)
	return source, err
}

// Del removes a source.
func Del(id string, params *stripe.CustomerSourceParams) (*stripe.PaymentSource, error) {
	return getC().Del(id, params)
}

// Del removes a source.
func (s Client) Del(id string, params *stripe.CustomerSourceParams) (*stripe.PaymentSource, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	if params.Customer == nil {
		return nil, errors.New("Invalid source params: customer needs to be set")
	}

	source := &stripe.PaymentSource{}
	path := stripe.FormatURLPath("/v1/customers/%s/sources/%s", stripe.StringValue(params.Customer), id)
	err := s.B.Call(http.MethodDelete, path, s.Key, params, source)
	return source, err
}

// List returns a list of sources.
func List(params *stripe.SourceListParams) *Iter {
	return getC().List(params)
}

// List returns a list of sources.
func (s Client) List(listParams *stripe.SourceListParams) *Iter {
	var outerErr error
	var path string

	if listParams == nil {
		outerErr = errors.New("params should not be nil")
	} else if listParams.Customer == nil {
		outerErr = errors.New("Invalid source params: customer needs to be set")
	} else {
		path = stripe.FormatURLPath("/v1/customers/%s/sources",
			stripe.StringValue(listParams.Customer))
	}

	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.SourceList{}

		if outerErr != nil {
			return nil, list.ListMeta, outerErr
		}

		err := s.B.CallRaw(http.MethodGet, path, s.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Verify verifies a source which is used for bank accounts.
func Verify(id string, params *stripe.SourceVerifyParams) (*stripe.PaymentSource, error) {
	return getC().Verify(id, params)
}

// Verify verifies a source which is used for bank accounts.
func (s Client) Verify(id string, params *stripe.SourceVerifyParams) (*stripe.PaymentSource, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	var path string
	if params.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources/%s/verify",
			stripe.StringValue(params.Customer), id)
	} else if len(params.Values) > 0 {
		path = stripe.FormatURLPath("/v1/sources/%s/verify", id)
	} else {
		return nil, errors.New("Only customer bank accounts or sources can be verified in this manner")
	}

	source := &stripe.PaymentSource{}
	err := s.B.Call(http.MethodPost, path, s.Key, params, source)
	return source, err
}

// Iter is an iterator for sources.
type Iter struct {
	*stripe.Iter
}

// PaymentSource returns the source which the iterator is currently pointing to.
func (i *Iter) PaymentSource() *stripe.PaymentSource {
	return i.Current().(*stripe.PaymentSource)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
