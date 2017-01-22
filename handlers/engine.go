package handlers

import (
	"github.com/WeCanHearYou/wchy-api/context"
	_ "github.com/lib/pq" //

	"github.com/gin-gonic/gin"
)

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx context.WchyContext) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/status", Status(ctx))
	return router
}
