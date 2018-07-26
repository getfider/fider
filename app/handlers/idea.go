package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/csv"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// Index is the default home page
func Index() web.HandlerFunc {
	return func(c web.Context) error {
		ideas, err := c.Services().Ideas.Search(
			c.QueryParam("q"),
			c.QueryParam("f"),
			c.QueryParam("l"),
			c.QueryParamAsArray("t"),
		)
		if err != nil {
			return c.Failure(err)
		}

		tags, err := c.Services().Tags.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		stats, err := c.Services().Ideas.CountPerStatus()
		if err != nil {
			return c.Failure(err)
		}

		description := ""
		if c.Tenant().WelcomeMessage != "" {
			description = markdown.PlainText(c.Tenant().WelcomeMessage)
		} else {
			description = "We'd love to hear what you're thinking about. What can we do better? This is the place for you to vote, discuss and share ideas."
		}

		return c.Page(web.Props{
			Description: description,
			Data: web.Map{
				"ideas":          ideas,
				"tags":           tags,
				"countPerStatus": stats,
			},
		})
	}
}

// SearchIdeas return existing ideas based on search criteria
func SearchIdeas() web.HandlerFunc {
	return func(c web.Context) error {
		ideas, err := c.Services().Ideas.Search(
			c.QueryParam("q"),
			c.QueryParam("f"),
			c.QueryParam("l"),
			c.QueryParamAsArray("t"),
		)
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
		idea, err := ideas.Add(input.Model.Title, input.Model.Description)
		if err != nil {
			return c.Failure(err)
		}

		if err := ideas.AddSupporter(idea, c.User()); err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutNewIdea(idea))

		return c.Ok(idea)
	}
}

// UpdateIdea updates an existing idea of current tenant
func UpdateIdea() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.UpdateIdea)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		_, err := c.Services().Ideas.Update(input.Idea, input.Model.Title, input.Model.Description)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// DeleteIdea deletes an existing idea of current tenant
func DeleteIdea() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.DeleteIdea)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Ideas.SetResponse(input.Idea, input.Model.Text, models.IdeaDeleted)
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

		comments, err := ideas.GetCommentsByIdea(idea)
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

		return c.Page(web.Props{
			Title:       idea.Title,
			Description: markdown.PlainText(idea.Description),
			Data: web.Map{
				"comments":   comments,
				"subscribed": subscribed,
				"idea":       idea,
				"tags":       tags,
			},
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

		_, err = c.Services().Ideas.AddComment(idea, input.Model.Content)
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutNewComment(idea, input.Model))

		return c.Ok(web.Map{})
	}
}

// UpdateComment changes an existing comment with new content
func UpdateComment() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.EditComment)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Ideas.UpdateComment(input.Model.ID, input.Model.Content)
		if err != nil {
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

		idea, err := c.Services().Ideas.GetByNumber(input.Model.Number)
		if err != nil {
			return c.Failure(err)
		}

		prevStatus := idea.Status
		if input.Model.Status == models.IdeaDuplicate {
			err = c.Services().Ideas.MarkAsDuplicate(idea, input.Original)
		} else {
			err = c.Services().Ideas.SetResponse(idea, input.Model.Text, input.Model.Status)
		}
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutStatusChange(idea, prevStatus))

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

func addOrRemove(c web.Context, addOrRemove func(idea *models.Idea, user *models.User) error) error {
	number, err := c.ParamAsInt("number")
	if err != nil {
		return c.Failure(err)
	}

	idea, err := c.Services().Ideas.GetByNumber(number)
	if err != nil {
		return c.Failure(err)
	}

	err = addOrRemove(idea, c.User())
	if err != nil {
		return c.Failure(err)
	}

	return c.Ok(web.Map{})
}

// ExportIdeasToCSV returns a CSV with all ideas
func ExportIdeasToCSV() web.HandlerFunc {
	return func(c web.Context) error {

		ideas, err := c.Services().Ideas.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		bytes, err := csv.FromPosts(ideas)
		if err != nil {
			return c.Failure(err)
		}

		return c.Attachment("ideas.csv", "text/csv", bytes)
	}
}
