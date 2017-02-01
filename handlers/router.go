package handlers

import (
	"net/http"
	"strings"

	"github.com/WeCanHearYou/wchy/context"

	"github.com/gin-gonic/gin"
)

func multiTenant(ctx context.WchyContext) gin.HandlerFunc {
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
func GetMainEngine(ctx context.WchyContext) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(multiTenant(ctx))
	router.LoadHTMLGlob("views/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": c.MustGet("Tenant"),
		})
	})

	api := router.Group("/api")
	{
		api.GET("/status", Status(ctx))
		api.GET("/tenants/:domain", TenantByDomain(ctx))
	}
	return router
}
