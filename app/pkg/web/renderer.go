package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/getfider/fider/app/models/dto"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"

	"io/ioutil"

	"github.com/getfider/fider/app/pkg/crypto"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/markdown"
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

type distAsset struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type assetsFile struct {
	Entrypoints struct {
		Main struct {
			Assets []distAsset `json:"assets"`
		} `json:"main"`
	} `json:"entrypoints"`
	ChunkGroups map[string]struct {
		Assets []distAsset `json:"assets"`
	} `json:"namedChunkGroups"`
}

//Renderer is the default HTML Render
type Renderer struct {
	templates     map[string]*template.Template
	assets        *clientAssets
	chunkedAssets map[string]*clientAssets
	mutex         sync.RWMutex
	reactRenderer *ReactRenderer
}

// NewRenderer creates a new Renderer
func NewRenderer() *Renderer {
	return &Renderer{
		templates:     make(map[string]*template.Template),
		mutex:         sync.RWMutex{},
		reactRenderer: NewReactRenderer("ssr.js"),
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

	r.assets = getClientAssets(file.Entrypoints.Main.Assets)
	r.chunkedAssets = make(map[string]*clientAssets)

	for chunkName, chunkGroup := range file.ChunkGroups {
		r.chunkedAssets[chunkName] = getClientAssets(chunkGroup.Assets)
	}

	return nil
}

func getClientAssets(assets []distAsset) *clientAssets {
	clientAssets := &clientAssets{
		CSS: make([]string, 0),
		JS:  make([]string, 0),
	}

	for _, asset := range assets {
		if strings.HasSuffix(asset.Name, ".map") {
			continue
		}

		assetURL := "/assets/" + asset.Name
		if strings.HasSuffix(asset.Name, ".css") {
			clientAssets.CSS = append(clientAssets.CSS, assetURL)
		} else if strings.HasSuffix(asset.Name, ".js") {
			clientAssets.JS = append(clientAssets.JS, assetURL)
		}
	}

	return clientAssets
}

//Render a template based on parameters
func (r *Renderer) Render(w io.Writer, statusCode int, templateName string, props Props, ctx *Context) {
	var err error

	if r.assets == nil || env.IsDevelopment() {
		if err := r.loadAssets(); err != nil {
			panic(err)
		}
	}

	public := make(Map)
	private := make(Map)
	if props.Data == nil {
		props.Data = make(Map)
	}

	tenant := ctx.Tenant()
	tenantName := "Fider"
	if tenant != nil {
		tenantName = tenant.Name
	}

	title := tenantName
	if props.Title != "" {
		title = fmt.Sprintf("%s · %s", props.Title, tenantName)
	}

	public["title"] = title

	if props.Description != "" {
		description := strings.Replace(props.Description, "\n", " ", -1)
		public["description"] = fmt.Sprintf("%.150s", description)
	}

	if props.ChunkName != "" {
		private["chunkAssets"] = r.chunkedAssets[props.ChunkName]
	}

	private["assets"] = r.assets
	private["logo"] = LogoURL(ctx)

	if tenant == nil || tenant.LogoBlobKey == "" {
		private["favicon"] = AssetsURL(ctx, "/static/favicon")
	} else {
		private["favicon"] = AssetsURL(ctx, "/static/favicon/%s", tenant.LogoBlobKey)
	}

	private["currentURL"] = ctx.Request.URL.String()
	if canonicalURL := ctx.Value("Canonical-URL"); canonicalURL != nil {
		private["canonicalURL"] = canonicalURL
	}

	oauthProviders := &query.ListActiveOAuthProviders{
		Result: make([]*dto.OAuthProviderOption, 0),
	}
	if !ctx.IsAuthenticated() && statusCode >= 200 && statusCode < 300 {
		err = bus.Dispatch(ctx, oauthProviders)
		if err != nil {
			panic(errors.Wrap(err, "failed to get list of providers"))
		}
	}

	public["contextID"] = ctx.ContextID()
	public["sessionID"] = ctx.SessionID()
	public["tenant"] = tenant
	public["props"] = props.Data
	public["settings"] = &Map{
		"mode":            env.Config.HostMode,
		"environment":     env.Config.Environment,
		"googleAnalytics": env.Config.GoogleAnalytics,
		"domain":          env.MultiTenantDomain(),
		"hasLegal":        env.HasLegal(),
		"baseURL":         ctx.BaseURL(),
		"assetsURL":       AssetsURL(ctx, ""),
		"oauth":           oauthProviders.Result,
	}

	if ctx.IsAuthenticated() {
		u := ctx.User()
		public["user"] = &Map{
			"id":              u.ID,
			"name":            u.Name,
			"email":           u.Email,
			"role":            u.Role,
			"status":          u.Status,
			"avatarType":      u.AvatarType,
			"avatarURL":       u.AvatarURL,
			"avatarBlobKey":   u.AvatarBlobKey,
			"isAdministrator": u.IsAdministrator(),
			"isCollaborator":  u.IsCollaborator(),
		}
	}

	// Only index.html template uses React, other templates are already SSR
	if env.Config.Experimental_SSR_SEO && ctx.Request.IsCrawler() && templateName == "index.html" {
		html, err := r.reactRenderer.Render(ctx.Request.URL, public)
		if err != nil {
			log.Errorf(ctx, "Failed to render react page: @{Error}", dto.Props{
				"Error": err.Error(),
			})
		}
		if html != "" {
			templateName = "ssr.html"
			props.Data["html"] = template.HTML(html)
		}
	}

	tmpl, ok := r.templates[templateName]
	if !ok || env.IsDevelopment() {
		tmpl = r.add(templateName)
	}

	err = tmpl.Execute(w, Map{
		"public":  public,
		"private": private,
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to execute template %s", templateName))
	}
}
