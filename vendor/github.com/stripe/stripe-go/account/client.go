// Package account provides API functions related to accounts.
//
// For more details, see: https://stripe.com/docs/api/go#accounts.
package account

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke APIs related to accounts.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new account.
func New(params *stripe.AccountParams) (*stripe.Account, error) {
	return getC().New(params)
}

// New creates a new account.
func (c Client) New(params *stripe.AccountParams) (*stripe.Account, error) {
	acct := &stripe.Account{}
	err := c.B.Call(http.MethodPost, "/v1/accounts", c.Key, params, acct)
	return acct, err
}

// Get retrieves the authenticating account.
func Get() (*stripe.Account, error) {
	return getC().Get()
}

// Get retrieves the authenticating account.
func (c Client) Get() (*stripe.Account, error) {
	account := &stripe.Account{}
	err := c.B.Call(http.MethodGet, "/v1/account", c.Key, nil, account)
	return account, err
}

// GetByID retrieves an account.
func GetByID(id string, params *stripe.AccountParams) (*stripe.Account, error) {
	return getC().GetByID(id, params)
}

// GetByID retrieves an account.
func (c Client) GetByID(id string, params *stripe.AccountParams) (*stripe.Account, error) {
	path := stripe.FormatURLPath("/v1/accounts/%s", id)
	account := &stripe.Account{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, account)
	return account, err
}

// Update updates an account.
func Update(id string, params *stripe.AccountParams) (*stripe.Account, error) {
	return getC().Update(id, params)
}

// Update updates an account.
func (c Client) Update(id string, params *stripe.AccountParams) (*stripe.Account, error) {
	path := stripe.FormatURLPath("/v1/accounts/%s", id)
	acct := &stripe.Account{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, acct)
	return acct, err
}

// Del deletes an account.
func Del(id string, params *stripe.AccountParams) (*stripe.Account, error) {
	return getC().Del(id, params)
}

// Del deletes an account.
func (c Client) Del(id string, params *stripe.AccountParams) (*stripe.Account, error) {
	path := stripe.FormatURLPath("/v1/accounts/%s", id)
	acct := &stripe.Account{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, acct)
	return acct, err
}

// Reject rejects an account.
func Reject(id string, params *stripe.AccountRejectParams) (*stripe.Account, error) {
	return getC().Reject(id, params)
}

// Reject rejects an account.
func (c Client) Reject(id string, params *stripe.AccountRejectParams) (*stripe.Account, error) {
	path := stripe.FormatURLPath("/v1/accounts/%s/reject", id)
	acct := &stripe.Account{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, acct)
	return acct, err
}

// List returns an iterator that iterates all accounts.
func List(params *stripe.AccountListParams) *Iter {
	return getC().List(params)
}

// List returns an iterator that iterates all accounts.
func (c Client) List(listParams *stripe.AccountListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.AccountList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/accounts", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for accounts.
type Iter struct {
	*stripe.Iter
}

// Account returns the account which the iterator is currently pointing to.
func (i *Iter) Account() *stripe.Account {
	return i.Current().(*stripe.Account)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
