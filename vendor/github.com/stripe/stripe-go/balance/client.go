// Package balance provides the /balance APIs
package balance

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /balance and transaction-related APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get retrieves an account's balance
func Get(params *stripe.BalanceParams) (*stripe.Balance, error) {
	return getC().Get(params)
}

// Get retrieves an account's balance
func (c Client) Get(params *stripe.BalanceParams) (*stripe.Balance, error) {
	balance := &stripe.Balance{}
	err := c.B.Call(http.MethodGet, "/v1/balance", c.Key, params, balance)
	return balance, err
}

// GetBalanceTransaction retrieves a balance transaction
func GetBalanceTransaction(id string, params *stripe.BalanceTransactionParams) (*stripe.BalanceTransaction, error) {
	return getC().GetBalanceTransaction(id, params)
}

// GetBalanceTransaction retrieves a balance transaction
func (c Client) GetBalanceTransaction(id string, params *stripe.BalanceTransactionParams) (*stripe.BalanceTransaction, error) {
	path := stripe.FormatURLPath("/v1/balance/history/%s", id)
	balance := &stripe.BalanceTransaction{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, balance)
	return balance, err
}

// List returns a list of balance transactions.
func List(params *stripe.BalanceTransactionListParams) *Iter {
	return getC().List(params)
}

// List returns a list of balance transactions.
func (c Client) List(listParams *stripe.BalanceTransactionListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.BalanceTransactionList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/balance/history", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for balance transactions.
type Iter struct {
	*stripe.Iter
}

// BalanceTransaction returns the balance transaction which the iterator is currently pointing to.
func (i *Iter) BalanceTransaction() *stripe.BalanceTransaction {
	return i.Current().(*stripe.BalanceTransaction)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
