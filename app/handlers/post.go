package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/csv"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// Index is the default home page
func Index() web.HandlerFunc {
	return func(c web.Context) error {
		posts, err := c.Services().Posts.Search(
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

		stats, err := c.Services().Posts.CountPerStatus()
		if err != nil {
			return c.Failure(err)
		}

		description := ""
		if c.Tenant().WelcomeMessage != "" {
			description = markdown.PlainText(c.Tenant().WelcomeMessage)
		} else {
			description = "We'd love to hear what you're thinking about. What can we do better? This is the place for you to vote, discuss and share posts."
		}

		return c.Page(web.Props{
			Description: description,
			Data: web.Map{
				"posts":          posts,
				"tags":           tags,
				"countPerStatus": stats,
			},
		})
	}
}

// PostDetails shows details of given Post by id
func PostDetails() web.HandlerFunc {
	return func(c web.Context) error {
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.Failure(err)
		}

		posts := c.Services().Posts
		post, err := posts.GetByNumber(number)
		if err != nil {
			return c.Failure(err)
		}

		comments, err := posts.GetCommentsByPost(post)
		if err != nil {
			return c.Failure(err)
		}

		tags, err := c.Services().Tags.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		subscribed, err := c.Services().Users.HasSubscribedTo(post.ID)
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:       post.Title,
			Description: markdown.PlainText(post.Description),
			Data: web.Map{
				"comments":   comments,
				"subscribed": subscribed,
				"post":       post,
				"tags":       tags,
			},
		})
	}
}

// PostComment creates a new comment on given post
func PostComment() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.AddNewComment)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		post, err := c.Services().Posts.GetByNumber(input.Model.Number)
		if err != nil {
			return c.Failure(err)
		}

		_, err = c.Services().Posts.AddComment(post, input.Model.Content)
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutNewComment(post, input.Model))

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

		err := c.Services().Posts.UpdateComment(input.Model.ID, input.Model.Content)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ExportPostsToCSV returns a CSV with all posts
func ExportPostsToCSV() web.HandlerFunc {
	return func(c web.Context) error {

		posts, err := c.Services().Posts.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		bytes, err := csv.FromPosts(posts)
		if err != nil {
			return c.Failure(err)
		}

		return c.Attachment("posts.csv", "text/csv", bytes)
	}
}
