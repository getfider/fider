package handlers

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/labstack/echo"
)

// Index is the default home page
func Index() web.HandlerFunc {
	return func(c web.Context) error {
		ideas, err := c.Services().Ideas.GetAll(c.Tenant().ID)
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
func PostIdea() web.HandlerFunc {
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

		ideas := c.Services().Ideas
		idea, err := ideas.Save(c.Tenant().ID, c.User().ID, input.Title, input.Description)
		if err != nil {
			return c.Failure(err)
		}

		if err := ideas.AddSupporter(c.Tenant().ID, c.User().ID, idea.ID); err != nil {
			return c.Failure(err)
		}

		return c.JSON(200, idea)
	}
}

// IdeaDetails shows details of given Idea by id
func IdeaDetails() web.HandlerFunc {
	return func(c web.Context) error {
		tenant := c.Tenant()

		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		ideas := c.Services().Ideas
		idea, err := ideas.GetByNumber(tenant.ID, number)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
			return c.Failure(err)
		}

		comments, err := ideas.GetCommentsByIdeaID(tenant.ID, idea.ID)
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
func PostComment() web.HandlerFunc {
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

		ideas := c.Services().Ideas
		idea, err := ideas.GetByNumber(c.Tenant().ID, ideaNumber)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
			return c.Failure(err)
		}

		_, err = ideas.AddComment(c.User().ID, idea.ID, input.Content)
		if err != nil {
			return c.Failure(err)
		}

		return c.JSON(200, echo.Map{})
	}
}

type setResponseInput struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
}

// SetResponse changes current idea staff response
func SetResponse() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(setResponseInput)
		if err := c.Bind(input); err != nil {
			return c.Failure(err)
		}

		ideaNumber, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		ideas := c.Services().Ideas
		idea, err := ideas.GetByNumber(c.Tenant().ID, ideaNumber)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
			return c.Failure(err)
		}

		err = ideas.SetResponse(c.Tenant().ID, idea.ID, input.Text, c.User().ID, input.Status)
		if err != nil {
			return c.Failure(err)
		}

		return c.JSON(200, echo.Map{})
	}
}

// AddSupporter adds current user to given idea list of supporters
func AddSupporter() web.HandlerFunc {
	return func(c web.Context) error {
		return addOrRemoveSupporter(c, c.Services().Ideas.AddSupporter)
	}
}

// RemoveSupporter removes current user from given idea list of supporters
func RemoveSupporter() web.HandlerFunc {
	return func(c web.Context) error {
		return addOrRemoveSupporter(c, c.Services().Ideas.RemoveSupporter)
	}
}

func addOrRemoveSupporter(c web.Context, addOrRemove func(tenantId, userId, ideaId int) error) error {
	ideaNumber, err := c.ParamAsInt("number")
	if err != nil {
		return c.Failure(err)
	}

	idea, err := c.Services().Ideas.GetByNumber(c.Tenant().ID, ideaNumber)
	if err != nil {
		if err == app.ErrNotFound {
			return c.NotFound()
		}
		return c.Failure(err)
	}

	err = addOrRemove(c.Tenant().ID, c.User().ID, idea.ID)
	if err != nil {
		return c.Failure(err)
	}

	return c.JSON(200, echo.Map{})
}
