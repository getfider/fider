package feedback

import (
	"net/http"

	"github.com/WeCanHearYou/wechy/app/identity"
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
		ideas, err := h.ideaService.GetAll(tenant.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.Render(200, "index.html", echo.Map{
			"Tenant": tenant,
			"Ideas":  ideas,
		})
	}
}

type newIdeaInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Post creates a new idea on current tenant
func (h IndexHandlder) Post() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := new(newIdeaInput)
		c.Bind(input)

		tenant := c.Get("Tenant").(*identity.Tenant)
		claims := c.Get("Claims").(*identity.WechyClaims)

		idea, err := h.ideaService.Save(tenant.ID, claims.UserID, input.Title, input.Description)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(200, echo.Map{
			"idea": idea,
		})
	}
}
