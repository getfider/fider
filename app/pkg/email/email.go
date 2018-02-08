package email

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
)

var cache = make(map[string]*template.Template, 0)

// Message represents what is sent by e-mail
type Message struct {
	To      string
	From    string
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

// Sender is used to send e-mails
type Sender interface {
	Send(templateName, from string, to Recipient) (*Message, error)
}

//NoopSender does not send e-mails
type NoopSender struct {
}

//NewNoopSender creates a new NoopSender
func NewNoopSender() *NoopSender {
	return &NoopSender{}
}

//Send an e-mail
func (s *NoopSender) Send(templateName, from string, to Recipient) (*Message, error) {
	msg := RenderMessage(templateName, to.Params)
	msg.To = to.Address
	msg.From = fmt.Sprintf("%s <%s>", from, NoReply)
	return msg, nil
}
