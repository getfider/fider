// Package discount provides the discount-related APIs
package discount

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
)

// Client is used to invoke discount-related APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Del removes a discount from a customer.
func Del(customerID string, params *stripe.DiscountParams) (*stripe.Discount, error) {
	return getC().Del(customerID, params)
}

// Del removes a discount from a customer.
func (c Client) Del(customerID string, params *stripe.DiscountParams) (*stripe.Discount, error) {
	path := stripe.FormatURLPath("/v1/customers/%s/discount", customerID)
	discount := &stripe.Discount{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, discount)
	return discount, err
}

// DelSubscription removes a discount from a customer's subscription.
func DelSubscription(subscriptionID string, params *stripe.DiscountParams) (*stripe.Discount, error) {
	return getC().DelSub(subscriptionID, params)
}

// DelSub removes a discount from a customer's subscription.
func (c Client) DelSub(subscriptionID string, params *stripe.DiscountParams) (*stripe.Discount, error) {
	path := stripe.FormatURLPath("/v1/subscriptions/%s/discount", subscriptionID)
	discount := &stripe.Discount{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, discount)

	return discount, err
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
