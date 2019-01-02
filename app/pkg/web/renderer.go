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
	templates     map[string]*template.Template
	logger        log.Logger
	settings      *models.SystemSettings
	assets        *clientAssets
	chunkedAssets map[string]*clientAssets
	mutex         sync.RWMutex
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
		ChunkGroups map[string]struct {
			Assets []string `json:"assets"`
		} `json:"namedChunkGroups"`
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

func getClientAssets(assets []string) *clientAssets {
	clientAssets := &clientAssets{
		CSS: make([]string, 0),
		JS:  make([]string, 0),
	}

	for _, asset := range assets {
		if strings.HasSuffix(asset, ".map") {
			continue
		}

		assetURL := "/assets/" + asset
		if strings.HasSuffix(asset, ".css") {
			clientAssets.CSS = append(clientAssets.CSS, assetURL)
		} else if strings.HasSuffix(asset, ".js") {
			clientAssets.JS = append(clientAssets.JS, assetURL)
		}
	}

	return clientAssets
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

	public := make(Map, 0)
	private := make(Map, 0)

	tenantName := "Fider"
	if ctx.Tenant() != nil {
		tenantName = ctx.Tenant().Name
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
	private["logo"] = ctx.LogoURL()
	private["favicon"] = ctx.FaviconURL()
	private["currentURL"] = ctx.Request.URL.String()
	if canonicalURL := ctx.Get("Canonical-URL"); canonicalURL != nil {
		private["canonicalURL"] = canonicalURL
	}

	oauthProviders := make([]*oauth.ProviderOption, 0)
	if !ctx.IsAuthenticated() && ctx.Services() != nil {
		oauthProviders, err = ctx.Services().OAuth.ListActiveProviders()
		if err != nil {
			panic(errors.Wrap(err, "failed to get list of providers"))
		}
	}

	public["contextID"] = ctx.ContextID()
	public["tenant"] = ctx.Tenant()
	public["props"] = props.Data
	public["settings"] = &Map{
		"mode":            r.settings.Mode,
		"buildTime":       r.settings.BuildTime,
		"version":         r.settings.Version,
		"environment":     r.settings.Environment,
		"compiler":        r.settings.Compiler,
		"googleAnalytics": r.settings.GoogleAnalytics,
		"stripePublicKey": env.Config.Stripe.PublicKey,
		"domain":          r.settings.Domain,
		"hasLegal":        r.settings.HasLegal,
		"baseURL":         ctx.BaseURL(),
		"tenantAssetsURL": ctx.TenantAssetsURL(""),
		"globalAssetsURL": ctx.GlobalAssetsURL(""),
		"oauth":           oauthProviders,
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

	err = tmpl.Execute(w, Map{
		"public":  public,
		"private": private,
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to execute template %s", name))
	}
}
