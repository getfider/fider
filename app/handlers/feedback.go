package handlers

import (
	"strings"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/WeCanHearYou/wechy/app/storage"
	"github.com/labstack/echo"
)

// AllHandlers contains multiple feedback HTTP handlers
type AllHandlers struct {
	ideas storage.Idea
}

// Handlers handles feedback based page
func Handlers(ideas storage.Idea) AllHandlers {
	return AllHandlers{ideas}
}

// List all tenant ideas
func (h AllHandlers) List() web.HandlerFunc {
	return func(c web.Context) error {
		ideas, err := h.ideas.GetAll(c.Tenant().ID)
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(echo.Map{
			"ideas": ideas,
		})
	}
}

type newIdeaInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// PostIdea creates a new idea on current tenant
func (h AllHandlers) PostIdea() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(newIdeaInput)
		if err := c.Bind(input); err != nil {
			return c.Failure(err)
		}

		if strings.Trim(input.Title, " ") == "" {
			return c.JSON(400, echo.Map{
				"message": "Idea title is required.",
			})
		}

		idea, err := h.ideas.Save(c.Tenant().ID, c.User().ID, input.Title, input.Description)
		if err != nil {
			return c.Failure(err)
		}

		if err := h.ideas.AddSupporter(c.User().ID, idea.ID); err != nil {
			return c.Failure(err)
		}

		return c.JSON(200, idea)
	}
}

// Details shows details of given Idea by id
func (h AllHandlers) Details() web.HandlerFunc {
	return func(c web.Context) error {
		tenant := c.Tenant()

		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		idea, err := h.ideas.GetByNumber(tenant.ID, number)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
			return c.Failure(err)
		}

		comments, err := h.ideas.GetCommentsByIdeaID(tenant.ID, idea.ID)
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(echo.Map{
			"comments": comments,
			"idea":     idea,
		})
	}
}

type newCommentInput struct {
	Content string `json:"content"`
}

// PostComment creates a new comment on given idea
func (h AllHandlers) PostComment() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(newCommentInput)
		if err := c.Bind(input); err != nil {
			return c.Failure(err)
		}

		if strings.Trim(input.Content, " ") == "" {
			return c.JSON(400, echo.Map{
				"message": "Comment is required.",
			})
		}

		ideaNumber, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		idea, err := h.ideas.GetByNumber(c.Tenant().ID, ideaNumber)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
			return c.Failure(err)
		}

		_, err = h.ideas.AddComment(c.User().ID, idea.ID, input.Content)
		if err != nil {
			return c.Failure(err)
		}

		return c.JSON(200, echo.Map{})
	}
}

// AddSupporter adds current user to given idea list of supporters
func (h AllHandlers) AddSupporter() web.HandlerFunc {
	return func(c web.Context) error {
		return addOrRemoveSupporter(c, h, h.ideas.AddSupporter)
	}
}

// RemoveSupporter removes current user from given idea list of supporters
func (h AllHandlers) RemoveSupporter() web.HandlerFunc {
	return func(c web.Context) error {
		return addOrRemoveSupporter(c, h, h.ideas.RemoveSupporter)
	}
}

func addOrRemoveSupporter(c web.Context, h AllHandlers, addOrRemove func(userId, ideaId int) error) error {
	ideaNumber, err := c.ParamAsInt("number")
	if err != nil {
		return c.Failure(err)
	}

	idea, err := h.ideas.GetByNumber(c.Tenant().ID, ideaNumber)
	if err != nil {
		if err == app.ErrNotFound {
			return c.NotFound()
		}
		return c.Failure(err)
	}

	err = addOrRemove(c.User().ID, idea.ID)
	if err != nil {
		return c.Failure(err)
	}

	return c.JSON(200, echo.Map{})
}
