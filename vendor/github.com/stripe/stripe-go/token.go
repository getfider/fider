package stripe

// TokenType is the list of allowed values for a token's type.
type TokenType string

// List of values that TokenType can take.
const (
	TokenTypeAccount     TokenType = "account"
	TokenTypeCard        TokenType = "card"
	TokenTypeBankAccount TokenType = "bank_account"
	TokenTypePII         TokenType = "pii"
)

// TokenParams is the set of parameters that can be used when creating a token.
// For more details see https://stripe.com/docs/api#create_card_token and https://stripe.com/docs/api#create_bank_account_token.
type TokenParams struct {
	Params      `form:"*"`
	BankAccount *BankAccountParams `form:"bank_account"`
	Card        *CardParams        `form:"card"`
	Customer    *string            `form:"customer"`

	// Email is an undocumented parameter used by Stripe Checkout
	// It may be removed from the API without notice.
	Email *string `form:"email"`

	PII *PIIParams `form:"pii"`
}

// Token is the resource representing a Stripe token.
// For more details see https://stripe.com/docs/api#tokens.
type Token struct {
	BankAccount *BankAccount `json:"bank_account"`
	Card        *Card        `json:"card"`
	ClientIP    string       `json:"client_ip"`
	Created     int64        `json:"created"`

	// Email is an undocumented field but included for all tokens created
	// with Stripe Checkout.
	Email string `json:"email"`

	ID       string    `json:"id"`
	Livemode bool      `json:"livemode"`
	Type     TokenType `json:"type"`
	Used     bool      `json:"used"`
}

// PIIParams are parameters for personal identifiable information (PII).
type PIIParams struct {
	Params           `form:"*"`
	PersonalIDNumber *string `form:"personal_id_number"`
}
