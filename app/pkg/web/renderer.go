package web

import (
	"fmt"
	"html/template"
	"io"

	"io/ioutil"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/labstack/echo"
)

//HTMLRenderer renderer
type HTMLRenderer struct {
	templates map[string]*template.Template
	logger    echo.Logger
	settings  *models.AppSettings
}

// NewHTMLRenderer creates a new HTMLRenderer
func NewHTMLRenderer(settings *models.AppSettings, logger echo.Logger) *HTMLRenderer {
	renderer := &HTMLRenderer{nil, logger, settings}
	renderer.templates = make(map[string]*template.Template)

	renderer.add("index.html")
	renderer.add("404.html")
	renderer.add("500.html")

	return renderer
}

//Render a template based on parameters
func (r *HTMLRenderer) add(name string) *template.Template {
	base := env.Path("/views/base.html")
	file := env.Path("/views", name)
	tpl, err := template.ParseFiles(base, file)
	if err != nil {
		panic(err)
	}

	r.templates[name] = tpl
	return tpl
}

//Render a template based on parameters
func (r *HTMLRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	ctx := Context{Context: c}

	tmpl, ok := r.templates[name]
	if !ok {
		panic(fmt.Errorf("The template '%s' does not exist", name))
	}

	if env.IsDevelopment() {
		tmpl = r.add(name)
	}

	m := data.(Map)
	m["baseUrl"] = ctx.BaseURL()
	m["tenant"] = ctx.Tenant()
	m["auth"] = Map{
		"endpoint": ctx.AuthEndpoint(),
		"providers": Map{
			oauth.GoogleProvider:   oauth.IsProviderEnabled(oauth.GoogleProvider),
			oauth.FacebookProvider: oauth.IsProviderEnabled(oauth.FacebookProvider),
		},
	}
	m["settings"] = r.settings

	if renderVars := ctx.RenderVars(); renderVars != nil {
		for key, value := range renderVars {
			m[key] = value
		}
	}

	if ctx.IsAuthenticated() {
		m["__User"] = ctx.User()
	}

	files, _ := ioutil.ReadDir("dist/js")
	if len(files) > 0 {
		m["__JavaScriptBundle"] = files[0].Name()
	}

	files, _ = ioutil.ReadDir("dist/css")
	if len(files) > 0 {
		m["__StyleBundle"] = files[0].Name()
	}

	return tmpl.Execute(w, m)
}
