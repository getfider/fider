package handlers

import (
	"github.com/WeCanHearYou/wchy-api/context"
	"github.com/gin-gonic/gin"
)

type tenantHandlers struct {
	ctx context.WchyContext
}

// TenantByDomain creates a new TenantByDomain HTTP handler
func TenantByDomain(ctx context.WchyContext) gin.HandlerFunc {
	return tenantHandlers{ctx: ctx}.byDomain()
}

func (h tenantHandlers) byDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant := h.ctx.Tenant.GetByDomain(c.Param("domain"))
		if tenant == nil {
			c.Status(404)
		} else {
			c.JSON(200, tenant)
		}
	}
}
