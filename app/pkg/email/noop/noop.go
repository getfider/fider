package noop

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
)

type request struct {
	Tenant       *models.Tenant
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
func (s *Sender) Send(tenant *models.Tenant, templateName string, params email.Params, from string, to email.Recipient) error {
	s.Requests = append(s.Requests, &request{
		Tenant:       tenant,
		TemplateName: templateName,
		Params:       params,
		From:         from,
		To:           []email.Recipient{to},
	})
	return nil

}

// BatchSend an email to multiple recipients
func (s *Sender) BatchSend(tenant *models.Tenant, templateName string, params email.Params, from string, to []email.Recipient) error {
	if len(to) > 0 {
		s.Requests = append(s.Requests, &request{
			Tenant:       tenant,
			TemplateName: templateName,
			Params:       params,
			From:         from,
			To:           to,
		})
	}

	return nil
}
