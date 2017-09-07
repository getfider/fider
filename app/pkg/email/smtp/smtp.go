package smtp

import (
	"fmt"
	gosmtp "net/smtp"
)

//Sender is used to send e-mails
type Sender struct {
	host     string
	port     string
	username string
	password string
}

//NewSender creates a new mailgun e-mail sender
func NewSender(host, port, username, password string) *Sender {
	return &Sender{host, port, username, password}
}

//Send an e-mail
func (s *Sender) Send(from, to, subject, message string) error {

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	body := ""
	for k, v := range headers {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body += "\r\n" + message

	servername := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := gosmtp.PlainAuth("", s.username, s.password, s.host)
	return gosmtp.SendMail(servername, auth, from, []string{to}, []byte(body))
}
