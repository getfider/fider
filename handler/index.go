package handler

import (
	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/model"
	"github.com/gin-gonic/gin"
)

type indexHandlder struct {
	ctx *context.WchyContext
}

// Index handles initial page
func Index(ctx *context.WchyContext) gin.HandlerFunc {
	return indexHandlder{ctx: ctx}.get()
}

func (h indexHandlder) get() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant := c.MustGet("Tenant").(*model.Tenant)
		ideas, _ := h.ctx.Idea.GetAll(tenant.ID)
		c.HTML(200, "index.html", gin.H{
			"Tenant": tenant,
			"Ideas":  ideas,
		})
	}
}
