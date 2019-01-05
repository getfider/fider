// Package card provides the /cards APIs
package card

import (
	"errors"
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /cards APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new card.
func New(params *stripe.CardParams) (*stripe.Card, error) {
	return getC().New(params)
}

// New creates a new card.
func (c Client) New(params *stripe.CardParams) (*stripe.Card, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	var path string
	if params.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts",
			stripe.StringValue(params.Account))
	} else if params.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources",
			stripe.StringValue(params.Customer))
	} else if params.Recipient != nil {
		path = stripe.FormatURLPath("/v1/recipients/%s/cards",
			stripe.StringValue(params.Recipient))
	} else {
		return nil, errors.New("Invalid card params: either account, customer or recipient need to be set")
	}

	body := &form.Values{}

	// Note that we call this special append method instead of the standard one
	// from the form package. We should not use form's because doing so will
	// include some parameters that are undesirable here.
	params.AppendToAsCardSourceOrExternalAccount(body, nil)

	// Because card creation uses the custom append above, we have to make an
	// explicit call using a form and CallRaw instead of the standard Call
	// (which takes a set of parameters).
	card := &stripe.Card{}
	err := c.B.CallRaw(http.MethodPost, path, c.Key, body, &params.Params, card)
	return card, err
}

// Get returns the details of a card.
func Get(id string, params *stripe.CardParams) (*stripe.Card, error) {
	return getC().Get(id, params)
}

// Get returns the details of a card.
func (c Client) Get(id string, params *stripe.CardParams) (*stripe.Card, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	var path string
	if params.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts/%s",
			stripe.StringValue(params.Account), id)
	} else if params.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources/%s",
			stripe.StringValue(params.Customer), id)
	} else if params.Recipient != nil {
		path = stripe.FormatURLPath("/v1/recipients/%s/cards/%s",
			stripe.StringValue(params.Recipient), id)
	} else {
		return nil, errors.New("Invalid card params: either account, customer or recipient need to be set")
	}

	card := &stripe.Card{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, card)
	return card, err
}

// Update updates a card's properties.
func Update(id string, params *stripe.CardParams) (*stripe.Card, error) {
	return getC().Update(id, params)
}

// Update updates a card's properties.
func (c Client) Update(id string, params *stripe.CardParams) (*stripe.Card, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	var path string
	if params.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts/%s",
			stripe.StringValue(params.Account), id)
	} else if params.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources/%s",
			stripe.StringValue(params.Customer), id)
	} else if params.Recipient != nil {
		path = stripe.FormatURLPath("/v1/recipients/%s/cards/%s",
			stripe.StringValue(params.Recipient), id)
	} else {
		return nil, errors.New("Invalid card params: either account, customer or recipient need to be set")
	}

	card := &stripe.Card{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, card)
	return card, err
}

// Del removes a card.
func Del(id string, params *stripe.CardParams) (*stripe.Card, error) {
	return getC().Del(id, params)
}

// Del removes a card.
func (c Client) Del(id string, params *stripe.CardParams) (*stripe.Card, error) {
	if params == nil {
		return nil, errors.New("params should not be nil")
	}

	var path string
	if params.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts/%s", stripe.StringValue(params.Account), id)
	} else if params.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources/%s", stripe.StringValue(params.Customer), id)
	} else if params.Recipient != nil {
		path = stripe.FormatURLPath("/v1/recipients/%s/cards/%s", stripe.StringValue(params.Recipient), id)
	} else {
		return nil, errors.New("Invalid card params: either account, customer or recipient need to be set")
	}

	card := &stripe.Card{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, card)
	return card, err
}

// List returns a list of cards.
func List(params *stripe.CardListParams) *Iter {
	return getC().List(params)
}

// List returns a list of cards.
func (c Client) List(listParams *stripe.CardListParams) *Iter {
	var path string
	var outerErr error

	// Because the external accounts and sources endpoints can return non-card
	// objects, CardListParam's `AppendTo` will add the filter `object=card` to
	// make sure that only cards come back with the response.
	if listParams == nil {
		outerErr = errors.New("params should not be nil")
	} else if listParams.Account != nil {
		path = stripe.FormatURLPath("/v1/accounts/%s/external_accounts",
			stripe.StringValue(listParams.Account))
	} else if listParams.Customer != nil {
		path = stripe.FormatURLPath("/v1/customers/%s/sources",
			stripe.StringValue(listParams.Customer))
	} else if listParams.Recipient != nil {
		path = stripe.FormatURLPath("/v1/recipients/%s/cards", stripe.StringValue(listParams.Recipient))
	} else {
		outerErr = errors.New("Invalid card params: either account, customer or recipient need to be set")
	}

	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.CardList{}

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

// Iter is an iterator for cards.
type Iter struct {
	*stripe.Iter
}

// Card returns the card which the iterator is currently pointing to.
func (i *Iter) Card() *stripe.Card {
	return i.Current().(*stripe.Card)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
