package handler

import (
	"net/http"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/service"
	"github.com/labstack/echo"
)

type tenantHandlers struct {
	ctx *context.WchyContext
}

// TenantByDomain creates a new TenantByDomain HTTP handler
func TenantByDomain(ctx *context.WchyContext) echo.HandlerFunc {
	return tenantHandlers{ctx: ctx}.byDomain()
}

func (h tenantHandlers) byDomain() echo.HandlerFunc {
	return func(c echo.Context) error {
		tenant, err := h.ctx.Tenant.GetByDomain(c.Param("domain"))

		if err == service.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, tenant)
	}
}
