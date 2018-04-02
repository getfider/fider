package smtp

import (
	"fmt"
	gosmtp "net/smtp"

	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
)

//Sender is used to send emails
type Sender struct {
	logger   log.Logger
	host     string
	port     string
	username string
	password string
}

//NewSender creates a new mailgun email sender
func NewSender(logger log.Logger, host, port, username, password string) *Sender {
	return &Sender{logger, host, port, username, password}
}

//Send an email
func (s *Sender) Send(templateName string, params email.Params, from string, to email.Recipient) error {
	if !email.CanSendTo(to.Address) {
		s.logger.Warnf("Skipping email to %s <%s> due to whitelist.", to.Name, to.Address)
		return nil
	}

	s.logger.Debugf("Sending email to %s with template %s and params %s.", to.Address, templateName, to.Params)

	message := email.RenderMessage(templateName, params.Merge(to.Params))
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", from, email.NoReply)
	headers["To"] = fmt.Sprintf("%s <%s>", to.Name, to.Address)
	headers["Subject"] = message.Subject
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	body := ""
	for k, v := range headers {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body += "\r\n" + message.Body

	servername := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := gosmtp.PlainAuth("", s.username, s.password, s.host)
	err := gosmtp.SendMail(servername, auth, email.NoReply, []string{to.Address}, []byte(body))
	if err != nil {
		return errors.Wrap(err, "failed to send email with template %s", templateName)
	}
	s.logger.Debugf("Email sent.")
	return nil
}

// BatchSend an email to multiple recipients
func (s *Sender) BatchSend(templateName string, params email.Params, from string, to []email.Recipient) error {
	for _, r := range to {
		if err := s.Send(templateName, params, from, r); err != nil {
			return errors.Wrap(err, "failed to batch send email to %d recipients", len(to))
		}
	}
	return nil
}
