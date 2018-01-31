package email

import (
	"bytes"
	"html/template"
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
func RenderMessage(templateName string, params map[string]interface{}) Message {
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
	return Message{
		Subject: strings.TrimLeft(lines[0], "subject: "),
		Body:    strings.TrimLeft(strings.Join(lines[2:], "\n"), " "),
	}
}

// NoReply is the default 'from' address
var NoReply = env.MustGet("NOREPLY_EMAIL")

// Sender is used to send e-mails
type Sender interface {
	Send(from, to, templateName string, params map[string]interface{}) error
}

//NoopSender does not send e-mails
type NoopSender struct {
}

//NewNoopSender creates a new NoopSender
func NewNoopSender() *NoopSender {
	return &NoopSender{}
}

//Send an e-mail
func (s *NoopSender) Send(from, to, templateName string, params map[string]interface{}) error {
	return nil
}
