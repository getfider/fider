package stripe

import (
	"encoding/json"
	"strconv"

	"github.com/stripe/stripe-go/form"
)

// BankAccountStatus is the list of allowed values for the bank account's status.
type BankAccountStatus string

// List of values that BankAccountStatus can take.
const (
	BankAccountStatusErrored            BankAccountStatus = "errored"
	BankAccountStatusNew                BankAccountStatus = "new"
	BankAccountStatusValidated          BankAccountStatus = "validated"
	BankAccountStatusVerificationFailed BankAccountStatus = "verification_failed"
	BankAccountStatusVerified           BankAccountStatus = "verified"
)

// BankAccountAccountHolderType is the list of allowed values for the bank account holder type.
type BankAccountAccountHolderType string

// List of values that BankAccountAccountHolderType can take.
const (
	BankAccountAccountHolderTypeCompany    BankAccountAccountHolderType = "company"
	BankAccountAccountHolderTypeIndividual BankAccountAccountHolderType = "individual"
)

// BankAccountParams is the set of parameters that can be used when updating a
// bank account.
//
// Note that while form annotations are used for updates, bank accounts have
// some unusual logic on creates that necessitates manual handling of all
// parameters. See AppendToAsSourceOrExternalAccount.
type BankAccountParams struct {
	Params `form:"*"`

	// Account is the identifier of the parent account under which bank
	// accounts are nested.
	Account *string `form:"-"`

	AccountHolderName  *string `form:"account_holder_name"`
	AccountHolderType  *string `form:"account_holder_type"`
	AccountNumber      *string `form:"account_number"`
	Country            *string `form:"country"`
	Currency           *string `form:"currency"`
	Customer           *string `form:"-"`
	DefaultForCurrency *bool   `form:"default_for_currency"`
	RoutingNumber      *string `form:"routing_number"`

	// Token is a token referencing an external account like one returned from
	// Stripe.js.
	Token *string `form:"-"`

	// ID is used when tokenizing a bank account for shared customers
	ID *string `form:"*"`
}

// AppendToAsSourceOrExternalAccount appends the given BankAccountParams as
// either a source or external account.
//
// It may look like an AppendTo from the form package, but it's not, and is
// only used in the special case where we use `bankaccount.New`. It's needed
// because we have some weird encoding logic here that can't be handled by the
// form package (and it's special enough that it wouldn't be desirable to have
// it do so).
//
// This is not a pattern that we want to push forward, and this largely exists
// because the bank accounts endpoint is a little unusual. There is one other
// resource like it, which is cards.
func (a *BankAccountParams) AppendToAsSourceOrExternalAccount(body *form.Values) {
	// Rather than being called in addition to `AppendTo`, this function
	// *replaces* `AppendTo`, so we must also make sure to handle the encoding
	// of `Params` so metadata and the like is included in the encoded payload.
	form.AppendTo(body, a.Params)

	isCustomer := a.Customer != nil

	var sourceType string
	if isCustomer {
		sourceType = "source"
	} else {
		sourceType = "external_account"
	}

	// Use token (if exists) or a dictionary containing a user’s bank account details.
	if a.Token != nil {
		body.Add(sourceType, StringValue(a.Token))

		if a.DefaultForCurrency != nil {
			body.Add("default_for_currency", strconv.FormatBool(BoolValue(a.DefaultForCurrency)))
		}
	} else {
		body.Add(sourceType+"[object]", "bank_account")
		body.Add(sourceType+"[country]", StringValue(a.Country))
		body.Add(sourceType+"[account_number]", StringValue(a.AccountNumber))
		body.Add(sourceType+"[currency]", StringValue(a.Currency))

		// These are optional and the API will fail if we try to send empty
		// values in for them, so make sure to check that they're actually set
		// before encoding them.
		if a.AccountHolderName != nil {
			body.Add(sourceType+"[account_holder_name]", StringValue(a.AccountHolderName))
		}

		if a.AccountHolderType != nil {
			body.Add(sourceType+"[account_holder_type]", StringValue(a.AccountHolderType))
		}

		if a.RoutingNumber != nil {
			body.Add(sourceType+"[routing_number]", StringValue(a.RoutingNumber))
		}

		if a.DefaultForCurrency != nil {
			body.Add(sourceType+"[default_for_currency]", strconv.FormatBool(BoolValue(a.DefaultForCurrency)))
		}
	}
}

// BankAccountListParams is the set of parameters that can be used when listing bank accounts.
type BankAccountListParams struct {
	ListParams `form:"*"`

	// The identifier of the parent account under which the bank accounts are
	// nested. Either Account or Customer should be populated.
	Account *string `form:"-"`

	// The identifier of the parent customer under which the bank accounts are
	// nested. Either Account or Customer should be populated.
	Customer *string `form:"-"`
}

// AppendTo implements custom encoding logic for BankAccountListParams
// so that we can send the special required `object` field up along with the
// other specified parameters.
func (p *BankAccountListParams) AppendTo(body *form.Values, keyParts []string) {
	body.Add(form.FormatKey(append(keyParts, "object")), "bank_account")
}

// BankAccount represents a Stripe bank account.
type BankAccount struct {
	AccountHolderName  string                       `json:"account_holder_name"`
	AccountHolderType  BankAccountAccountHolderType `json:"account_holder_type"`
	BankName           string                       `json:"bank_name"`
	Country            string                       `json:"country"`
	Currency           Currency                     `json:"currency"`
	Customer           *Customer                    `json:"customer"`
	DefaultForCurrency bool                         `json:"default_for_currency"`
	Deleted            bool                         `json:"deleted"`
	Fingerprint        string                       `json:"fingerprint"`
	ID                 string                       `json:"id"`
	Last4              string                       `json:"last4"`
	Metadata           map[string]string            `json:"metadata"`
	RoutingNumber      string                       `json:"routing_number"`
	Status             BankAccountStatus            `json:"status"`
}

// BankAccountList is a list object for bank accounts.
type BankAccountList struct {
	ListMeta
	Data []*BankAccount `json:"data"`
}

// UnmarshalJSON handles deserialization of a BankAccount.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (b *BankAccount) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		b.ID = id
		return nil
	}

	type bankAccount BankAccount
	var v bankAccount
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BankAccount(v)
	return nil
}
