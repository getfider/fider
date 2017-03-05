package router

import (
	"fmt"
	"html/template"
	"io"

	"os"

	"github.com/WeCanHearYou/wchy/app/auth"
	"github.com/WeCanHearYou/wchy/app/env"
	"github.com/WeCanHearYou/wchy/app/util"
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
		path = os.Getenv("GOPATH") + "/src/github.com/WeCanHearYou/wchy/" + path
	}

	renderer.add("index.html")

	return renderer
}

//Render a template based on parameters
func (r *HTMLRenderer) add(name string) {
	tpl, err := template.ParseFiles(path+"base.html", path+name)
	if err != nil {
		r.logger.Error(err)
	}

	r.templates[name] = tpl
}

//Render a template based on parameters
func (r *HTMLRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if env.IsDevelopment() {
		r.add(name)
	}

	tmpl, ok := r.templates[name]
	if !ok {
		panic(fmt.Errorf("The template %s does not exist", name))
	}

	protocol := "http://"
	if c.Request().TLS != nil {
		protocol = "https://"
	}

	//TODO: refactor (and move somewhere else?)
	m := data.(echo.Map)
	claims, ok := c.Get("Claims").(*auth.WchyClaims)

	m["AuthEndpoint"] = os.Getenv("AUTH_ENDPOINT")
	if ok {
		m["Claims"] = claims
		m["Gravatar"] = util.MD5(claims.UserEmail)
	}
	m["CurrentUrl"] = protocol + c.Request().Host + c.Request().URL.String()

	return tmpl.Execute(w, m)
}
