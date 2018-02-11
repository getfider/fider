package email

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
)

var cache = make(map[string]*template.Template, 0)

// Message represents what is sent by e-mail
type Message struct {
	Subject string
	Body    string
}

// RenderMessage returns the HTML of an e-mail based on template and params
func RenderMessage(templateName string, params map[string]interface{}) *Message {
	tpl, ok := cache[templateName]
	if !ok {
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

// Recipient contains details of who is receiving the e-mail
type Recipient struct {
	Name    string
	Address string
	Params  map[string]interface{}
}

// NewRecipient creates a new Recipient
func NewRecipient(address string, params map[string]interface{}) Recipient {
	return Recipient{
		Address: address,
		Params:  params,
	}
}

var whitelist = env.GetEnvOrDefault("EMAIL_WHITELIST", "")
var whitelistRegex = regexp.MustCompile(whitelist)

// SetWhitelist can be used to change e-mail whitelist during rutime
func SetWhitelist(s string) {
	whitelist = s
	whitelistRegex = regexp.MustCompile(whitelist)
}

// CanSendTo returns true if Fider is allowed to send e-mail to given address
func CanSendTo(address string) bool {
	if whitelist == "" {
		return true
	}
	return whitelistRegex.MatchString(address)
}

// Sender is used to send e-mails
type Sender interface {
	Send(templateName, from string, to Recipient) error
	BatchSend(templateName, from string, to []Recipient) error
}

// NoopSender does not send e-mails
type NoopSender struct {
}

// NewNoopSender creates a new NoopSender
func NewNoopSender() *NoopSender {
	return &NoopSender{}
}

// Send an e-mail
func (s *NoopSender) Send(templateName, from string, to Recipient) error {
	return nil

}

// BatchSend an e-mail to multiple recipients
func (s *NoopSender) BatchSend(templateName, from string, to []Recipient) error {
	return nil
}
