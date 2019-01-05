package stripe

import (
	"encoding/json"

	"github.com/stripe/stripe-go/form"
)

// AccountType is the type of an account.
type AccountType string

// List of values that AccountType can take.
const (
	AccountTypeCustom   AccountType = "custom"
	AccountTypeExpress  AccountType = "express"
	AccountTypeStandard AccountType = "standard"
)

// ExternalAccountType is the type of an external account.
type ExternalAccountType string

// List of values that ExternalAccountType can take.
const (
	ExternalAccountTypeBankAccount ExternalAccountType = "bank_account"
	ExternalAccountTypeCard        ExternalAccountType = "card"
)

// LegalEntityType describes the types for a legal entity.
type LegalEntityType string

// List of values that LegalEntityType can take.
const (
	LegalEntityTypeCompany    LegalEntityType = "company"
	LegalEntityTypeIndividual LegalEntityType = "individual"
)

// IdentityVerificationDetailsCode is a machine-readable code specifying the
// verification state of a legal entity
type IdentityVerificationDetailsCode string

// List of values that IdentityVerificationDetailsCode can take.
const (
	IdentityVerificationDetailsCodeFailedKeyedIdentity       IdentityVerificationDetailsCode = "failed_keyed_identity"
	IdentityVerificationDetailsCodeFailedOther               IdentityVerificationDetailsCode = "failed_other"
	IdentityVerificationDetailsCodeScanCorrupt               IdentityVerificationDetailsCode = "scan_corrupt"
	IdentityVerificationDetailsCodeScanFailedGreyscale       IdentityVerificationDetailsCode = "scan_failed_greyscale"
	IdentityVerificationDetailsCodeScanFailedOther           IdentityVerificationDetailsCode = "scan_failed_other"
	IdentityVerificationDetailsCodeScanIDCountryNotSupported IdentityVerificationDetailsCode = "scan_id_country_not_supported"
	IdentityVerificationDetailsCodeScanIDTypeNotSupported    IdentityVerificationDetailsCode = "scan_id_type_not_supported"
	IdentityVerificationDetailsCodeScanNameMismatch          IdentityVerificationDetailsCode = "scan_name_mismatch"
	IdentityVerificationDetailsCodeScanNotReadable           IdentityVerificationDetailsCode = "scan_not_readable"
	IdentityVerificationDetailsCodeScanNotUploaded           IdentityVerificationDetailsCode = "scan_not_uploaded"
)

// IdentityVerificationDisabledReason describes the valid reason to disable account
type IdentityVerificationDisabledReason string

// List of values that IdentityVerificationDisabledReason can take.
const (
	IdentityVerificationDisabledReasonFieldsNeeded           IdentityVerificationDisabledReason = "fields_needed"
	IdentityVerificationDisabledReasonListed                 IdentityVerificationDisabledReason = "listed"
	IdentityVerificationDisabledReasonOther                  IdentityVerificationDisabledReason = "other"
	IdentityVerificationDisabledReasonRejectedFraud          IdentityVerificationDisabledReason = "rejected.fraud"
	IdentityVerificationDisabledReasonRejectedListed         IdentityVerificationDisabledReason = "rejected.listed"
	IdentityVerificationDisabledReasonRejectedOther          IdentityVerificationDisabledReason = "rejected.other"
	IdentityVerificationDisabledReasonRejectedTermsOfService IdentityVerificationDisabledReason = "rejected.terms_of_service"
	IdentityVerificationDisabledReasonUnderReview            IdentityVerificationDisabledReason = "under_review"
)

// IdentityVerificationStatus describes the different statuses for identity verification.
type IdentityVerificationStatus string

// List of values that IdentityVerificationStatus can take.
const (
	IdentityVerificationStatusPending    IdentityVerificationStatus = "pending"
	IdentityVerificationStatusUnverified IdentityVerificationStatus = "unverified"
	IdentityVerificationStatusVerified   IdentityVerificationStatus = "verified"
)

