package stripe

import "encoding/json"

// RecipientTransferDestinationType consts represent valid recipient_transfer destinations.
type RecipientTransferDestinationType string

// List of values that RecipientTransferDestinationType can take.
const (
	RecipientTransferDestinationBankAccount RecipientTransferDestinationType = "bank_account"
	RecipientTransferDestinationCard        RecipientTransferDestinationType = "card"
)

// RecipientTransferFailureCode is the list of allowed values for the recipient_transfer's failure code.
type RecipientTransferFailureCode string

// List of values that RecipientTransferFailureCode can take.
const (
	RecipientTransferFailureCodeAccountClosed         RecipientTransferFailureCode = "account_closed"
	RecipientTransferFailureCodeAccountFrozen         RecipientTransferFailureCode = "account_frozen"
	RecipientTransferFailureCodeBankAccountRestricted RecipientTransferFailureCode = "bank_account_restricted"
	RecipientTransferFailureCodeBankOwnershipChanged  RecipientTransferFailureCode = "bank_ownership_changed"
	RecipientTransferFailureCodeDebitNotAuthorized    RecipientTransferFailureCode = "debit_not_authorized"
	RecipientTransferFailureCodeCouldNotProcess       RecipientTransferFailureCode = "could_not_process"
	RecipientTransferFailureCodeInsufficientFunds     RecipientTransferFailureCode = "insufficient_funds"
	RecipientTransferFailureCodeInvalidAccountNumber  RecipientTransferFailureCode = "invalid_account_number"
	RecipientTransferFailureCodeInvalidCurrency       RecipientTransferFailureCode = "invalid_currency"
	RecipientTransferFailureCodeNoAccount             RecipientTransferFailureCode = "no_account"
)

// RecipientTransferSourceType is the list of allowed values for the recipient_transfer's source_type field.
type RecipientTransferSourceType string

// List of values that RecipientTransferSourceType can take.
const (
	RecipientTransferSourceTypeAlipayAccount   RecipientTransferSourceType = "alipay_account"
	RecipientTransferSourceTypeBankAccount     RecipientTransferSourceType = "bank_account"
	RecipientTransferSourceTypeBitcoinReceiver RecipientTransferSourceType = "bitcoin_receiver"
	RecipientTransferSourceTypeCard            RecipientTransferSourceType = "card"
)

// RecipientTransferStatus is the list of allowed values for the recipient_transfer's status.
type RecipientTransferStatus string

// List of values that RecipientTransferStatus can take.
const (
	RecipientTransferStatusFailed    RecipientTransferStatus = "failed"
	RecipientTransferStatusInTransit RecipientTransferStatus = "in_transit"
	RecipientTransferStatusPaid      RecipientTransferStatus = "paid"
	RecipientTransferStatusPending   RecipientTransferStatus = "pending"
)

// RecipientTransferType is the list of allowed values for the recipient_transfer's type.
type RecipientTransferType string

// List of values that RecipientTransferType can take.
const (
	RecipientTransferTypeBankAccount RecipientTransferType = "bank_account"
	RecipientTransferTypeCard        RecipientTransferType = "card"
)

// RecipientTransferMethodType represents the type of recipient_transfer
type RecipientTransferMethodType string

// List of values that RecipientTransferMethodType can take.
const (
	RecipientTransferMethodInstant  RecipientTransferMethodType = "instant"
	RecipientTransferMethodStandard RecipientTransferMethodType = "standard"
)

// RecipientTransferDestination describes the destination of a RecipientTransfer.
// The Type should indicate which object is fleshed out
// For more details see https://stripe.com/docs/api/go#recipient_transfer_object
type RecipientTransferDestination struct {
	BankAccount *BankAccount                     `json:"-"`
	Card        *Card                            `json:"-"`
	ID          string                           `json:"id"`
	Type        RecipientTransferDestinationType `json:"object"`
}

// RecipientTransfer is the resource representing a Stripe recipient_transfer.
// For more details see https://stripe.com/docs/api#recipient_transfers.
type RecipientTransfer struct {
	Amount              int64                        `json:"amount"`
	AmountReversed      int64                        `json:"amount_reversed"`
	BalanceTransaction  *BalanceTransaction          `json:"balance_transaction"`
	BankAccount         *BankAccount                 `json:"bank_account"`
	Card                *Card                        `json:"card"`
	Created             int64                        `json:"created"`
	Currency            Currency                     `json:"currency"`
	Date                int64                        `json:"date"`
	Description         string                       `json:"description"`
	Destination         string                       `json:"destination"`
	FailureCode         RecipientTransferFailureCode `json:"failure_code"`
	FailureMessage      string                       `json:"failure_message"`
	ID                  string                       `json:"id"`
	Livemode            bool                         `json:"livemode"`
	Metadata            map[string]string            `json:"metadata"`
	Method              RecipientTransferMethodType  `json:"method"`
	Recipient           *Recipient                   `json:"recipient"`
	Reversals           *ReversalList                `json:"reversals"`
	Reversed            bool                         `json:"reversed"`
	SourceTransaction   *BalanceTransactionSource    `json:"source_transaction"`
	SourceType          RecipientTransferSourceType  `json:"source_type"`
	StatementDescriptor string                       `json:"statement_descriptor"`
	Status              RecipientTransferStatus      `json:"status"`
	Type                RecipientTransferType        `json:"type"`
}

// UnmarshalJSON handles deserialization of a RecipientTransfer.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (t *RecipientTransfer) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		t.ID = id
		return nil
	}

	type recipientTransfer RecipientTransfer
	var v recipientTransfer
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*t = RecipientTransfer(v)
	return nil
}

// UnmarshalJSON handles deserialization of a RecipientTransferDestination.
// This custom unmarshaling is needed because the specific
// type of destination it refers to is specified in the JSON
func (d *RecipientTransferDestination) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		d.ID = id
		return nil
	}

	type recipientTransferDestination RecipientTransferDestination
	var v recipientTransferDestination
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var err error
	*d = RecipientTransferDestination(v)

	switch d.Type {
	case RecipientTransferDestinationBankAccount:
		err = json.Unmarshal(data, &d.BankAccount)
	case RecipientTransferDestinationCard:
		err = json.Unmarshal(data, &d.Card)
	}

	return err
}

// MarshalJSON handles serialization of a RecipientTransferDestination.
// This custom marshaling is needed because we can only send a string
// ID as a destination, even though it can be expanded to a full
// object when retrieving
func (d *RecipientTransferDestination) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.ID)
}
