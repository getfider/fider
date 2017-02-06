package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"

	"github.com/WeCanHearYou/wchy/env"
	"github.com/gin-gonic/gin"
)

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *context.WchyContext) *gin.Engine {
	router := gin.New()

	if env.IsDevelopment() {
		router.Use(gin.Logger())
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(MultiTenant(ctx))

	if env.IsTest() {
		router.LoadHTMLGlob(filepath.Join(os.Getenv("GOPATH"), "src/github.com/WeCanHearYou/wchy/views/*"))
	} else {
		router.LoadHTMLGlob("views/*")
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": c.MustGet("Tenant"),
		})
	})

	api := router.Group("/api")
	{
		api.GET("/status", handler.Status(ctx))
		api.GET("/tenants/:domain", handler.TenantByDomain(ctx))
	}
	return router
}