// PayoutInterval describes the payout interval.
type PayoutInterval string

// List of values that PayoutInterval can take.
const (
	PayoutIntervalDaily   PayoutInterval = "daily"
	PayoutIntervalManual  PayoutInterval = "manual"
	PayoutIntervalMonthly PayoutInterval = "monthly"
	PayoutIntervalWeekly  PayoutInterval = "weekly"
)

// AccountRejectReason describes the valid reason to reject an account
type AccountRejectReason string

// List of values that AccountRejectReason can take.
const (
	AccountRejectReasonFraud          AccountRejectReason = "fraud"
	AccountRejectReasonOther          AccountRejectReason = "other"
	AccountRejectReasonTermsOfService AccountRejectReason = "terms_of_service"
)

// AccountParams are the parameters allowed during account creation/updates.
type AccountParams struct {
	Params                    `form:"*"`
	AccountToken              *string                       `form:"account_token"`
	BusinessName              *string                       `form:"business_name"`
	BusinessPrimaryColor      *string                       `form:"business_primary_color"`
	BusinessURL               *string                       `form:"business_url"`
	Country                   *string                       `form:"country"`
	DeclineChargeOn           *AccountDeclineSettingsParams `form:"decline_charge_on"`
	DebitNegativeBalances     *bool                         `form:"debit_negative_balances"`
	DefaultCurrency           *string                       `form:"default_currency"`
	Email                     *string                       `form:"email"`
	ExternalAccount           *AccountExternalAccountParams `form:"external_account"`
	FromRecipient             *string                       `form:"from_recipient"`
	LegalEntity               *LegalEntityParams            `form:"legal_entity"`
	PayoutSchedule            *PayoutScheduleParams         `form:"payout_schedule"`
	PayoutStatementDescriptor *string                       `form:"payout_statement_descriptor"`
	ProductDescription        *string                       `form:"product_description"`
	StatementDescriptor       *string                       `form:"statement_descriptor"`
	SupportEmail              *string                       `form:"support_email"`
	SupportPhone              *string                       `form:"support_phone"`
	SupportURL                *string                       `form:"support_url"`
	TOSAcceptance             *TOSAcceptanceParams          `form:"tos_acceptance"`
	Type                      *string                       `form:"type"`
}

// LegalEntityParams represents a legal_entity during account creation/updates.
type LegalEntityParams struct {
	AdditionalOwners []*AdditionalOwnerParams `form:"additional_owners"`

	// AdditionalOwnersEmpty can be set to clear a legal entity's additional
	// owners.
	AdditionalOwnersEmpty bool `form:"additional_owners,empty"`

	Address              *AccountAddressParams       `form:"address"`
	AddressKana          *AccountAddressParams       `form:"address_kana"`
	AddressKanji         *AccountAddressParams       `form:"address_kanji"`
	BusinessName         *string                     `form:"business_name"`
	BusinessNameKana     *string                     `form:"business_name_kana"`
	BusinessNameKanji    *string                     `form:"business_name_kanji"`
	BusinessTaxID        *string                     `form:"business_tax_id"`
	BusinessVATID        *string                     `form:"business_vat_id"`
	DOB                  *DOBParams                  `form:"dob"`
	FirstName            *string                     `form:"first_name"`
	FirstNameKana        *string                     `form:"first_name_kana"`
	FirstNameKanji       *string                     `form:"first_name_kanji"`
	Gender               *string                     `form:"gender"`
	LastName             *string                     `form:"last_name"`
	LastNameKana         *string                     `form:"last_name_kana"`
	LastNameKanji        *string                     `form:"last_name_kanji"`
	MaidenName           *string                     `form:"maiden_name"`
	PersonalAddress      *AccountAddressParams       `form:"personal_address"`
	PersonalAddressKana  *AccountAddressParams       `form:"personal_address_kana"`
	PersonalAddressKanji *AccountAddressParams       `form:"personal_address_kanji"`
	PersonalIDNumber     *string                     `form:"personal_id_number"`
	PhoneNumber          *string                     `form:"phone_number"`
	SSNLast4             *string                     `form:"ssn_last_4"`
	Type                 *string                     `form:"type"`
	Verification         *IdentityVerificationParams `form:"verification"`
}

