package vat

import (
	"errors"
	"time"
)

var (
	// ErrInvalidVATNumberFormat is returned when given VAT Number is in a invalid format
	ErrInvalidVATNumberFormat = errors.New("vat: invalid VAT number format")
	// ErrServiceUnreachable is returned when we cannot connect to VIES VAT Service
	ErrServiceUnreachable = errors.New("vat: VIES service is unreachable")
	// ErrVIESInvalidInput is returned based on VIES response
	ErrVIESInvalidInput = errors.New("vat: VAT number is invalid")
	// ErrVIESGlobalMaxConcurrentRequest is returned based on VIES response
	ErrVIESGlobalMaxConcurrentRequest = errors.New("vat: the maximum number of global concurrent requests has been reached, try again later")
	// ErrVIESMSMaxConcurrentRequest is returned based on VIES response
	ErrVIESMSMaxConcurrentRequest = errors.New("vat: the maximum number of concurrent requests for this member state has been reached, try again later")
	// ErrVIESServiceUnavailable is returned based on VIES response
	ErrVIESServiceUnavailable = errors.New("vat: an error was encountered either at the network level or the web application level, try again later")
	// ErrVIESMSUnavailable is returned based on VIES response
	ErrVIESMSUnavailable = errors.New("vat: the application at the member state is not replying or not available, try again later")
	// ErrVIESTimeout is returned based on VIES response
	ErrVIESTimeout = errors.New("vat: the application did not receive a reply within the allocated time period, try again later")
)

// Response is the model returned from VIES query
type Response struct {
	IsValid     bool
	VATNumber   string
	CountryCode string
	RequestDate time.Time
	Name        string
	Address     string
}

// Query VIES service to check if VAT is valid
func Query(vatNumber string) (*Response, error) {
	vatNumber = sanitizeVATNumber(vatNumber)

	isValidFormat, _ := ValidateNumberFormat(vatNumber)
	if !isValidFormat {
		return nil, ErrInvalidVATNumberFormat
	}

	var (
		response *Response
		err      error
	)

	retry(4, 1*time.Second, func() bool {
		response, err = sendRequestVIES(vatNumber)
		// only retry if there's an error and it's different than INVALID_INPUT
		return err != nil && err != ErrVIESInvalidInput
	})

	return response, err
}
