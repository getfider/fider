// Package bankaccount provides the /bank_accounts APIs
package bankaccount

import (
	"errors"
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /bank_accounts APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new bank account.
func New(params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return getC().New(params)
}

// New creates a new bank account.
func (c Client) New(params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	var path string
	if params.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources", stripe.StringValue(params.Customer))
	} else if params.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts", stripe.StringValue(params.Account))
	} else {
		return nil, errors.New("Invalid bank account params: either Customer or Account need to be set")
	}

	body := &form.Values{}

	// Note that we call this special append method instead of the standard one
	// from the form package. We should not use form's because doing so will
	// include some parameters that are undesirable here.
	params.AppendToAsSourceOrExternalAccount(body)

	// Because bank account creation uses the custom append above, we have to
	// make an explicit call using a form and CallRaw instead of the standard
	// Call (which takes a set of parameters).
	ba := &stripe.BankAccount{}
	err := c.B.CallRaw(http.MethodPost, path, c.Key, body, &params.Params, ba)
	return ba, err
}

// Get returns the details of a bank account.
func Get(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return getC().Get(id, params)
}

// Get returns the details of a bank account.
func (c Client) Get(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	var path string
	if params != nil && params.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources/%s", stripe.StringValue(params.Customer), id)
	} else if params != nil && params.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts/%s", stripe.StringValue(params.Account), id)
	} else {
		return nil, errors.New("Invalid bank account params: either Customer or Account need to be set")
	}

	ba := &stripe.BankAccount{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, ba)
	return ba, err
}

// Update updates a bank account.
func Update(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return getC().Update(id, params)
}

// Update updates a bank account.
func (c Client) Update(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	var path string
	if params.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources/%s", stripe.StringValue(params.Customer), id)
	} else if params.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts/%s", stripe.StringValue(params.Account), id)
	} else {
		return nil, errors.New("Invalid bank account params: either Customer or Account need to be set")
	}

	ba := &stripe.BankAccount{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, ba)
	return ba, err
}

// Del removes a bank account.
func Del(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return getC().Del(id, params)
}

// Del removes a bank account.
func (c Client) Del(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	var path string
	if params.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources/%s", stripe.StringValue(params.Customer), id)
	} else if params.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts/%s", stripe.StringValue(params.Account), id)
	} else {
		return nil, errors.New("Invalid bank account params: either Customer or Account need to be set")
	}

	ba := &stripe.BankAccount{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, ba)
	return ba, err
}

// List returns a list of bank accounts.
func List(params *stripe.BankAccountListParams) *Iter {
	return getC().List(params)
}

// List returns a list of bank accounts.
func (c Client) List(listParams *stripe.BankAccountListParams) *Iter {
	var path string
	var outerErr error

	// There's no bank accounts list URL, so we use one sources or external
	// accounts. An override on BankAccountListParam's `AppendTo` will add the
	// filter `object=bank_account` to make sure that only bank accounts come
	// back with the response.
	if listParams == nil {
		outerErr = errors.New("params should not be nil")
	} else if listParams.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources",
			stripe.StringValue(listParams.Customer))
	} else if listParams.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts",
			stripe.StringValue(listParams.Account))
	} else {
		outerErr = errors.New("Invalid bank account params: either Customer or Account need to be set")
	}

	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.BankAccountList{}

		if outerErr != nil {
			return nil, list.ListMeta, outerErr
		}

		err := c.B.CallRaw(http.MethodGet, path, c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for bank accounts.
type Iter struct {
	*stripe.Iter
}

// BankAccount returns the bank account which the iterator is currently pointing to.
func (i *Iter) BankAccount() *stripe.BankAccount {
	return i.Current().(*stripe.BankAccount)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
