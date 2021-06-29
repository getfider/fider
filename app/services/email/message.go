package email

import (
	"bytes"
	"context"
	"strings"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/tpl"
)

// Message represents what is sent by email
type Message struct {
	Subject string
	Body    string
}

// RenderMessage returns the HTML of an email based on template and params
func RenderMessage(ctx context.Context, templateName string, params dto.Props) *Message {
	tmpl := tpl.GetTemplate("/views/email/base_email.html", "/views/email/"+templateName+".html")
	var bf bytes.Buffer
	if err := tpl.Render(ctx, tmpl, &bf, params.Merge(dto.Props{
		"logo": params["logo"],
	})); err != nil {
		panic(err)
	}

	lines := strings.Split(bf.String(), "\n")
	body := strings.TrimLeft(strings.Join(lines[2:], "\n"), " ")

	return &Message{
		Subject: strings.TrimLeft(lines[0], "subject: "),
		Body:    body,
	}
}
