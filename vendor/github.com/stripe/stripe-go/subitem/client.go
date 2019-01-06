// Package subitem provides the /subscription_items APIs
package subitem

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /subscriptions APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new subscription item.
func New(params *stripe.SubscriptionItemParams) (*stripe.SubscriptionItem, error) {
	return getC().New(params)
}

// New creates a new subscription item.
func (c Client) New(params *stripe.SubscriptionItemParams) (*stripe.SubscriptionItem, error) {
	item := &stripe.SubscriptionItem{}
	err := c.B.Call(http.MethodPost, "/v1/subscription_items", c.Key, params, item)
	return item, err
}

// Get returns the details of a subscription item.
func Get(id string, params *stripe.SubscriptionItemParams) (*stripe.SubscriptionItem, error) {
	return getC().Get(id, params)
}

// Get returns the details of a subscription item.
func (c Client) Get(id string, params *stripe.SubscriptionItemParams) (*stripe.SubscriptionItem, error) {
	path := stripe.FormatURLPath("/v1/subscription_items/%s", id)
	item := &stripe.SubscriptionItem{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, item)
	return item, err
}

// Update updates a subscription item's properties.
func Update(id string, params *stripe.SubscriptionItemParams) (*stripe.SubscriptionItem, error) {
	return getC().Update(id, params)
}

// Update updates a subscription item's properties.
func (c Client) Update(id string, params *stripe.SubscriptionItemParams) (*stripe.SubscriptionItem, error) {
	path := stripe.FormatURLPath("/v1/subscription_items/%s", id)
	subi := &stripe.SubscriptionItem{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, subi)
	return subi, err
}

// Del removes a subscription item.
func Del(id string, params *stripe.SubscriptionItemParams) (*stripe.SubscriptionItem, error) {
	return getC().Del(id, params)
}

// Del removes a subscription item.
func (c Client) Del(id string, params *stripe.SubscriptionItemParams) (*stripe.SubscriptionItem, error) {
	path := stripe.FormatURLPath("/v1/subscription_items/%s", id)
	item := &stripe.SubscriptionItem{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, item)

	return item, err
}

// List returns a list of subscription items.
func List(params *stripe.SubscriptionItemListParams) *Iter {
	return getC().List(params)
}

// List returns a list of subscription items.
func (c Client) List(listParams *stripe.SubscriptionItemListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.SubscriptionItemList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/subscription_items", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for subscription items.
type Iter struct {
	*stripe.Iter
}

// SubscriptionItem returns the subscription item which the iterator is currently pointing to.
func (i *Iter) SubscriptionItem() *stripe.SubscriptionItem {
	return i.Current().(*stripe.SubscriptionItem)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
