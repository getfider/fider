package handlers

import (
	"github.com/WeCanHearYou/wchy-api/context"
	"github.com/gin-gonic/gin"
)

type tenantHandlers struct {
	ctx context.WchyContext
}

func TenantByDomain(ctx context.WchyContext) gin.HandlerFunc {
	return tenantHandlers{ctx: ctx}.byDomain()
}

func (h tenantHandlers) byDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{})
	}
}
