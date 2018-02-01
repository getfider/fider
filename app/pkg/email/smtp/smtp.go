package smtp

import (
	"fmt"
	gosmtp "net/smtp"

	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/log"
)

//Sender is used to send e-mails
type Sender struct {
	logger   log.Logger
	host     string
	port     string
	username string
	password string
}

//NewSender creates a new mailgun e-mail sender
func NewSender(logger log.Logger, host, port, username, password string) *Sender {
	return &Sender{logger, host, port, username, password}
}

//Send an e-mail
func (s *Sender) Send(from, to, templateName string, params map[string]interface{}) error {

	message := email.RenderMessage(templateName, params)
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", from, email.NoReply)
	headers["To"] = to
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
	err := gosmtp.SendMail(servername, auth, email.NoReply, []string{to}, []byte(body))
	if err == nil {
		s.logger.Debugf("E-mail sent.")
	} else {
		s.logger.Errorf("Failed to send e-mail")
	}
	return err
}
