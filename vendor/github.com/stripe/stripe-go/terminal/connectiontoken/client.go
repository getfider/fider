// Package connectiontoken provides API functions related to terminal connection tokens
package connectiontoken

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
)

// Client is used to invoke /terminal/connection_tokens APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new terminal connection token.
func New(params *stripe.TerminalConnectionTokenParams) (*stripe.TerminalConnectionToken, error) {
	return getC().New(params)
}

// New creates a new terminal connection token.
func (c Client) New(params *stripe.TerminalConnectionTokenParams) (*stripe.TerminalConnectionToken, error) {
	connectiontoken := &stripe.TerminalConnectionToken{}
	err := c.B.Call(http.MethodPost, "/v1/terminal/connection_tokens", c.Key, params, connectiontoken)
	return connectiontoken, err
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
