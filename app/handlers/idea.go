package handlers

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// Index is the default home page
func Index() web.HandlerFunc {
	return func(c web.Context) error {
		ideas, err := c.Services().Ideas.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		tags, err := c.Services().Tags.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Map{
			"ideas": ideas,
			"tags":  tags,
		})
	}
}

// GetIdeas return basic model of all tenant ideas
func GetIdeas() web.HandlerFunc {
	return func(c web.Context) error {
		ideas, err := c.Services().Ideas.GetAllBasic()
		if err != nil {
			return c.Failure(err)
		}
		return c.Ok(ideas)
	}
}

// SearchIdeas return existing ideas based on search criteria
func SearchIdeas() web.HandlerFunc {
	return func(c web.Context) error {
		query := c.QueryParam("q")
		c.Logger().Infof(query)
		ideas, err := c.Services().Ideas.Search(query)
		if err != nil {
			return c.Failure(err)
		}
		return c.Ok(ideas)
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

		c.Enqueue(tasks.NotifyAboutNewIdea(idea))

		return c.Ok(idea)
	}
}

// UpdateIdea updates an existing ideaof current tenant
func UpdateIdea() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.UpdateIdea)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		_, err := c.Services().Ideas.Update(input.Model.Number, input.Model.Title, input.Model.Description)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
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
			return c.Failure(err)
		}

		comments, err := ideas.GetCommentsByIdea(number)
		if err != nil {
			return c.Failure(err)
		}

		tags, err := c.Services().Tags.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		subscribed, err := c.Services().Users.HasSubscribedTo(idea.ID)
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Map{
			"comments":   comments,
			"subscribed": subscribed,
			"idea":       idea,
			"tags":       tags,
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

		idea, err := c.Services().Ideas.GetByNumber(input.Model.Number)
		if err != nil {
			return c.Failure(err)
		}

		_, err = c.Services().Ideas.AddComment(input.Model.Number, input.Model.Content, c.User().ID)
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutNewComment(idea, input.Model))

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

		idea, err := c.Services().Ideas.GetByNumber(input.Model.Number)
		if err != nil {
			return c.Failure(err)
		}

		if input.Model.Status == models.IdeaDuplicate {
			err = c.Services().Ideas.MarkAsDuplicate(input.Model.Number, input.Original.Number, c.User().ID)
		} else {
			err = c.Services().Ideas.SetResponse(input.Model.Number, input.Model.Text, c.User().ID, input.Model.Status)
		}
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutStatusChange(idea, input.Model))

		return c.Ok(web.Map{})
	}
}

// AddSupporter adds current user to given idea list of supporters
func AddSupporter() web.HandlerFunc {
	return func(c web.Context) error {
		return addOrRemove(c, c.Services().Ideas.AddSupporter)
	}
}

// RemoveSupporter removes current user from given idea list of supporters
func RemoveSupporter() web.HandlerFunc {
	return func(c web.Context) error {
		return addOrRemove(c, c.Services().Ideas.RemoveSupporter)
	}
}

// Subscribe adds current user to list of subscribers of given idea
func Subscribe() web.HandlerFunc {
	return func(c web.Context) error {
		return addOrRemove(c, c.Services().Ideas.AddSubscriber)
	}
}

// Unsubscribe removes current user from list of subscribers of given idea
func Unsubscribe() web.HandlerFunc {
	return func(c web.Context) error {
		return addOrRemove(c, c.Services().Ideas.RemoveSubscriber)
	}
}

func addOrRemove(c web.Context, addOrRemove func(number, userID int) error) error {
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
