// Package loginlink provides the /login_links APIs
package loginlink

import (
	"errors"
	"net/http"

	stripe "github.com/stripe/stripe-go"
)

// Client is used to invoke /login_links APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a login link for an account.
func New(params *stripe.LoginLinkParams) (*stripe.LoginLink, error) {
	return getC().New(params)
}

// New creates a login link for an account.
func (c Client) New(params *stripe.LoginLinkParams) (*stripe.LoginLink, error) {
	if params.Account == nil {
		return nil, errors.New("Invalid login link params: Account must be set")
	}

	path := stripe.FormatURLPath("/v1/accounts/%s/login_links", stripe.StringValue(params.Account))
	loginLink := &stripe.LoginLink{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, loginLink)
	return loginLink, err
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
