package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"sync"

	"io/ioutil"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/crypto"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/oauth"
)

var templateFunctions = template.FuncMap{
	"md5": func(input string) string {
		return crypto.MD5(input)
	},
	"markdown": func(input string) template.HTML {
		return markdown.Full(input)
	},
}

type clientAssets struct {
	CSS []string
	JS  []string
}

//Renderer is the default HTML Render
type Renderer struct {
	templates map[string]*template.Template
	logger    log.Logger
	settings  *models.SystemSettings
	assets    *clientAssets
	mutex     sync.RWMutex
}

// NewRenderer creates a new Renderer
func NewRenderer(settings *models.SystemSettings, logger log.Logger) *Renderer {
	return &Renderer{
		templates: make(map[string]*template.Template),
		logger:    logger,
		settings:  settings,
		mutex:     sync.RWMutex{},
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

func (r *Renderer) loadAssets() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.assets != nil && env.IsProduction() {
		return nil
	}

	type assetsFile struct {
		Entrypoints struct {
			Main struct {
				Assets []string `json:"assets"`
			} `json:"main"`
		} `json:"entrypoints"`
	}

	assetsFilePath := "/dist/assets.json"
	if env.IsTest() {
		// Load a fake assets.json for Unit Testing
		assetsFilePath = "/app/pkg/web/testdata/assets.json"
	}

	jsonFile, err := os.Open(env.Path(assetsFilePath))
	if err != nil {
		return errors.Wrap(err, "failed to open file: assets.json")
	}
	defer jsonFile.Close()

	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	file := &assetsFile{}
	err = json.Unmarshal([]byte(jsonBytes), file)
	if err != nil {
		return errors.Wrap(err, "failed to parse file: assets.json")
	}

	r.assets = &clientAssets{
		CSS: make([]string, 0),
		JS:  make([]string, 0),
	}

	for _, asset := range file.Entrypoints.Main.Assets {
		if strings.HasSuffix(asset, ".map") {
			continue
		}

		assetURL := "/assets/" + asset
		if strings.HasSuffix(asset, ".css") {
			r.assets.CSS = append(r.assets.CSS, assetURL)
		} else if strings.HasSuffix(asset, ".js") {
			r.assets.JS = append(r.assets.JS, assetURL)
		}
	}

	return nil
}

//Render a template based on parameters
func (r *Renderer) Render(w io.Writer, name string, props Props, ctx *Context) {
	var err error

	if r.assets == nil || env.IsDevelopment() {
		err := r.loadAssets()
		if err != nil && !env.IsTest() {
			panic(err)
		}
	}

	tmpl, ok := r.templates[name]
	if !ok || env.IsDevelopment() {
		tmpl = r.add(name)
	}

	m := make(Map, 0)

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

	m["__assets"] = r.assets
	m["__logo"] = ctx.LogoURL()
	m["__favicon"] = ctx.FaviconURL()
	m["__contextID"] = ctx.ContextID()
	m["__currentURL"] = ctx.Request.URL.String()
	m["__tenant"] = ctx.Tenant()
	m["__props"] = props.Data

	oauthProviders := make([]*oauth.ProviderOption, 0)
	if !ctx.IsAuthenticated() && ctx.Services() != nil {
		oauthProviders, err = ctx.Services().OAuth.ListActiveProviders()
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
		"tenantAssetsURL": ctx.TenantAssetsURL(""),
		"globalAssetsURL": ctx.GlobalAssetsURL(""),
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
