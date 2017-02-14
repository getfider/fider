package router

import (
	"fmt"
	"html/template"
	"io"
	"log"

	"os"

	"github.com/WeCanHearYou/wchy/env"
	"github.com/labstack/echo"
)

//HTMLRenderer renderer
type HTMLRenderer struct {
	templates map[string]*template.Template
}

// NewHTMLRenderer creates a new HTMLRenderer
func NewHTMLRenderer() *HTMLRenderer {
	renderer := &HTMLRenderer{}
	renderer.templates = make(map[string]*template.Template)

	path := "views/"
	if env.IsTest() {
		path = os.Getenv("GOPATH") + "/src/github.com/WeCanHearYou/wchy/" + path
	}

	renderer.add(path, "index.html")

	return renderer
}

//Render a template based on parameters
func (r *HTMLRenderer) add(path, name string) {
	tpl, err := template.ParseFiles(path+"base.html", path+name)
	if err != nil {
		log.Fatal(err)
	}

	r.templates[name] = tpl
}

//Render a template based on parameters
func (r *HTMLRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist", name)
	}

	return tmpl.Execute(w, data)
}
