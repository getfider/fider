package smtp

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	gosmtp "net/smtp"
	"net/url"
	"strconv"
	"time"

	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
)

func authenticate(username string, password string, host string) gosmtp.Auth {
	if username == "" && password == "" {
		return nil
	}
	return AgnosticAuth("", username, password, host)
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
	send     func(string, string, gosmtp.Auth, string, []string, []byte) error
}

//NewSender creates a new mailgun email sender
func NewSender(logger log.Logger, host, port, username, password string) *Sender {
	return &Sender{logger, host, port, username, password, sendMail}
}

//ReplaceSend can be used to mock internal send function
func (s *Sender) ReplaceSend(send func(string, string, gosmtp.Auth, string, []string, []byte) error) {
	s.send = send
}

//Send an email
func (s *Sender) Send(ctx email.Context, templateName string, params email.Params, from string, to email.Recipient) error {
	if to.Address == "" {
		return nil
	}

	u, err := url.Parse(ctx.BaseURL())
	localname := "localhost"
	if err == nil {
		localname = u.Hostname()
	}

	if !email.CanSendTo(to.Address) {
		s.logger.Warnf("Skipping email to '@{Name} <@{Address}>'.", log.Props{
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

	message := email.RenderMessage(ctx, templateName, params.Merge(to.Params))
	b := builder{}
	b.Set("From", email.NewRecipient(from, email.NoReply, email.Params{}).String())
	b.Set("Reply-To", email.NoReply)
	b.Set("To", to.String())
	b.Set("Subject", message.Subject)
	b.Set("MIME-version", "1.0")
	b.Set("Content-Type", "text/html; charset=\"UTF-8\"")
	b.Set("Date", time.Now().Format(time.RFC1123Z))
	b.Set("Message-ID", generateMessageID(localname))
	b.Body(message.Body)

	servername := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := authenticate(s.username, s.password, s.host)
	err = s.send(localname, servername, auth, email.NoReply, []string{to.Address}, b.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to send email with template %s", templateName)
	}
	s.logger.Debug("Email sent.")
	return nil
}

// BatchSend an email to multiple recipients
func (s *Sender) BatchSend(ctx email.Context, templateName string, params email.Params, from string, to []email.Recipient) error {
	for _, r := range to {
		if err := s.Send(ctx, templateName, params, from, r); err != nil {
			return errors.Wrap(err, "failed to batch send email to %d recipients", len(to))
		}
	}
	return nil
}

func sendMail(localName, serverAddress string, a gosmtp.Auth, from string, to []string, msg []byte) error {
	host, _, _ := net.SplitHostPort(serverAddress)
	c, err := gosmtp.Dial(serverAddress)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello(localName); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}
	if a != nil {
		if ok, _ := c.Extension("AUTH"); !ok {
			return errors.New("smtp: server doesn't support AUTH")
		}
		if err = c.Auth(a); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func generateMessageID(localName string) string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
	  panic(err)
	}
	randStr := hex.EncodeToString(buf)
	messageID := fmt.Sprintf("<%s.%s@%s>", randStr, timestamp, localName)
	return messageID
}
