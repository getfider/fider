package handlers

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/web"
)

// Index is the default home page
func Index() web.HandlerFunc {
	return func(c web.Context) error {
		ideas, err := c.Services().Ideas.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Map{
			"ideas": ideas,
		})
	}
}

// PostIdea creates a new idea on current tenant
func PostIdea() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.CreateNewIdea)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		ideas := c.Services().Ideas
		idea, err := ideas.Add(input.Model.Title, input.Model.Description, c.User().ID)
		if err != nil {
			return c.Failure(err)
		}

		if err := ideas.AddSupporter(idea.Number, c.User().ID); err != nil {
			return c.Failure(err)
		}

		return c.Ok(idea)
	}
}

// UpdateIdea updates an existing ideaof current tenant
func UpdateIdea() web.HandlerFunc {
	return func(c web.Context) error {
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		input := new(actions.CreateNewIdea)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		idea, err := c.Services().Ideas.GetByNumber(number)
		if idea.CanBeChangedBy(c.User()) {
			idea, err = c.Services().Ideas.Update(number, input.Model.Title, input.Model.Description)
			if err != nil {
				return c.Failure(err)
			}
		} else {
			return c.Unauthorized()
		}

		return c.Ok(idea)
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

		return c.Page(web.Map{
			"comments": comments,
			"idea":     idea,
		})
	}
}

// PostComment creates a new comment on given idea
func PostComment() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.AddNewComment)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		ideaNumber, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		ideas := c.Services().Ideas
		_, err = ideas.AddComment(ideaNumber, input.Model.Content, c.User().ID)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// SetResponse changes current idea staff response
func SetResponse() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.SetResponse)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		ideaNumber, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		ideas := c.Services().Ideas
		err = ideas.SetResponse(ideaNumber, input.Model.Text, c.User().ID, input.Model.Status)
		if err != nil {
			if err == app.ErrNotFound {
				return c.NotFound()
			}
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
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

	return c.Ok(web.Map{})
}
