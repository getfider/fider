package stripe

import (
	"encoding/json"
)

// PaymentIntentCancellationReason is the list of allowed values for the cancelation reason.
type PaymentIntentCancellationReason string

// List of values that PaymentIntentCaptureMethod can take.
const (
	PaymentIntentCancellationReasonDuplicate           PaymentIntentCancellationReason = "duplicate"
	PaymentIntentCancellationReasonFraudulent          PaymentIntentCancellationReason = "fraudulent"
	PaymentIntentCancellationReasonRequestedByCustomer PaymentIntentCancellationReason = "requested_by_customer"
)

// PaymentIntentCaptureMethod is the list of allowed values for the capture method.
type PaymentIntentCaptureMethod string

// List of values that PaymentIntentCaptureMethod can take.
const (
	PaymentIntentCaptureMethodAutomatic PaymentIntentCaptureMethod = "automatic"
	PaymentIntentCaptureMethodManual    PaymentIntentCaptureMethod = "manual"
)

// PaymentIntentConfirmationMethod is the list of allowed values for the confirmation method.
type PaymentIntentConfirmationMethod string

// List of values that PaymentIntentConfirmationMethod can take.
const (
	PaymentIntentConfirmationMethodPublishable PaymentIntentConfirmationMethod = "publishable"
	PaymentIntentConfirmationMethodSecret      PaymentIntentConfirmationMethod = "secret"
)

// PaymentIntentNextActionType is the list of allowed values for the next action's type.
type PaymentIntentNextActionType string

// List of values that PaymentIntentNextActionType can take.
const (
	PaymentIntentNextActionAuthorizeWithURL PaymentIntentNextActionType = "authorize_with_url"
)

// PaymentIntentStatus is the list of allowed values for the payment intent's status.
type PaymentIntentStatus string

// List of values that PaymentIntentStatus can take.
const (
	PaymentIntentStatusCanceled             PaymentIntentStatus = "canceled"
	PaymentIntentStatusProcessing           PaymentIntentStatus = "processing"
	PaymentIntentStatusRequiresCapture      PaymentIntentStatus = "requires_capture"
	PaymentIntentStatusRequiresConfirmation PaymentIntentStatus = "requires_confirmation"
	PaymentIntentStatusRequiresSource       PaymentIntentStatus = "requires_source"
	PaymentIntentStatusRequiresSourceAction PaymentIntentStatus = "requires_source_action"
	PaymentIntentStatusSucceeded            PaymentIntentStatus = "succeeded"
)

// PaymentIntentCancelParams is the set of parameters that can be used when canceling a payment intent.
type PaymentIntentCancelParams struct {
	Params             `form:"*"`
	CancellationReason *string `form:"cancellation_reason"`
}

// PaymentIntentCaptureParams is the set of parameters that can be used when capturing a payment intent.
type PaymentIntentCaptureParams struct {
	Params               `form:"*"`
	AmountToCapture      *int64                           `form:"amount_to_capture"`
	ApplicationFeeAmount *int64                           `form:"application_fee_amount"`
	TransferData         *PaymentIntentTransferDataParams `form:"transfer_data"`
}

// PaymentIntentConfirmParams is the set of parameters that can be used when confirming a payment intent.
type PaymentIntentConfirmParams struct {
	Params               `form:"*"`
	ReceiptEmail         *string                `form:"receipt_email"`
	ReturnURL            *string                `form:"return_url"`
	SaveSourceToCustomer *bool                  `form:"save_source_to_customer"`
	Shipping             *ShippingDetailsParams `form:"shipping"`
	Source               *string                `form:"source"`
}

// PaymentIntentTransferDataParams is the set of parameters allowed for the transfer hash.
type PaymentIntentTransferDataParams struct {
	Amount      *int64  `form:"amount"`
	Destination *string `form:"destination"`
}

