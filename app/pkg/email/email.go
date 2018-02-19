package email

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
)

var cache = make(map[string]*template.Template, 0)

// Params used to replace variables on emails
type Params map[string]interface{}

// Merge given params into current params
func (p Params) Merge(p2 Params) Params {
	for k, v := range p2 {
		p[k] = v
	}
	return p
}

// Message represents what is sent by email
type Message struct {
	Subject string
	Body    string
}

// RenderMessage returns the HTML of an email based on template and params
func RenderMessage(templateName string, params Params) *Message {
	tpl, ok := cache[templateName]
	if !ok || env.IsDevelopment() {
		var err error
		file := env.Path("/views/templates", templateName+".tpl")
		tpl, err = template.ParseFiles(file)
		if err != nil {
			panic(err)
		}
		cache[templateName] = tpl
	}

	var bf bytes.Buffer
	tpl.Execute(&bf, params)
	lines := strings.Split(bf.String(), "\n")
	return &Message{
		Subject: strings.TrimLeft(lines[0], "subject: "),
		Body:    strings.TrimLeft(strings.Join(lines[2:], "\n"), " "),
	}
}

// NoReply is the default 'from' address
var NoReply = env.MustGet("NOREPLY_EMAIL")

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

var whitelist = env.GetEnvOrDefault("EMAIL_WHITELIST", "")
var whitelistRegex = regexp.MustCompile(whitelist)

// SetWhitelist can be used to change email whitelist during rutime
func SetWhitelist(s string) {
	whitelist = s
	whitelistRegex = regexp.MustCompile(whitelist)
}

// CanSendTo returns true if Fider is allowed to send email to given address
func CanSendTo(address string) bool {
	if strings.TrimSpace(address) == "" {
		return false
	}
	if whitelist == "" {
		return true
	}
	return whitelistRegex.MatchString(address)
}

// Sender is used to send emails
type Sender interface {
	Send(templateName string, params Params, from string, to Recipient) error
	BatchSend(templateName string, params Params, from string, to []Recipient) error
}

// NoopSender does not send emails
type NoopSender struct {
}

// NewNoopSender creates a new NoopSender
func NewNoopSender() *NoopSender {
	return &NoopSender{}
}

// Send an email
func (s *NoopSender) Send(templateName string, params Params, from string, to Recipient) error {
	return nil

}

// BatchSend an email to multiple recipients
func (s *NoopSender) BatchSend(templateName string, params Params, from string, to []Recipient) error {
	return nil
}
