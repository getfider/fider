package router

import (
	"fmt"

	"net/http"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/labstack/echo"
)

func errorHandler(e error, c echo.Context) {
	fmt.Println(e)
	c.NoContent(http.StatusInternalServerError)
}

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *context.WchyContext) *echo.Echo {
	router := echo.New()

	router.Renderer = NewHTMLRenderer()
	router.HTTPErrorHandler = errorHandler

	router.Static("/favicon.ico", "public/imgs/favicon.ico")
	router.Static("/public", "public")
	router.Static("/assets", "node_modules") //TODO: Don't expose node_modules

	router.GET("/oauth", handler.OAuth(ctx))

	app := router.Group("", MultiTenant(ctx))
	{
		app.GET("/", handler.Index(ctx))

		api := app.Group("/api")
		{
			api.GET("/status", handler.Status(ctx))
			api.GET("/tenants/:domain", handler.TenantByDomain(ctx))
		}
	}

	return router
}
