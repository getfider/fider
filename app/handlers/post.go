package handlers

import (
	"github.com/getfider/fider/app/pkg/csv"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/web"
)

// Index is the default home page
func Index() web.HandlerFunc {
	return func(c web.Context) error {
		posts, err := c.Services().Posts.Search(
			c.QueryParam("query"),
			c.QueryParam("view"),
			c.QueryParam("limit"),
			c.QueryParamAsArray("tags"),
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
			return c.NotFound()
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

		votes, err := c.Services().Posts.ListVotes(post, 6)
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
				"votes":      votes,
			},
		})
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
