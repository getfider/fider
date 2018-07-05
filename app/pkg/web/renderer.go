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
	templates map[string]*template.Template
	logger    log.Logger
	settings  *models.SystemSettings
	assets    map[string]string
}

// NewRenderer creates a new Renderer
func NewRenderer(settings *models.SystemSettings, logger log.Logger) *Renderer {
	return &Renderer{
		templates: make(map[string]*template.Template),
		logger:    logger,
		settings:  settings,
		assets:    make(map[string]string, 0),
	}
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

func (r *Renderer) getBundle(folder, prefix, suffix string) string {
	files, _ := ioutil.ReadDir(env.Path(folder))
	if len(files) > 0 {
		for _, file := range files {
			fileName := file.Name()
			if strings.HasPrefix(fileName, prefix) && strings.HasSuffix(fileName, suffix) {
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
	var err error

	if len(r.assets) == 0 || env.IsDevelopment() {
		r.assets["main.js"] = r.getBundle("/dist/js", "main", "js")
		r.assets["vendor.js"] = r.getBundle("/dist/js", "vendor", "js")
		r.assets["main.css"] = r.getBundle("/dist/css", "main", "css")
		r.assets["icons.woff2"] = r.getBundle("/dist/fonts", "icons", "woff2")
	}

	tmpl, ok := r.templates[name]
	if !ok || env.IsDevelopment() {
		tmpl = r.add(name)
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

	m["__title"] = title

	if props.Description != "" {
		description := strings.Replace(props.Description, "\n", " ", -1)
		m["__description"] = fmt.Sprintf("%.150s", description)
	}

	m["__vendorBundle"] = ctx.GlobalAssetsURL("/assets/js/%s", r.assets["vendor.js"])
	m["__jsBundle"] = ctx.GlobalAssetsURL("/assets/js/%s", r.assets["main.js"])
	m["__cssBundle"] = ctx.GlobalAssetsURL("/assets/css/%s", r.assets["main.css"])
	m["__fontBundle"] = ctx.GlobalAssetsURL("/assets/fonts/%s", r.assets["icons.woff2"])

	if ctx.Tenant() != nil && ctx.Tenant().LogoID > 0 {
		m["__logo"] = ctx.TenantAssetsURL("/logo/200/%d", ctx.Tenant().LogoID)
		m["__favicon"] = ctx.TenantAssetsURL("/logo/50/%d", ctx.Tenant().LogoID)
	} else {
		m["__logo"] = "https://getfider.com/images/logo-100x100.png"
		m["__favicon"] = ctx.GlobalAssetsURL("/favicon.ico")
	}

	m["__contextID"] = ctx.ContextID()
	m["__currentURL"] = ctx.Request.URL.String()
	m["__tenant"] = ctx.Tenant()

	oauthProviders := make([]*oauth.ProviderOption, 0)
	if !ctx.IsAuthenticated() && ctx.Services() != nil {
		oauthProviders, err = ctx.Services().OAuth.ListProviders()
		if err != nil {
			panic(errors.Wrap(err, "failed to get list of providers"))
		}
	}

	m["__settings"] = &Map{
		"mode":            r.settings.Mode,
		"buildTime":       r.settings.BuildTime,
		"version":         r.settings.Version,
		"environment":     r.settings.Environment,
		"compiler":        r.settings.Compiler,
		"googleAnalytics": r.settings.GoogleAnalytics,
		"domain":          r.settings.Domain,
		"hasLegal":        r.settings.HasLegal,
		"baseURL":         ctx.BaseURL(),
		"assetsURL":       ctx.TenantAssetsURL(""),
		"oauth":           oauthProviders,
	}

	if ctx.IsAuthenticated() {
		u := ctx.User()
		m["__user"] = &Map{
			"id":              u.ID,
			"name":            u.Name,
			"email":           u.Email,
			"role":            u.Role,
			"status":          u.Status,
			"isAdministrator": u.IsAdministrator(),
			"isCollaborator":  u.IsCollaborator(),
		}
	}

	err = tmpl.Execute(w, m)
	if err != nil {
		panic(errors.Wrap(err, "failed to execute template %s", name))
	}
}
