package stripe

import (
	"encoding/json"
	"strconv"

	"github.com/stripe/stripe-go/form"
)

// CardAvailablePayoutMethod is a set of available payout methods for the card.
type CardAvailablePayoutMethod string

// List of values that CardAvailablePayoutMethod can take.
const (
	CardAvailablePayoutMethodInstant  CardAvailablePayoutMethod = "Instant"
	CardAvailablePayoutMethodStandard CardAvailablePayoutMethod = "Standard"
)

// CardBrand is the list of allowed values for the card's brand.
type CardBrand string

// List of values that CardBrand can take.
const (
	CardBrandAmex       CardBrand = "American Express"
	CardBrandDiscover   CardBrand = "Discover"
	CardBrandDinersClub CardBrand = "Diners Club"
	CardBrandJCB        CardBrand = "JCB"
	CardBrandMasterCard CardBrand = "MasterCard"
	CardBrandUnknown    CardBrand = "Unknown"
	CardBrandUnionPay   CardBrand = "UnionPay"
	CardBrandVisa       CardBrand = "Visa"
)

// CardFunding is the list of allowed values for the card's funding.
type CardFunding string

// List of values that CardFunding can take.
const (
	CardFundingCredit  CardFunding = "credit"
	CardFundingDebit   CardFunding = "debit"
	CardFundingPrepaid CardFunding = "prepaid"
	CardFundingUnknown CardFunding = "unknown"
)

// CardTokenizationMethod is the list of allowed values for the card's tokenization method.
type CardTokenizationMethod string

// List of values that CardTokenizationMethod can take.
const (
	TokenizationMethodAndroidPay CardTokenizationMethod = "android_pay"
	TokenizationMethodApplePay   CardTokenizationMethod = "apple_pay"
)

// CardVerification is the list of allowed verification responses.
type CardVerification string

// List of values that CardVerification can take.
const (
	CardVerificationFail      CardVerification = "fail"
	CardVerificationPass      CardVerification = "pass"
	CardVerificationUnchecked CardVerification = "unchecked"
)

// cardSource is a string that's used to build card form parameters. It's a
// constant just to make mistakes less likely.
const cardSource = "source"

// CardParams is the set of parameters that can be used when creating or updating a card.
// For more details see https://stripe.com/docs/api#create_card and https://stripe.com/docs/api#update_card.
//
// Note that while form annotations are used for tokenization and updates,
// cards have some unusual logic on creates that necessitates manual handling
// of all parameters. See AppendToAsCardSourceOrExternalAccount.
type CardParams struct {
	Params             `form:"*"`
	Account            *string `form:"-"`
	AddressCity        *string `form:"address_city"`
	AddressCountry     *string `form:"address_country"`
	AddressLine1       *string `form:"address_line1"`
	AddressLine2       *string `form:"address_line2"`
	AddressState       *string `form:"address_state"`
	AddressZip         *string `form:"address_zip"`
	CVC                *string `form:"cvc"`
	Currency           *string `form:"currency"`
	Customer           *string `form:"-"`
	DefaultForCurrency *bool   `form:"default_for_currency"`
	ExpMonth           *string `form:"exp_month"`
	ExpYear            *string `form:"exp_year"`
	Name               *string `form:"name"`
	Number             *string `form:"number"`
	Recipient          *string `form:"-"`
	Token              *string `form:"-"`

	// ID is used when tokenizing a card for shared customers
	ID string `form:"*"`
}

