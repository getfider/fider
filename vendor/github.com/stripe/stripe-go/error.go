package stripe

import "encoding/json"

// ErrorType is the list of allowed values for the error's type.
type ErrorType string

// List of values that ErrorType can take.
const (
	ErrorTypeAPI            ErrorType = "api_error"
	ErrorTypeAPIConnection  ErrorType = "api_connection_error"
	ErrorTypeAuthentication ErrorType = "authentication_error"
	ErrorTypeCard           ErrorType = "card_error"
	ErrorTypeInvalidRequest ErrorType = "invalid_request_error"
	ErrorTypePermission     ErrorType = "more_permissions_required"
	ErrorTypeRateLimit      ErrorType = "rate_limit_error"
)

// ErrorCode is the list of allowed values for the error's code.
type ErrorCode string

// List of values that ErrorCode can take.
const (
	ErrorCodeCardDeclined       ErrorCode = "card_declined"
	ErrorCodeExpiredCard        ErrorCode = "expired_card"
	ErrorCodeIncorrectCVC       ErrorCode = "incorrect_cvc"
	ErrorCodeIncorrectZip       ErrorCode = "incorrect_zip"
	ErrorCodeIncorrectNumber    ErrorCode = "incorrect_number"
	ErrorCodeInvalidCVC         ErrorCode = "invalid_cvc"
	ErrorCodeInvalidExpiryMonth ErrorCode = "invalid_expiry_month"
	ErrorCodeInvalidExpiryYear  ErrorCode = "invalid_expiry_year"
	ErrorCodeInvalidNumber      ErrorCode = "invalid_number"
	ErrorCodeInvalidSwipeData   ErrorCode = "invalid_swipe_data"
	ErrorCodeMissing            ErrorCode = "missing"
	ErrorCodeProcessingError    ErrorCode = "processing_error"
	ErrorCodeRateLimit          ErrorCode = "rate_limit"
	ErrorCodeResourceMissing    ErrorCode = "resource_missing"
)

// Error is the response returned when a call is unsuccessful.
// For more details see  https://stripe.com/docs/api#errors.
type Error struct {
	ChargeID string    `json:"charge,omitempty"`
	Code     ErrorCode `json:"code,omitempty"`

	// Err contains an internal error with an additional level of granularity
	// that can be used in some cases to get more detailed information about
	// what went wrong. For example, Err may hold a CardError that indicates
	// exactly what went wrong during charging a card.
	Err error `json:"-"`

	HTTPStatusCode int       `json:"status,omitempty"`
	Msg            string    `json:"message"`
	Param          string    `json:"param,omitempty"`
	RequestID      string    `json:"request_id,omitempty"`
	Type           ErrorType `json:"type"`
}

// Error serializes the error object to JSON and returns it as a string.
func (e *Error) Error() string {
	ret, _ := json.Marshal(e)
	return string(ret)
}

// APIConnectionError is a failure to connect to the Stripe API.
type APIConnectionError struct {
	stripeErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *APIConnectionError) Error() string {
	return e.stripeErr.Error()
}

// APIError is a catch all for any errors not covered by other types (and
// should be extremely uncommon).
type APIError struct {
	stripeErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *APIError) Error() string {
	return e.stripeErr.Error()
}

// AuthenticationError is a failure to properly authenticate during a request.
type AuthenticationError struct {
	stripeErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *AuthenticationError) Error() string {
	return e.stripeErr.Error()
}

// PermissionError results when you attempt to make an API request
// for which your API key doesn't have the right permissions.
type PermissionError struct {
	stripeErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *PermissionError) Error() string {
	return e.stripeErr.Error()
}

// CardError are the most common type of error you should expect to handle.
// They result when the user enters a card that can't be charged for some
// reason.
type CardError struct {
	stripeErr   *Error
	DeclineCode string `json:"decline_code,omitempty"`
}

// Error serializes the error object to JSON and returns it as a string.
func (e *CardError) Error() string {
	return e.stripeErr.Error()
}

// InvalidRequestError is an error that occurs when a request contains invalid
// parameters.
type InvalidRequestError struct {
	stripeErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *InvalidRequestError) Error() string {
	return e.stripeErr.Error()
}

// RateLimitError occurs when the Stripe API is hit to with too many requests
// too quickly and indicates that the current request has been rate limited.
type RateLimitError struct {
	stripeErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *RateLimitError) Error() string {
	return e.stripeErr.Error()
}

// rawError deserializes the outer JSON object returned in an error response
// from the API.
type rawError struct {
	E *rawErrorInternal `json:"error,omitempty"`
}

// rawErrorInternal embeds Error to deserialize all the standard error fields,
// but also adds other fields that may or may not be present depending on error
// type to help with deserialization. (e.g. Declinecode).
type rawErrorInternal struct {
	*Error
	DeclineCode *string `json:"decline_code,omitempty"`
}
