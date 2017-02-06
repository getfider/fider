package handler

import (
	"net/http"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/service"
	"github.com/gin-gonic/gin"
)

type tenantHandlers struct {
	ctx *context.WchyContext
}

// TenantByDomain creates a new TenantByDomain HTTP handler
func TenantByDomain(ctx *context.WchyContext) gin.HandlerFunc {
	return tenantHandlers{ctx: ctx}.byDomain()
}

func (h tenantHandlers) byDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, err := h.ctx.Tenant.GetByDomain(c.Param("domain"))
		if err == service.ErrNotFound {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, tenant)
		}
	}
}
