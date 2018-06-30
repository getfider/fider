package smtp

import (
	"fmt"
	gosmtp "net/smtp"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
)

func authenticate(username string, password string, host string) gosmtp.Auth {
	if username == "" && password == "" {
		return nil
	}
	return gosmtp.PlainAuth("", username, password, host)
}

type builder struct {
	content string
}

func (b *builder) Set(key, value string) {
	b.content += fmt.Sprintf("%s: %s\r\n", key, value)
}

func (b *builder) Body(body string) {
	b.content += "\r\n" + body
}

func (b *builder) Bytes() []byte {
	return []byte(b.content)
}

//Sender is used to send emails
type Sender struct {
	logger   log.Logger
	host     string
	port     string
	username string
	password string
	send     func(string, gosmtp.Auth, string, []string, []byte) error
}

//NewSender creates a new mailgun email sender
func NewSender(logger log.Logger, host, port, username, password string) *Sender {
	return &Sender{logger, host, port, username, password, gosmtp.SendMail}
}

//ReplaceSend can be used to mock internal send function
func (s *Sender) ReplaceSend(send func(string, gosmtp.Auth, string, []string, []byte) error) {
	s.send = send
}

//Send an email
func (s *Sender) Send(tenant *models.Tenant, templateName string, params email.Params, from string, to email.Recipient) error {
	if to.Address == "" {
		return nil
	}

	if !email.CanSendTo(to.Address) {
		s.logger.Warnf("Skipping email to '@{Name} <@{Address}>' due to whitelist.", log.Props{
			"Name":    to.Name,
			"Address": to.Address,
		})
		return nil
	}

	s.logger.Debugf("Sending email to @{Address} with template @{TemplateName} and params @{Params}.", log.Props{
		"Address":      to.Address,
		"TemplateName": templateName,
		"Params":       to.Params,
	})

	message := email.RenderMessage(templateName, params.Merge(to.Params))
	b := builder{}
	b.Set("From", fmt.Sprintf("%s <%s>", from, email.NoReply))
	b.Set("To", fmt.Sprintf("%s <%s>", to.Name, to.Address))
	b.Set("Subject", message.Subject)
	b.Set("MIME-version", "1.0")
	b.Set("Content-Type", "text/html; charset=\"UTF-8\"")
	b.Body(message.Body)

	servername := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := authenticate(s.username, s.password, s.host)
	err := s.send(servername, auth, email.NoReply, []string{to.Address}, b.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to send email with template %s", templateName)
	}
	s.logger.Debug("Email sent.")
	return nil
}

// BatchSend an email to multiple recipients
func (s *Sender) BatchSend(tenant *models.Tenant, templateName string, params email.Params, from string, to []email.Recipient) error {
	for _, r := range to {
		if err := s.Send(tenant, templateName, params, from, r); err != nil {
			return errors.Wrap(err, "failed to batch send email to %d recipients", len(to))
		}
	}
	return nil
}
