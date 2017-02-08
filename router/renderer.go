package router

import (
	"os"
	"path/filepath"

	"github.com/WeCanHearYou/wchy/env"
	"github.com/gin-contrib/multitemplate"
)

func CreateTemplateRender() multitemplate.Render {
	r := multitemplate.New()

	path := filepath.Join(os.Getenv("GOPATH"), "src/github.com/WeCanHearYou/wchy/views")
	if !env.IsTest() {
		path = "views"
	}

	r.AddFromFiles("index", path+"/base.html", path+"/index.html")
	return r
}
