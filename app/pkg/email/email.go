package email

import (
	"bytes"
	"html/template"

	"github.com/getfider/fider/app/pkg/env"
)

var cache = make(map[string]*template.Template, 0)

// RenderMessage returns the HTML of an e-mail based on template and params
func RenderMessage(templateName string, params map[string]interface{}) string {
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
	return bf.String()
}

// NoReply is the default 'from' address
var NoReply = env.MustGet("NOREPLY_EMAIL")

// Sender is used to send e-mails
type Sender interface {
	Send(from, to, subject, templateName string, params map[string]interface{}) error
}

//NoopSender does not send e-mails
type NoopSender struct {
}

//NewNoopSender creates a new NoopSender
func NewNoopSender() *NoopSender {
	return &NoopSender{}
}

//Send an e-mail
func (s *NoopSender) Send(from, to, subject, templateName string, params map[string]interface{}) error {
	return nil
}
