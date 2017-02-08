package router

import (
	"net/http"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"

	"github.com/WeCanHearYou/wchy/env"
	"github.com/gin-gonic/gin"
)

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *context.WchyContext) *gin.Engine {
	router := gin.New()
	router.HTMLRender = CreateTemplateRender()

	if env.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(MultiTenant(ctx))
	router.GET("/", handler.Index(ctx))

	api := router.Group("/api")
	{
		api.GET("/status", handler.Status(ctx))
		api.GET("/tenants/:domain", handler.TenantByDomain(ctx))
	}

	router.StaticFS("/assets", http.Dir("node_modules"))

	return router
}
