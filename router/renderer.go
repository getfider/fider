package router

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

//HTMLRenderer renderer
type HTMLRenderer struct {
	Templates *template.Template
}

//Render a template based on parameters
func (t *HTMLRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
