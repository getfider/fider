// Package card provides API functions related to issuing cards.
//
// For more details, see: https://stripe.com/docs/api/go#issuing_cards
package card

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /issuing/cards APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new issuing card.
func New(params *stripe.IssuingCardParams) (*stripe.IssuingCard, error) {
	return getC().New(params)
}

// New creates a new issuing card.
func (c Client) New(params *stripe.IssuingCardParams) (*stripe.IssuingCard, error) {
	card := &stripe.IssuingCard{}
	err := c.B.Call(http.MethodPost, "/v1/issuing/cards", c.Key, params, card)
	return card, err
}

// Get returns the details of an issuing card.
func Get(id string, params *stripe.IssuingCardParams) (*stripe.IssuingCard, error) {
	return getC().Get(id, params)
}

// Get returns the details of an issuing card.
func (c Client) Get(id string, params *stripe.IssuingCardParams) (*stripe.IssuingCard, error) {
	path := stripe.FormatURLPath("/v1/issuing/cards/%s", id)
	card := &stripe.IssuingCard{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, card)
	return card, err
}

// Update updates an issuing card.
func Update(id string, params *stripe.IssuingCardParams) (*stripe.IssuingCard, error) {
	return getC().Update(id, params)
}

// Update updates an issuing card.
func (c Client) Update(id string, params *stripe.IssuingCardParams) (*stripe.IssuingCard, error) {
	path := stripe.FormatURLPath("/v1/issuing/cards/%s", id)
	card := &stripe.IssuingCard{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, card)
	return card, err
}

// List returns a list of issuing cards.
func List(params *stripe.IssuingCardListParams) *Iter {
	return getC().List(params)
}

// List returns a list of issuing cards.
func (c Client) List(listParams *stripe.IssuingCardListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.IssuingCardList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/issuing/cards", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Details retrieves an issuing card details.
func Details(id string, params *stripe.IssuingCardParams) (*stripe.IssuingCardDetails, error) {
	return getC().Details(id, params)
}

// Details retrieves an issuing card details.
func (c Client) Details(id string, params *stripe.IssuingCardParams) (*stripe.IssuingCardDetails, error) {
	path := stripe.FormatURLPath("/v1/issuing/cards/%s/details", id)
	cardDetails := &stripe.IssuingCardDetails{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, cardDetails)
	return cardDetails, err
}

// Iter is an iterator for issuing cards.
type Iter struct {
	*stripe.Iter
}

// IssuingCard returns the issuing card which the iterator is currently pointing to.
func (i *Iter) IssuingCard() *stripe.IssuingCard {
	return i.Current().(*stripe.IssuingCard)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
