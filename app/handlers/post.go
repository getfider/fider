package handlers

import (
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/csv"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/web"
)

// Index is the default home page
func Index() web.HandlerFunc {
	return func(c *web.Context) error {
		c.SetCanonicalURL("")

		posts, err := c.Services().Posts.Search(
			c.QueryParam("query"),
			c.QueryParam("view"),
			c.QueryParam("limit"),
			c.QueryParamAsArray("tags"),
		)
		if err != nil {
			return c.Failure(err)
		}

		getAllTags := &query.GetAllTags{}
		countPerStatus := &query.CountPostPerStatus{}
		if err := bus.Dispatch(c, getAllTags, countPerStatus); err != nil {
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
			ChunkName:   "Home.page",
			Data: web.Map{
				"posts":          posts,
				"tags":           getAllTags.Result,
				"countPerStatus": countPerStatus.Result,
			},
		})
	}
}

// PostDetails shows details of given Post by id
func PostDetails() web.HandlerFunc {
	return func(c *web.Context) error {
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

		subscribed, err := c.Services().Users.HasSubscribedTo(post.ID)
		if err != nil {
			return c.Failure(err)
		}

		getAllTags := &query.GetAllTags{}
		listVotes := &query.ListPostVotes{PostID: post.ID, Limit: 6}
		getAttachments := &query.GetAttachments{Post: post}
		if err := bus.Dispatch(c, getAllTags, listVotes, getAttachments); err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:       post.Title,
			Description: markdown.PlainText(post.Description),
			ChunkName:   "ShowPost.page",
			Data: web.Map{
				"comments":    comments,
				"subscribed":  subscribed,
				"post":        post,
				"tags":        getAllTags.Result,
				"votes":       listVotes.Result,
				"attachments": getAttachments.Result,
			},
		})
	}
}

// ExportPostsToCSV returns a CSV with all posts
func ExportPostsToCSV() web.HandlerFunc {
	return func(c *web.Context) error {

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
