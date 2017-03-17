package feedback

import (
	"github.com/WeCanHearYou/wechy/identity"
	"github.com/labstack/echo"
)

// IndexHandlder contains multiple feedback HTTP handlers
type IndexHandlder struct {
	ideaService IdeaService
}

// Index handles initial page
func Index(ideaService IdeaService) IndexHandlder {
	return IndexHandlder{ideaService}
}

// List all tenant ideas
func (h IndexHandlder) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		tenant := c.Get("Tenant").(*identity.Tenant)
		ideas, _ := h.ideaService.GetAll(tenant.ID)
		return c.Render(200, "index.html", echo.Map{
			"Tenant": tenant,
			"Ideas":  ideas,
		})
	}
}

type newIdea struct {
	Title string `json:"title"`
}

// Post creates a new idea on current tenant
func (h IndexHandlder) Post() echo.HandlerFunc {
	return func(c echo.Context) error {
		idea := new(newIdea)
		c.Bind(idea)

		return c.JSON(200, echo.Map{
			"tenant": c.Get("Tenant"),
			"claims": c.Get("Claims"),
			"idea":   idea.Title,
		})
	}
}
