package paddle

import "encoding/json"

type PaddleResponse struct {
	IsSuccess bool            `json:"success"`
	Response  json.RawMessage `json:"response"`
	Error     struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
}

type PaddleSubscriptionItem struct {
	SignupDate         string `json:"signup_date"`
	UpdateURL          string `json:"update_url"`
	CancelURL          string `json:"cancel_url"`
	PaymentInformation struct {
		PaymentMethod  string `json:"payment_method"`
		CardType       string `json:"card_type"`
		LastFourDigits string `json:"last_four_digits"`
		ExpiryDate     string `json:"expiry_date"`
	} `json:"payment_information"`
	LastPayment struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
		Date     string  `json:"date"`
	} `json:"last_payment"`
}
