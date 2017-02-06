package router

import (
	"strings"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/gin-gonic/gin"
)

// MultiTenant extract tenant information from hostname and inject it into current context
func MultiTenant(ctx *context.WchyContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		hostname := stripPort(c.Request.Host)
		tenant, err := ctx.Tenant.GetByDomain(hostname)
		if err == nil {
			c.Set("Tenant", tenant.Name)
			c.Next()
		} else {
			c.AbortWithStatus(404)
		}
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	return hostport[:colon]
}