// AccountAddressParams represents an address during account creation/updates.
type AccountAddressParams struct {
	City       *string `form:"city"`
	Country    *string `form:"country"`
	Line1      *string `form:"line1"`
	Line2      *string `form:"line2"`
	PostalCode *string `form:"postal_code"`
	State      *string `form:"state"`

	// Town/cho-me. Note that this is only used for Kana/Kanji representations
	// of an address.
	Town *string `form:"town"`
}

// AccountDeclineSettingsParams represents the parameters allowed for configuring
// declines on connected accounts.
type AccountDeclineSettingsParams struct {
	AVSFailure *bool `form:"avs_failure"`
	CVCFailure *bool `form:"cvc_failure"`
}

// DOBParams represents a DOB during account creation/updates.
type DOBParams struct {
	Day   *int64 `form:"day"`
	Month *int64 `form:"month"`
	Year  *int64 `form:"year"`
}

// TOSAcceptanceParams represents tos_acceptance during account creation/updates.
type TOSAcceptanceParams struct {
	Date      *int64  `form:"date"`
	IP        *string `form:"ip"`
	UserAgent *string `form:"user_agent"`
}

// AdditionalOwnerParams represents an additional owner during account creation/updates.
type AdditionalOwnerParams struct {
	Address          *AccountAddressParams       `form:"address"`
	DOB              *DOBParams                  `form:"dob"`
	FirstName        *string                     `form:"first_name"`
	LastName         *string                     `form:"last_name"`
	MaidenName       *string                     `form:"maiden_name"`
	PersonalIDNumber *string                     `form:"personal_id_number"`
	Verification     *IdentityVerificationParams `form:"verification"`
}

// IdentityVerificationParams represents a verification during account creation/updates.
type IdentityVerificationParams struct {
	Document     *string `form:"document"`
	DocumentBack *string `form:"document_back"`
}

// AccountListParams are the parameters allowed during account listing.
type AccountListParams struct {
	ListParams `form:"*"`
}

// AccountExternalAccountParams are the parameters allowed to reference an
// external account when creating an account. It should either have Token set
// or everything else.
type AccountExternalAccountParams struct {
	Params            `form:"*"`
	AccountNumber     *string `form:"account_number"`
	AccountHolderName *string `form:"account_holder_name"`
	AccountHolderType *string `form:"account_holder_type"`
	Country           *string `form:"country"`
	Currency          *string `form:"currency"`
	RoutingNumber     *string `form:"routing_number"`
	Token             *string `form:"token"`
}

// AppendTo implements custom encoding logic for AccountExternalAccountParams
// so that we can send the special required `object` field up along with the
// other specified parameters or the token value.
func (p *AccountExternalAccountParams) AppendTo(body *form.Values, keyParts []string) {
	if p.Token != nil {
		body.Add(form.FormatKey(keyParts), StringValue(p.Token))
	} else {
		body.Add(form.FormatKey(append(keyParts, "object")), "bank_account")
	}
}

// PayoutScheduleParams are the parameters allowed for payout schedules.
type PayoutScheduleParams struct {
	DelayDays        *int64  `form:"delay_days"`
	DelayDaysMinimum *bool   `form:"-"` // See custom AppendTo
	Interval         *string `form:"interval"`
	MonthlyAnchor    *int64  `form:"monthly_anchor"`
	WeeklyAnchor     *string `form:"weekly_anchor"`
}

// AppendTo implements custom encoding logic for PayoutScheduleParams
// so that we can send a special value for `delay_days` field if needed.
func (p *PayoutScheduleParams) AppendTo(body *form.Values, keyParts []string) {
	if BoolValue(p.DelayDaysMinimum) {
		body.Add(form.FormatKey(append(keyParts, "delay_days")), "minimum")
	}
}

