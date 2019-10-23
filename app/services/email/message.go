package email

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/env"
)

var cache = make(map[string]*template.Template)

// Message represents what is sent by email
type Message struct {
	Subject string
	Body    string
}

var baseTpl, _ = template.ParseFiles(env.Path("/views/templates/base_email.tpl"))

// RenderMessage returns the HTML of an email based on template and params
func RenderMessage(templateName string, params dto.Props) *Message {
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
	if err := tpl.Execute(&bf, params); err != nil {
		panic(err)
	}

	lines := strings.Split(bf.String(), "\n")
	body := strings.TrimLeft(strings.Join(lines[2:], "\n"), " ")

	bf.Reset()
	if err := baseTpl.Execute(&bf, dto.Props{
		"logo": params["logo"],
		"body": template.HTML(body),
	}); err != nil {
		panic(err)
	}

	return &Message{
		Subject: strings.TrimLeft(lines[0], "subject: "),
		Body:    bf.String(),
	}
}
