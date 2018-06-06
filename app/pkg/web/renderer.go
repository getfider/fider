package web

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"io/ioutil"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/md5"
	"github.com/getfider/fider/app/pkg/oauth"
)

var templateFunctions = template.FuncMap{
	"md5": func(input string) string {
		return md5.Hash(input)
	},
	"markdown": func(input string) template.HTML {
		return markdown.Parse(input)
	},
}

//Renderer is the default HTML Render
type Renderer struct {
	templates    map[string]*template.Template
	logger       log.Logger
	settings     *models.SystemSettings
	jsBundle     string
	vendorBundle string
	cssBundle    string
}

// NewRenderer creates a new Renderer
func NewRenderer(settings *models.SystemSettings, logger log.Logger) *Renderer {
	r := &Renderer{
		templates: make(map[string]*template.Template),
		logger:    logger,
		settings:  settings,
	}

	r.add("index.html")
	r.add("not-invited.html")
	r.add("legal.html")
	r.add("403.html")
	r.add("404.html")
	r.add("410.html")
	r.add("500.html")

	r.jsBundle = r.getBundle("/dist/js", "main")
	r.vendorBundle = r.getBundle("/dist/js", "vendor")
	r.cssBundle = r.getBundle("/dist/css", "main")

	return r
}

//Render a template based on parameters
func (r *Renderer) add(name string) *template.Template {
	base := env.Path("/views/base.html")
	file := env.Path("/views", name)
	tpl, err := template.New("base.html").Funcs(templateFunctions).ParseFiles(base, file)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse template %s", file))
	}

	r.templates[name] = tpl
	return tpl
}

func (r *Renderer) getBundle(folder, prefix string) string {
	files, _ := ioutil.ReadDir(env.Path(folder))
	if len(files) > 0 {
		for _, file := range files {
			fileName := file.Name()
			if strings.HasPrefix(fileName, prefix) {
				return fileName
			}
		}
	}

	// Panic if bundle is not available in production mode
	if env.IsProduction() {
		panic(fmt.Sprintf("Bundle not found: %s/%s.", folder, prefix))
	}

	return ""
}

//Render a template based on parameters
func (r *Renderer) Render(w io.Writer, name string, props Props, ctx *Context) {
	tmpl, ok := r.templates[name]
	if !ok {
		panic(fmt.Errorf("The template '%s' does not exist", name))
	}

	if env.IsDevelopment() {
		tmpl = r.add(name)
		r.jsBundle = r.getBundle("/dist/js", "main")
		r.vendorBundle = r.getBundle("/dist/js", "vendor")
		r.cssBundle = r.getBundle("/dist/css", "main")
	}

	m := props.Data
	if m == nil {
		m = make(Map, 0)
	}

	tenantName := "Fider"
	if ctx.Tenant() != nil {
		tenantName = ctx.Tenant().Name
	}

	title := tenantName
	if props.Title != "" {
		title = fmt.Sprintf("%s · %s", props.Title, tenantName)
	}

	m["__Title"] = title

	if props.Description != "" {
		description := strings.Replace(props.Description, "\n", " ", -1)
		m["__Description"] = fmt.Sprintf("%.150s", description)
	}

	m["__VendorBundle"] = ctx.GlobalAssetsURL("/assets/js/%s", r.vendorBundle)
	m["__JavaScriptBundle"] = ctx.GlobalAssetsURL("/assets/js/%s", r.jsBundle)
	m["__StyleBundle"] = ctx.GlobalAssetsURL("/assets/css/%s", r.cssBundle)
	m["__ContextID"] = ctx.ContextID()
	if ctx.Tenant() != nil && ctx.Tenant().LogoID > 0 {
		m["__logo"] = ctx.TenantAssetsURL("/logo/200/%d", ctx.Tenant().LogoID)
		m["__favicon"] = ctx.TenantAssetsURL("/logo/50/%d", ctx.Tenant().LogoID)
	} else {
		m["__logo"] = "https://getfider.com/images/logo-100x100.png"
		m["__favicon"] = "/favicon.ico"
	}

	m["system"] = r.settings
	m["baseURL"] = ctx.BaseURL()
	m["currentURL"] = ctx.Request.URL.String()
	m["tenant"] = ctx.Tenant()
	m["auth"] = Map{
		"endpoint": ctx.AuthEndpoint(),
		"providers": Map{
			oauth.GoogleProvider:   oauth.IsProviderEnabled(oauth.GoogleProvider),
			oauth.FacebookProvider: oauth.IsProviderEnabled(oauth.FacebookProvider),
			oauth.GitHubProvider:   oauth.IsProviderEnabled(oauth.GitHubProvider),
		},
	}

	if ctx.IsAuthenticated() {
		u := ctx.User()
		m["user"] = &Map{
			"id":              u.ID,
			"name":            u.Name,
			"email":           u.Email,
			"role":            u.Role,
			"status":          u.Status,
			"isAdministrator": u.IsAdministrator(),
			"isCollaborator":  u.IsCollaborator(),
		}
	}

	err := tmpl.Execute(w, m)
	if err != nil {
		panic(errors.Wrap(err, "failed to execute template %s", name))
	}
}
