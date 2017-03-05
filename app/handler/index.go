package handler

import (
	"github.com/WeCanHearYou/wchy/app/context"
	"github.com/WeCanHearYou/wchy/app/model"
	"github.com/labstack/echo"
)

type indexHandlder struct {
	ctx *context.WchyContext
}

// Index handles initial page
func Index(ctx *context.WchyContext) echo.HandlerFunc {
	return indexHandlder{ctx: ctx}.get()
}

func (h indexHandlder) get() echo.HandlerFunc {
	return func(c echo.Context) error {
		tenant := c.Get("Tenant").(*model.Tenant)
		ideas, _ := h.ctx.Idea.GetAll(tenant.ID)
		return c.Render(200, "index.html", echo.Map{
			"Tenant": tenant,
			"Ideas":  ideas,
		})
	}
}
