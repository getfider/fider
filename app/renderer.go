package app

import (
	"fmt"
	"html/template"
	"io"

	"os"

	"io/ioutil"

	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/toolbox/env"
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

	//TODO: refactor (and move somewhere else?)
	m := data.(echo.Map)
	m["AuthEndpoint"] = os.Getenv("AUTH_ENDPOINT")

	claims, ok := c.Get("Claims").(*identity.WechyClaims)
	if ok {
		m["User"] = &identity.User{
			ID:    claims.UserID,
			Name:  claims.UserName,
			Email: claims.UserEmail,
		}
	}

	files, _ := ioutil.ReadDir("dist/js")
	if len(files) > 0 {
		m["JavaScriptBundle"] = files[0].Name()
	}

	return tmpl.Execute(w, m)
}