// AppendToAsCardSourceOrExternalAccount appends the given CardParams as either a
// card or external account.
//
// It may look like an AppendTo from the form package, but it's not, and is
// only used in the special case where we use `card.New`. It's needed because
// we have some weird encoding logic here that can't be handled by the form
// package (and it's special enough that it wouldn't be desirable to have it do
// so).
//
// This is not a pattern that we want to push forward, and this largely exists
// because the cards endpoint is a little unusual. There is one other resource
// like it, which is bank account.
func (c *CardParams) AppendToAsCardSourceOrExternalAccount(body *form.Values, keyParts []string) {
	// Rather than being called in addition to `AppendTo`, this function
	// *replaces* `AppendTo`, so we must also make sure to handle the encoding
	// of `Params` so metadata and the like is included in the encoded payload.
	form.AppendToPrefixed(body, c.Params, keyParts)

	if c.DefaultForCurrency != nil {
		body.Add(form.FormatKey(append(keyParts, "default_for_currency")), strconv.FormatBool(BoolValue(c.DefaultForCurrency)))
	}

	if c.Token != nil {
		if c.Account != nil {
			body.Add(form.FormatKey(append(keyParts, "external_account")), StringValue(c.Token))
		} else {
			body.Add(form.FormatKey(append(keyParts, cardSource)), StringValue(c.Token))
		}
	}

	if c.Number != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "object")), "card")
		body.Add(form.FormatKey(append(keyParts, cardSource, "number")), StringValue(c.Number))
	}

	if c.CVC != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "cvc")), StringValue(c.CVC))
	}

	if c.Currency != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "currency")), StringValue(c.Currency))
	}

	if c.ExpMonth != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "exp_month")), StringValue(c.ExpMonth))
	}

	if c.ExpYear != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "exp_year")), StringValue(c.ExpYear))
	}

	if c.Name != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "name")), StringValue(c.Name))
	}

	if c.AddressCity != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "address_city")), StringValue(c.AddressCity))
	}

	if c.AddressCountry != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "address_country")), StringValue(c.AddressCountry))
	}

	if c.AddressLine1 != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "address_line1")), StringValue(c.AddressLine1))
	}

	if c.AddressLine2 != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "address_line2")), StringValue(c.AddressLine2))
	}

	if c.AddressState != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "address_state")), StringValue(c.AddressState))
	}

	if c.AddressZip != nil {
		body.Add(form.FormatKey(append(keyParts, cardSource, "address_zip")), StringValue(c.AddressZip))
	}
}

// CardListParams is the set of parameters that can be used when listing cards.
// For more details see https://stripe.com/docs/api#list_cards.
type CardListParams struct {
	ListParams `form:"*"`
	Account    *string `form:"-"`
	Customer   *string `form:"-"`
	Recipient  *string `form:"-"`
}

// AppendTo implements custom encoding logic for CardListParams
// so that we can send the special required `object` field up along with the
// other specified parameters.
func (p *CardListParams) AppendTo(body *form.Values, keyParts []string) {
	if p.Account != nil || p.Customer != nil {
		body.Add(form.FormatKey(append(keyParts, "object")), "card")
	}
}

// Card is the resource representing a Stripe credit/debit card.
// For more details see https://stripe.com/docs/api#cards.
type Card struct {
	AddressCity            string                      `json:"address_city"`
	AddressCountry         string                      `json:"address_country"`
	AddressLine1           string                      `json:"address_line1"`
	AddressLine1Check      CardVerification            `json:"address_line1_check"`
	AddressLine2           string                      `json:"address_line2"`
	AddressState           string                      `json:"address_state"`
	AddressZip             string                      `json:"address_zip"`
	AddressZipCheck        CardVerification            `json:"address_zip_check"`
	AvailablePayoutMethods []CardAvailablePayoutMethod `json:"available_payout_methods"`
	Brand                  CardBrand                   `json:"brand"`
	CVCCheck               CardVerification            `json:"cvc_check"`
	Country                string                      `json:"country"`
	Currency               Currency                    `json:"currency"`
	Customer               *Customer                   `json:"customer"`
	DefaultForCurrency     bool                        `json:"default_for_currency"`
	Deleted                bool                        `json:"deleted"`

	// Description is a succinct summary of the card's information.
	//
	// Please note that this field is for internal use only and is not returned
	// as part of standard API requests.
	Description string `json:"description"`

	DynamicLast4 string      `json:"dynamic_last4"`
	ExpMonth     uint8       `json:"exp_month"`
	ExpYear      uint16      `json:"exp_year"`
	Fingerprint  string      `json:"fingerprint"`
	Funding      CardFunding `json:"funding"`
	ID           string      `json:"id"`

	// IIN is the card's "Issuer Identification Number".
	//
	// Please note that this field is for internal use only and is not returned
	// as part of standard API requests.
	IIN string `json:"iin"`

	// Issuer is a bank or financial institution that provides the card.
	//
	// Please note that this field is for internal use only and is not returned
	// as part of standard API requests.
	Issuer string `json:"issuer"`

	Last4              string                 `json:"last4"`
	Metadata           map[string]string      `json:"metadata"`
	Name               string                 `json:"name"`
	Recipient          *Recipient             `json:"recipient"`
	ThreeDSecure       *ThreeDSecure          `json:"three_d_secure"`
	TokenizationMethod CardTokenizationMethod `json:"tokenization_method"`
}

// CardList is a list object for cards.
type CardList struct {
	ListMeta
	Data []*Card `json:"data"`
}

// UnmarshalJSON handles deserialization of a Card.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (c *Card) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		c.ID = id
		return nil
	}

	type card Card
	var v card
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = Card(v)
	return nil
}
