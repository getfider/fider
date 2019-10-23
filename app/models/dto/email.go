package dto

import "net/mail"

// Recipient contains details of who is receiving the email
type Recipient struct {
	Name    string
	Address string
	Props   Props
}

// NewRecipient creates a new Recipient
func NewRecipient(name, address string, props Props) Recipient {
	return Recipient{
		Name:    name,
		Address: address,
		Props:   props,
	}
}

// Strings returns the RFC format to send emails via SMTP
func (r Recipient) String() string {
	if r.Address == "" {
		return ""
	}

	address := mail.Address{
		Name:    r.Name,
		Address: r.Address,
	}

	return address.String()
}
