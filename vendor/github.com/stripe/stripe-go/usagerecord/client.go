// Package usagerecord provides the /subscription_items/{SUBSCRIPTION_ITEM_ID}/usage_records APIs
package usagerecord

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
)

// Client is used to invoke APIs related to usage records.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new usage record.
func New(params *stripe.UsageRecordParams) (*stripe.UsageRecord, error) {
	return getC().New(params)
}

// New creates a new usage record.
func (c Client) New(params *stripe.UsageRecordParams) (*stripe.UsageRecord, error) {
	path := stripe.FormatURLPath("/v1/subscription_items/%s/usage_records", stripe.StringValue(params.SubscriptionItem))
	record := &stripe.UsageRecord{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, record)
	return record, err
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
