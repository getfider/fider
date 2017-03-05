package feedback

import (
	"github.com/WeCanHearYou/wchy/identity"
	"github.com/labstack/echo"
)

type indexHandlder struct {
	ideaService IdeaService
}

// Index handles initial page
func Index(ideaService IdeaService) echo.HandlerFunc {
	return indexHandlder{ideaService}.get()
}

func (h indexHandlder) get() echo.HandlerFunc {
	return func(c echo.Context) error {
		tenant := c.Get("Tenant").(*identity.Tenant)
		ideas, _ := h.ideaService.GetAll(tenant.ID)
		return c.Render(200, "index.html", echo.Map{
			"Tenant": tenant,
			"Ideas":  ideas,
		})
	}
}
