package tpl

import (
	"context"
	"html/template"
	"io"
	"path"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/i18n"
	"github.com/getfider/fider/app/pkg/markdown"
)

var cache = make(map[string]*template.Template)

// GetTemplate returns the parsed template. Each of the extra files becomes a template named after
// its file name, which the others can include with {{template "file.html" .}}.
func GetTemplate(baseFileName, templateFileName string, extraFileNames ...string) *template.Template {
	cacheKey := strings.Join(append([]string{templateFileName}, extraFileNames...), "|")
	tmpl, ok := cache[cacheKey]
	if ok && !env.IsDevelopment() {
		return tmpl
	}

	baseFile := env.Path(baseFileName)
	files := []string{baseFile, env.Path(templateFileName)}
	for _, extraFileName := range extraFileNames {
		files = append(files, env.Path(extraFileName))
	}

	tpl, err := template.New(path.Base(baseFile)).Funcs(templateFunctions).ParseFiles(files...)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse template %s", templateFileName))
	}

	cache[cacheKey] = tpl
	return tpl
}

func Render(ctx context.Context, tmpl *template.Template, w io.Writer, data any) error {
	if err := template.Must(tmpl.Clone()).Funcs(template.FuncMap{
		"translate": func(key string, params ...i18n.Params) string {
			return i18n.T(ctx, key, params...)
		},
		"markdown": func(input string) template.HTML {
			return markdown.Full(ctx, input, true)
		},
	}).Execute(w, data); err != nil {
		return err
	}
	return nil
}
