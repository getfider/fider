package tpl

import (
	"html"
	"strings"
	"text/template"

	"github.com/getfider/fider/app/pkg/markdown"
)

func GetTextTemplate(name string, rawText string) (*template.Template, error) {
	tpl, err := template.New(name).Funcs(templateFunctions).Funcs(template.FuncMap{
		"markdown": func(input string) string {
			// Kept HTML-escaped: webhook destinations such as Slack read "<...>" as
			// control sequences, and HTML destinations must not receive raw tags.
			return html.EscapeString(markdown.PlainText(input))
		},
	}).Parse(rawText)
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func Execute(tmpl *template.Template, data any) (string, error) {
	builder := &strings.Builder{}
	if err := tmpl.Execute(builder, data); err != nil {
		return "", err
	}
	return builder.String(), nil
}