// Account is the resource representing your Stripe account.
// For more details see https://stripe.com/docs/api/#account.
type Account struct {
	BusinessLogo          string                  `json:"business_logo"`
	BusinessName          string                  `json:"business_name"`
	BusinessPrimaryColor  string                  `json:"business_primary_color"`
	BusinessURL           string                  `json:"business_url"`
	ChargesEnabled        bool                    `json:"charges_enabled"`
	Country               string                  `json:"country"`
	Created               int64                   `json:"created"`
	DeclineChargeOn       *AccountDeclineSettings `json:"decline_charge_on"`
	DebitNegativeBalances bool                    `json:"debit_negative_balances"`
	DefaultCurrency       Currency                `json:"default_currency"`
	Deleted               bool                    `json:"deleted"`
	DetailsSubmitted      bool                    `json:"details_submitted"`
	Email                 string                  `json:"email"`
	ExternalAccounts      *ExternalAccountList    `json:"external_accounts"`
	ID                    string                  `json:"id"`

	Keys *struct {
		Publishable string `json:"publishable"`
		Secret      string `json:"secret"`
	} `json:"keys"`

	LegalEntity               *LegalEntity      `json:"legal_entity"`
	Metadata                  map[string]string `json:"metadata"`
	DisplayName               string            `json:"display_name"`
	PayoutSchedule            *PayoutSchedule   `json:"payout_schedule"`
	PayoutStatementDescriptor string            `json:"payout_statement_descriptor"`
	PayoutsEnabled            bool              `json:"payouts_enabled"`
	ProductDescription        string            `json:"product_description"`
	StatementDescriptor       string            `json:"statement_descriptor"`
	SupportAddress            *Address          `json:"support_address"`
	SupportEmail              string            `json:"support_email"`
	SupportPhone              string            `json:"support_phone"`
	SupportURL                string            `json:"support_url"`
	Timezone                  string            `json:"timezone"`

	TOSAcceptance *struct {
		Date      int64  `json:"date"`
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
	} `json:"tos_acceptance"`

	Type AccountType `json:"type"`

	Verification *struct {
		DisabledReason IdentityVerificationDisabledReason `json:"disabled_reason"`
		DueBy          int64                              `json:"due_by"`
		FieldsNeeded   []string                           `json:"fields_needed"`
	} `json:"verification"`
}

// UnmarshalJSON handles deserialization of an account.
// This custom unmarshaling is needed because the resulting
// property may be an ID or the full struct if it was expanded.
func (a *Account) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		a.ID = id
		return nil
	}

	type account Account
	var v account
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = Account(v)
	return nil
}

// AccountList is a list of accounts as returned from a list endpoint.
type AccountList struct {
	ListMeta
	Data []*Account `json:"data"`
}

// ExternalAccountList is a list of external accounts that may be either bank
// accounts or cards.
type ExternalAccountList struct {
	ListMeta

	// Values contains any external accounts (bank accounts and/or cards)
	// currently attached to this account.
	Data []*ExternalAccount `json:"data"`
}

// ExternalAccount is an external account (a bank account or card) that's
// attached to an account. It contains fields that will be conditionally
// populated depending on its type.
type ExternalAccount struct {
	// BankAccount is a bank account attached to an account. Populated only if
	// the external account is a bank account.
	BankAccount *BankAccount

	// Card is a card attached to an account. Populated only if the external
	// account is a card.
	Card *Card

	ID   string              `json:"id"`
	Type ExternalAccountType `json:"object"`
}

// UnmarshalJSON implements Unmarshaler.UnmarshalJSON.
func (ea *ExternalAccount) UnmarshalJSON(data []byte) error {
	type externalAccount ExternalAccount
	var v externalAccount
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var err error
	*ea = ExternalAccount(v)

	switch ea.Type {
	case ExternalAccountTypeBankAccount:
		err = json.Unmarshal(data, &ea.BankAccount)
	case ExternalAccountTypeCard:
		err = json.Unmarshal(data, &ea.Card)
	}

	return err
}

