package smtp

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	gosmtp "net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/services/email"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "SMTP"
}

func (s Service) Category() string {
	return "email"
}

func (s Service) Enabled() bool {
	return env.Config.Email.Type == "smtp"
}

func (s Service) Init() {
	bus.AddListener(sendMail)
	bus.AddHandler(fetchRecentSupressions)
}

func fetchRecentSupressions(ctx context.Context, c *query.FetchRecentSupressions) error {
	//not implemented for SMTP
	return nil
}

func sendMail(ctx context.Context, c *cmd.SendMail) {
	if c.Props == nil {
		c.Props = dto.Props{}
	}

	if c.From.Address == "" {
		c.From.Address = email.NoReply
	}

	for _, to := range c.To {
		if to.Address == "" {
			return
		}

		u, err := url.Parse(web.BaseURL(ctx))
		localname := "localhost"
		if err == nil {
			localname = u.Hostname()
		}

		if !email.CanSendTo(to.Address) {
			log.Warnf(ctx, "Skipping email to '@{Name} <@{Address}>'.", dto.Props{
				"Name":    to.Name,
				"Address": to.Address,
			})
			return
		}

		log.Debugf(ctx, "Sending email to @{Address} with template @{TemplateName} and params @{Props}.", dto.Props{
			"Address":      to.Address,
			"TemplateName": c.TemplateName,
			"Props":        to.Props,
		})

		message := email.RenderMessage(ctx, c.TemplateName, c.From.Address, c.Props.Merge(to.Props))
		b := builder{}
		b.Set("From", c.From.String())
		b.Set("Reply-To", c.From.Address)
		b.Set("To", to.String())
		b.Set("Subject", email.EncodeSubject(message.Subject))
		b.Set("MIME-version", "1.0")
		b.Set("Content-Type", "text/html; charset=\"UTF-8\"")
		b.Set("Date", time.Now().Format(time.RFC1123Z))
		b.Set("Message-ID", generateMessageID(localname))
		b.Body(message.Body)

		smtpConfig := env.Config.Email.SMTP
		servername := fmt.Sprintf("%s:%s", smtpConfig.Host, smtpConfig.Port)
		auth, err := authenticate(ctx, smtpConfig)
		if err != nil {
			panic(errors.Wrap(err, "failed to build smtp auth"))
		}

		err = Send(localname, servername, smtpConfig.EnableStartTLS, auth, email.NoReply, []string{to.Address}, b.Bytes())
		if err != nil {
			panic(errors.Wrap(err, "failed to send email with template %s", c.TemplateName))
		}
		log.Debug(ctx, "Email sent.")
	}
}

var Send = func(localName, serverAddress string, enableStartTLS bool, a gosmtp.Auth, from string, to []string, msg []byte) error {
	host, _, _ := net.SplitHostPort(serverAddress)
	c, err := gosmtp.Dial(serverAddress)
	if err != nil {
		return err
	}
	defer func() { _ = c.Close() }()
	if err = c.Hello(localName); err != nil {
		return err
	}
	if enableStartTLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			config := &tls.Config{ServerName: host}
			if err = c.StartTLS(config); err != nil {
				return err
			}
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

func authenticate(ctx context.Context, cfg env.SMTPConfig) (gosmtp.Auth, error) {
	mech := strings.ToLower(strings.TrimSpace(cfg.AuthMechanism))
	switch mech {
	case "", "agnostic":
		if cfg.Username == "" && cfg.Password == "" {
			return nil, nil
		}
		return AgnosticAuth("", cfg.Username, cfg.Password, cfg.Host), nil

	case "xoauth2":
		if cfg.Username == "" {
			return nil, errors.New("smtp: username is required for XOAUTH2")
		}

		scopes := splitCommaScopes(cfg.Scopes)
		token, err := getClientCredentialsToken(ctx, cfg.TokenUrl, cfg.ClientId, cfg.ClientSecret, scopes)
		if err != nil {
			return nil, err
		}

		return XOAuth2Auth(cfg.Username, token, cfg.Host), nil

	default:
		return nil, errors.New("smtp: unsupported auth mechanism")
	}
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
