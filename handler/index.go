package handler

import (
	"net/http"

	"github.com/WeCanHearYou/wchy/context"
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
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": c.MustGet("Tenant"),
		})
	}
}