// LegalEntity is the structure for properties related to an account's legal state.
type LegalEntity struct {
	AdditionalOwners         []*AdditionalOwner    `json:"additional_owners"`
	Address                  *AccountAddress       `json:"address"`
	AddressKana              *AccountAddress       `json:"address_kana"`
	AddressKanji             *AccountAddress       `json:"address_kanji"`
	BusinessName             string                `json:"business_name"`
	BusinessNameKana         string                `json:"business_name_kana"`
	BusinessNameKanji        string                `json:"business_name_kanji"`
	BusinessTaxIDProvided    bool                  `json:"business_tax_id_provided"`
	BusinessVATIDProvided    bool                  `json:"business_vat_id_provided"`
	DOB                      *DOB                  `json:"dob"`
	FirstName                string                `json:"first_name"`
	FirstNameKana            string                `json:"first_name_kana"`
	FirstNameKanji           string                `json:"first_name_kanji"`
	Gender                   string                `json:"gender"`
	LastName                 string                `json:"last_name"`
	LastNameKana             string                `json:"last_name_kana"`
	LastNameKanji            string                `json:"last_name_kanji"`
	MaidenName               string                `json:"maiden_name"`
	PersonalAddress          *AccountAddress       `json:"personal_address"`
	PersonalAddressKana      *AccountAddress       `json:"personal_address_kana"`
	PersonalAddressKanji     *AccountAddress       `json:"personal_address_kanji"`
	PersonalIDNumberProvided bool                  `json:"personal_id_number_provided"`
	PhoneNumber              string                `json:"phone_number"`
	SSNLast4Provided         bool                  `json:"ssn_last_4_provided"`
	Type                     LegalEntityType       `json:"type"`
	Verification             *IdentityVerification `json:"verification"`
}

// AccountAddress is the structure for an account address.
type AccountAddress struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`

	// Town/cho-me. Note that this is only used for Kana/Kanji representations
	// of an address.
	Town string `json:"town"`
}

// AccountDeclineSettings is the structure for an account's decline settings.
type AccountDeclineSettings struct {
	AVSFailure bool `json:"avs_failure"`
	CVCFailure bool `json:"cvc_failure"`
}

// DOB is a structure for an account owner's date of birth.
type DOB struct {
	Day   int64 `json:"day"`
	Month int64 `json:"month"`
	Year  int64 `json:"year"`
}

// AdditionalOwner is the structure for an account owner.
type AdditionalOwner struct {
	Address                  *AccountAddress      `json:"address"`
	DOB                      DOB                  `json:"dob"`
	FirstName                string               `json:"first_name"`
	LastName                 string               `json:"last_name"`
	MaidenName               string               `json:"maiden_name"`
	PersonalIDNumberProvided bool                 `json:"personal_id_number_provided"`
	Verification             IdentityVerification `json:"verification"`
}

// IdentityVerification is the structure for an account's verification.
type IdentityVerification struct {
	Details      string                          `json:"details"`
	DetailsCode  IdentityVerificationDetailsCode `json:"details_code"`
	Document     *File                           `json:"document"`
	DocumentBack *File                           `json:"document_back"`
	Status       IdentityVerificationStatus      `json:"status"`
}

// PayoutSchedule is the structure for an account's payout schedule.
type PayoutSchedule struct {
	DelayDays     int64          `json:"delay_days"`
	Interval      PayoutInterval `json:"interval"`
	MonthlyAnchor int64          `json:"monthly_anchor"`
	WeeklyAnchor  string         `json:"weekly_anchor"`
}

// AccountRejectParams is the structure for the Reject function.
type AccountRejectParams struct {
	Params `form:"*"`
	Reason *string `form:"reason"`
}
