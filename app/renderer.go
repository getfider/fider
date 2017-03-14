package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"

	"os"

	"io/ioutil"

	"github.com/WeCanHearYou/wechy/identity"
	"github.com/WeCanHearYou/wechy/toolbox/env"
	"github.com/labstack/echo"
)

//HTMLRenderer renderer
type HTMLRenderer struct {
	templates map[string]*template.Template
	logger    echo.Logger
}

var path string

// NewHTMLRenderer creates a new HTMLRenderer
func NewHTMLRenderer(logger echo.Logger) *HTMLRenderer {
	renderer := &HTMLRenderer{nil, logger}
	renderer.templates = make(map[string]*template.Template)

	path = "views/"
	if env.IsTest() {
		path = os.Getenv("GOPATH") + "/src/github.com/WeCanHearYou/wechy/" + path
	}

	//TODO: load all templates automatically
	renderer.add("index.html")
	renderer.add("404.html")
	renderer.add("500.html")

	return renderer
}

//Render a template based on parameters
func (r *HTMLRenderer) add(name string) *template.Template {
	tpl, err := template.ParseFiles(path+"base.html", path+name)
	if err != nil {
		r.logger.Error(err)
	}

	r.templates[name] = tpl
	return tpl
}

//Render a template based on parameters
func (r *HTMLRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := r.templates[name]
	if !ok {
		panic(fmt.Errorf("The template '%s' does not exist", name))
	}

	if env.IsDevelopment() {
		tmpl = r.add(name)
	}

	protocol := "http://"
	if c.Request().TLS != nil {
		protocol = "https://"
	}

	//TODO: refactor (and move somewhere else?)
	m := data.(echo.Map)
	claims, ok := c.Get("Claims").(*identity.WechyClaims)

	m["AuthEndpoint"] = os.Getenv("AUTH_ENDPOINT")
	if ok {
		m["Claims"] = claims
		m["Gravatar"] = toMDF5(claims.UserEmail)
	}
	m["CurrentUrl"] = protocol + c.Request().Host + c.Request().URL.String()

	files, _ := ioutil.ReadDir("dist/js")
	if len(files) > 0 {
		m["JavaScriptBundle"] = files[0].Name()
	}

	return tmpl.Execute(w, m)
}

func toMDF5(content string) string {
	hasher := md5.New()
	hasher.Write([]byte(content))
	return fmt.Sprintf("%v", hex.EncodeToString(hasher.Sum(nil)))
}
