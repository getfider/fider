package router

import (
	"html/template"

	"os"

	"fmt"

	"net/http"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/env"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/labstack/echo"
)

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *context.WchyContext) *echo.Echo {
	router := echo.New()

	path := "views/*.html"
	if env.IsTest() {
		path = os.Getenv("GOPATH") + "/src/github.com/WeCanHearYou/wchy/" + path
	}

	router.Renderer = &HTMLRenderer{
		Templates: template.Must(template.ParseGlob(path)),
	}

	router.HTTPErrorHandler = func(e error, c echo.Context) {
		fmt.Println(e)
		c.NoContent(http.StatusInternalServerError)
	}
	router.Use(MultiTenant(ctx))
	router.GET("/", handler.Index(ctx))

	api := router.Group("/api")
	{
		api.GET("/status", handler.Status(ctx))
		api.GET("/tenants/:domain", handler.TenantByDomain(ctx))
	}

	router.Static("/assets", "node_modules")

	return router
}
