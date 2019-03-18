package smtp

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	gosmtp "net/smtp"
	"strconv"
	"time"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/services/email"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Category() string {
	return "email"
}

func (s Service) Enabled() bool {
	return env.Config.Email.Mailgun.APIKey == ""
}

func (s Service) Init() {
	bus.AddEventListener(sendEmail)
}

func sendEmail(ctx context.Context, cmd *email.SendMessageCommand) {
	if cmd.Params == nil {
		cmd.Params = email.Params{}
	}

	for _, to := range cmd.To {
		if to.Address == "" {
			return
		}

		// EXPERIMENTAL-BUS how do I get baseURL hERE?
		// u, err := url.Parse(ctx.BaseURL())
		localname := "localhost"
		// if err == nil {
		// 	localname = u.Hostname()
		// }

		if !email.CanSendTo(to.Address) {
			log.Warnf(ctx, "Skipping email to '@{Name} <@{Address}>'.", log.Props{
				"Name":    to.Name,
				"Address": to.Address,
			})
			return
		}

		log.Debugf(ctx, "Sending email to @{Address} with template @{TemplateName} and params @{Params}.", log.Props{
			"Address":      to.Address,
			"TemplateName": cmd.TemplateName,
			"Params":       to.Params,
		})

		message := email.RenderMessage(cmd.TemplateName, cmd.Params.Merge(to.Params))
		b := builder{}
		b.Set("From", email.NewRecipient(cmd.From, email.NoReply, email.Params{}).String())
		b.Set("Reply-To", email.NoReply)
		b.Set("To", to.String())
		b.Set("Subject", message.Subject)
		b.Set("MIME-version", "1.0")
		b.Set("Content-Type", "text/html; charset=\"UTF-8\"")
		b.Set("Date", time.Now().Format(time.RFC1123Z))
		b.Set("Message-ID", generateMessageID(localname))
		b.Body(message.Body)

		smtpConfig := env.Config.Email.SMTP
		servername := fmt.Sprintf("%s:%s", smtpConfig.Host, smtpConfig.Port)
		auth := authenticate(smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)
		err := Send(localname, servername, auth, email.NoReply, []string{to.Address}, b.Bytes())
		if err != nil {
			panic(errors.Wrap(err, "failed to send email with template %s", cmd.TemplateName))
		}
		log.Debug(ctx, "Email sent.")
	}
}

var Send = func(localName, serverAddress string, a gosmtp.Auth, from string, to []string, msg []byte) error {
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
