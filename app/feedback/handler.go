package feedback

import (
	"strings"

	"github.com/WeCanHearYou/wechy/app"
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
func (h IndexHandlder) List() app.HandlerFunc {
	return func(c app.Context) error {
		ideas, err := h.ideaService.GetAll(c.Tenant().ID)
		if err != nil {
			return c.Failure(err)
		}

		return c.Render(200, "index.html", echo.Map{
			"Ideas": ideas,
		})
	}
}

type newIdeaInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Post creates a new idea on current tenant
func (h IndexHandlder) Post() app.HandlerFunc {
	return func(c app.Context) error {
		input := new(newIdeaInput)
		if err := c.Bind(input); err != nil {
			return c.Failure(err)
		}

		if strings.Trim(input.Title, " ") == "" {
			return c.JSON(400, echo.Map{
				"message": "Idea title is required.",
			})
		}

		idea, err := h.ideaService.Save(c.Tenant().ID, c.Claims().UserID, input.Title, input.Description)
		if err != nil {
			return c.Failure(err)
		}

		return c.JSON(200, echo.Map{
			"idea": idea,
		})
	}
}

// Details shows details of given Idea by id
func (h IndexHandlder) Details() app.HandlerFunc {
	return func(c app.Context) error {
		tenant := c.Tenant()

		ideaID, err := c.ParamAsInt("id")
		if err != nil {
			return c.Failure(err)
		}

		idea, err := h.ideaService.GetByID(tenant.ID, ideaID)
		if err != nil {
			return c.Failure(err)
		}

		comments, err := h.ideaService.GetCommentsByIdeaID(tenant.ID, idea.ID)
		if err != nil {
			return c.Failure(err)
		}

		return c.Render(200, "idea.html", echo.Map{
			"Comments": comments,
			"Idea":     idea,
		})
	}
}

type newCommentInput struct {
	Content string `json:"content"`
}

// PostComment creates a new comment on given idea
func (h IndexHandlder) PostComment() app.HandlerFunc {
	return func(c app.Context) error {
		input := new(newCommentInput)
		if err := c.Bind(input); err != nil {
			return c.Failure(err)
		}

		if strings.Trim(input.Content, " ") == "" {
			return c.JSON(400, echo.Map{
				"message": "Comment is required.",
			})
		}

		ideaID, err := c.ParamAsInt("id")
		if err != nil {
			return c.Failure(err)
		}

		id, err := h.ideaService.AddComment(c.Claims().UserID, ideaID, input.Content)
		if err != nil {
			return c.Failure(err)
		}

		return c.JSON(200, echo.Map{
			"comment_id": id,
		})
	}
}
