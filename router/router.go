package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"

	"github.com/WeCanHearYou/wchy/env"
	"github.com/gin-gonic/gin"
)

func multiTenant(ctx *context.WchyContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		hostname := stripPort(c.Request.Host)
		tenant, err := ctx.Tenant.GetByDomain(hostname)
		if err == nil {
			c.Set("Tenant", tenant.Name)
			c.Next()
		} else {
			c.Status(404)
		}
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *context.WchyContext) *gin.Engine {
	router := gin.New()

	if env.IsDevelopment() {
		router.Use(gin.Logger())
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(multiTenant(ctx))

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
