// Package token provides the /tokens APIs
package token

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
)

// Client is used to invoke /tokens APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new token.
func New(params *stripe.TokenParams) (*stripe.Token, error) {
	return getC().New(params)
}

// New creates a new token.
func (c Client) New(params *stripe.TokenParams) (*stripe.Token, error) {
	tok := &stripe.Token{}
	err := c.B.Call(http.MethodPost, "/v1/tokens", c.Key, params, tok)
	return tok, err
}

// Get returns the details of a token.
func Get(id string, params *stripe.TokenParams) (*stripe.Token, error) {
	return getC().Get(id, params)
}

// Get returns the details of a token.
func (c Client) Get(id string, params *stripe.TokenParams) (*stripe.Token, error) {
	path := stripe.FormatURLPath("/v1/tokens/%s", id)
	token := &stripe.Token{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, token)

	return token, err
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