// PaymentIntentParams is the set of parameters that can be used when handling a payment intent.
type PaymentIntentParams struct {
	Params               `form:"*"`
	AllowedSourceTypes   []*string                        `form:"allowed_source_types"`
	Amount               *int64                           `form:"amount"`
	ApplicationFeeAmount *int64                           `form:"application_fee_amount"`
	Confirm              *bool                            `form:"confirm"`
	CaptureMethod        *string                          `form:"capture_method"`
	Currency             *string                          `form:"currency"`
	Customer             *string                          `form:"customer"`
	Description          *string                          `form:"description"`
	OnBehalfOf           *string                          `form:"on_behalf_of"`
	ReceiptEmail         *string                          `form:"receipt_email"`
	ReturnURL            *string                          `form:"return_url"`
	SaveSourceToCustomer *bool                            `form:"save_source_to_customer"`
	Shipping             *ShippingDetailsParams           `form:"shipping"`
	Source               *string                          `form:"source"`
	StatementDescriptor  *string                          `form:"statement_descriptor"`
	TransferData         *PaymentIntentTransferDataParams `form:"transfer_data"`
	TransferGroup        *string                          `form:"transfer_group"`
}

// PaymentIntentListParams is the set of parameters that can be used when listing payment intents.
// For more details see https://stripe.com/docs/api#list_payouts.
type PaymentIntentListParams struct {
	ListParams `form:"*"`
}

// PaymentIntentLastPaymentError represents the last error happening on a payment intent.
type PaymentIntentLastPaymentError struct {
	Charge      string         `json:"charge"`
	Code        string         `json:"code"`
	DeclineCode string         `json:"decline_code"`
	DocURL      string         `json:"doc_url"`
	Message     string         `json:"message"`
	Param       string         `json:"param"`
	Source      *PaymentSource `json:"source"`
	Type        ErrorType      `json:"type"`
}

// PaymentIntentSourceActionAuthorizeWithURL represents the resource for the next action of type
// "authorize_with_url".
type PaymentIntentSourceActionAuthorizeWithURL struct {
	ReturnURL string `json:"return_url"`
	URL       string `json:"url"`
}

// PaymentIntentSourceAction represents the type of action to take on a payment intent.
type PaymentIntentSourceAction struct {
	AuthorizeWithURL *PaymentIntentSourceActionAuthorizeWithURL `json:"authorize_with_url"`
	Type             PaymentIntentNextActionType                `json:"type"`
}

// PaymentIntentTransferData represents the information for the transfer associated with a payment intent.
type PaymentIntentTransferData struct {
	Amount      int64    `json:"amount"`
	Destination *Account `json:"destination"`
}

// PaymentIntent is the resource representing a Stripe payout.
// For more details see https://stripe.com/docs/api#payment_intents.
type PaymentIntent struct {
	AllowedSourceTypes  []string                        `json:"allowed_source_types"`
	Amount              int64                           `json:"amount"`
	AmountCapturable    int64                           `json:"amount_capturable"`
	AmountReceived      int64                           `json:"amount_received"`
	Application         *Application                    `json:"application"`
	ApplicationFee      int64                           `json:"application_fee"`
	CanceledAt          int64                           `json:"canceled_at"`
	CaptureMethod       PaymentIntentCaptureMethod      `json:"capture_method"`
	Charges             *ChargeList                     `json:"charges"`
	ClientSecret        string                          `json:"client_secret"`
	ConfirmationMethod  PaymentIntentConfirmationMethod `json:"confirmation_method"`
	Created             int64                           `json:"created"`
	Currency            string                          `json:"currency"`
	Customer            *Customer                       `json:"customer"`
	Description         string                          `json:"description"`
	LastPaymentError    *PaymentIntentLastPaymentError  `json:"last_payment_error"`
	Livemode            bool                            `json:"livemode"`
	ID                  string                          `json:"id"`
	Metadata            map[string]string               `json:"metadata"`
	NextSourceAction    *PaymentIntentSourceAction      `json:"next_source_action"`
	OnBehalfOf          *Account                        `json:"on_behalf_of"`
	ReceiptEmail        string                          `json:"receipt_email"`
	Review              *Review                         `json:"review"`
	Shipping            ShippingDetails                 `json:"shipping"`
	Source              *PaymentSource                  `json:"source"`
	StatementDescriptor string                          `json:"statement_descriptor"`
	Status              PaymentIntentStatus             `json:"status"`
	TransferData        *PaymentIntentTransferData      `json:"transfer_data"`
	TransferGroup       string                          `json:"transfer_group"`
}

// PaymentIntentList is a list of payment intents as retrieved from a list endpoint.
type PaymentIntentList struct {
	ListMeta
	Data []*PaymentIntent `json:"data"`
}

// UnmarshalJSON handles deserialization of a Payment Intent.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (p *PaymentIntent) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		p.ID = id
		return nil
	}

	type paymentintent PaymentIntent
	var v paymentintent
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PaymentIntent(v)
	return nil
}
