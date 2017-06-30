package handlers

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/labstack/echo"
)

// Index is the default home page
func Index() web.HandlerFunc {
	return func(c web.Context) error {
		ideas, err := c.Services().Ideas.GetAll()
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
				"message": "Title is required.",
			})
		}

		if len(input.Title) < 10 || len(strings.Split(input.Title, " ")) < 3 {
			return c.JSON(400, echo.Map{
				"message": "Title needs to be more descriptive.",
			})
		}

		ideas := c.Services().Ideas
		idea, err := ideas.Add(input.Title, input.Description, c.User().ID)
		if err != nil {
			return c.Failure(err)
		}

		if err := ideas.AddSupporter(idea.Number, c.User().ID); err != nil {
			return c.Failure(err)
		}

		return c.JSON(200, idea)
	}
}

// IdeaDetails shows details of given Idea by id
func IdeaDetails() web.HandlerFunc {
	return func(c web.Context) error {
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		ideas := c.Services().Ideas
		idea, err := ideas.GetByNumber(number)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
			return c.Failure(err)
		}

		comments, err := ideas.GetCommentsByIdea(number)
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
		_, err = ideas.AddComment(ideaNumber, input.Content, c.User().ID)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
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

		if input.Status < models.IdeaNew || input.Status > models.IdeaDeclined {
			return c.JSON(400, echo.Map{
				"message": "Status is invalid.",
			})
		}

		if strings.Trim(input.Text, " ") == "" {
			return c.JSON(400, echo.Map{
				"message": "Text is required.",
			})
		}

		ideaNumber, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		ideas := c.Services().Ideas
		err = ideas.SetResponse(ideaNumber, input.Text, c.User().ID, input.Status)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
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

func addOrRemoveSupporter(c web.Context, addOrRemove func(number, userID int) error) error {
	ideaNumber, err := c.ParamAsInt("number")
	if err != nil {
		return c.Failure(err)
	}

	err = addOrRemove(ideaNumber, c.User().ID)
	if err != nil {
		if err == app.ErrNotFound {
			return c.NotFound()
		}
		return c.Failure(err)
	}

	return c.JSON(200, echo.Map{})
}
