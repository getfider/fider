package noop

import (
	"github.com/getfider/fider/app/pkg/email"
)

type request struct {
	Context      email.Context
	TemplateName string
	Params       email.Params
	From         string
	To           []email.Recipient
}

// Sender does not send emails
type Sender struct {
	Requests []*request
}

// NewSender creates a new NoopSender
func NewSender() *Sender {
	return &Sender{
		Requests: make([]*request, 0),
	}
}

// Send an email
func (s *Sender) Send(ctx email.Context, templateName string, params email.Params, from string, to email.Recipient) error {
	s.Requests = append(s.Requests, &request{
		Context:      ctx,
		TemplateName: templateName,
		Params:       params,
		From:         from,
		To:           []email.Recipient{to},
	})
	return nil

}

// BatchSend an email to multiple recipients
func (s *Sender) BatchSend(ctx email.Context, templateName string, params email.Params, from string, to []email.Recipient) error {
	if len(to) > 0 {
		s.Requests = append(s.Requests, &request{
			Context:      ctx,
			TemplateName: templateName,
			Params:       params,
			From:         from,
			To:           to,
		})
	}

	return nil
}
