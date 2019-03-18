package email

import (
	"net/mail"

	"github.com/getfider/fider/app/pkg/env"
)

// Params used to replace variables on emails
type Params map[string]interface{}

// Merge given params into current params
func (p Params) Merge(p2 Params) Params {
	if p == nil {
		p = Params{}
	}
	for k, v := range p2 {
		p[k] = v
	}
	return p
}

// NoReply is the default 'from' address
var NoReply = env.Config.Email.NoReply

// Recipient contains details of who is receiving the email
type Recipient struct {
	Name    string
	Address string
	Params  Params
}

// NewRecipient creates a new Recipient
func NewRecipient(name, address string, params Params) Recipient {
	return Recipient{
		Name:    name,
		Address: address,
		Params:  params,
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

type SendMessageCommand struct {
	From         string
	To           []Recipient
	TemplateName string
	Params       Params
}
