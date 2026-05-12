package email

import (
	"bytes"
	"context"
	"html"
	"mime"
	"strings"
	"unicode"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/tpl"
)

// Message represents what is sent by email
type Message struct {
	Subject string
	Body    string
}

// EncodeSubject encodes the subject header according to RFC 2047 if it contains non-ASCII characters
func EncodeSubject(subject string) string {
	// Check if the subject contains non-ASCII characters
	hasNonASCII := false
	for _, r := range subject {
		if r > unicode.MaxASCII {
			hasNonASCII = true
			break
		}
	}
	
	// If it's pure ASCII, no encoding is needed
	if !hasNonASCII {
		return subject
	}
	
	// Use Go's built-in MIME header encoding which implements RFC 2047
	return mime.QEncoding.Encode("utf-8", subject)
}

// RenderMessage returns the HTML of an email based on template and params
func RenderMessage(ctx context.Context, templateName string, fromAddress string, params dto.Props) *Message {
	noreply := fromAddress == NoReply

	tmpl := tpl.GetTemplate("/views/email/base_email.html", "/views/email/"+templateName+".html")
	var bf bytes.Buffer
	if err := tpl.Render(ctx, tmpl, &bf, params.Merge(dto.Props{
		"logo":    params["logo"],
		"noreply": noreply,
	})); err != nil {
		panic(err)
	}

	lines := strings.Split(bf.String(), "\n")
	body := strings.TrimLeft(strings.Join(lines[2:], "\n"), " ")

	// The subject is rendered through html/template (because it shares the
	// template set with the HTML body), which escapes characters like '+',
	// '"', '&', '<', '>' to numeric entities (e.g. "&#43;"). SMTP headers
	// are plain text and clients do not decode HTML entities in them, so
	// undo the escaping for the subject only.
	subject := html.UnescapeString(strings.TrimLeft(lines[0], "subject: "))

	return &Message{
		Subject: subject,
		Body:    body,
	}
}
