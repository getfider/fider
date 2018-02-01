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
	jsBundle  string
	cssBundle string
}

// NewRenderer creates a new Renderer
func NewRenderer(settings *models.AppSettings, logger log.Logger) *Renderer {
	r := &Renderer{
		templates: make(map[string]*template.Template),
		logger:    logger,
		settings:  settings,
	}

	r.add("index.html")
	r.add("403.html")
	r.add("404.html")
	r.add("500.html")

	r.jsBundle = r.getBundle("/dist/js")
	r.cssBundle = r.getBundle("/dist/css")

	return r
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

func (r *Renderer) getBundle(folder string) string {
	files, _ := ioutil.ReadDir(env.Path(folder))
	if len(files) > 0 {
		return files[0].Name()
	}
	panic(fmt.Sprintf("Bundle not found: %s.", folder))
}

//Render a template based on parameters
func (r *Renderer) Render(w io.Writer, name string, data interface{}, ctx *Context) error {
	tmpl, ok := r.templates[name]
	if !ok {
		panic(fmt.Errorf("The template '%s' does not exist", name))
	}

	if env.IsDevelopment() {
		tmpl = r.add(name)
		r.jsBundle = r.getBundle("/dist/js")
		r.cssBundle = r.getBundle("/dist/css")
	}

	m := data.(Map)
	m["__JavaScriptBundle"] = r.jsBundle
	m["__StyleBundle"] = r.cssBundle
	m["settings"] = r.settings

	m["baseURL"] = ctx.BaseURL()
	m["tenant"] = ctx.Tenant()
	m["auth"] = Map{
		"endpoint": ctx.AuthEndpoint(),
		"providers": Map{
			oauth.GoogleProvider:   oauth.IsProviderEnabled(oauth.GoogleProvider),
			oauth.FacebookProvider: oauth.IsProviderEnabled(oauth.FacebookProvider),
			oauth.GitHubProvider:   oauth.IsProviderEnabled(oauth.GitHubProvider),
		},
	}

	if renderVars := ctx.RenderVars(); renderVars != nil {
		for key, value := range renderVars {
			m[key] = value
		}
	}

	if ctx.IsAuthenticated() {
		u := ctx.User()
		m["user"] = &Map{
			"id":              u.ID,
			"name":            u.Name,
			"email":           u.Email,
			"role":            u.Role,
			"isAdministrator": u.IsAdministrator(),
			"isCollaborator":  u.IsCollaborator(),
		}
	}

	return tmpl.Execute(w, m)
}
