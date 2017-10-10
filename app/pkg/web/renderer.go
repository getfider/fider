package web

import (
	"fmt"
	"html/template"
	"io"

	"io/ioutil"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/oauth"
)

//Renderer is the default HTML Render
type Renderer struct {
	templates map[string]*template.Template
	logger    log.Logger
	settings  *models.AppSettings
}

// NewRenderer creates a new Renderer
func NewRenderer(settings *models.AppSettings, logger log.Logger) *Renderer {
	renderer := &Renderer{nil, logger, settings}
	renderer.templates = make(map[string]*template.Template)

	renderer.add("index.html")
	renderer.add("403.html")
	renderer.add("404.html")
	renderer.add("500.html")

	return renderer
}

//Render a template based on parameters
func (r *Renderer) add(name string) *template.Template {
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
func (r *Renderer) Render(w io.Writer, name string, data interface{}, ctx *Context) error {
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
			oauth.GitHubProvider:   oauth.IsProviderEnabled(oauth.GitHubProvider),
		},
	}
	m["settings"] = r.settings

	if renderVars := ctx.RenderVars(); renderVars != nil {
		for key, value := range renderVars {
			m[key] = value
		}
	}

	if ctx.IsAuthenticated() {
		m["user"] = ctx.User()
		m["email"] = ctx.User().Email
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
