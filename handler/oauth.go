package handler

import (
	"net/http"
	"time"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/labstack/echo"
)

type oauthHandler struct {
	ctx *context.WchyContext
}

// OAuth creates a new OAuth HTTP handler
func OAuth(ctx *context.WchyContext) echo.HandlerFunc {
	return oauthHandler{ctx: ctx}.get()
}

func (h oauthHandler) get() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"now": time.Now().Format("2006.01.02.150405"),
		})
	}
}
