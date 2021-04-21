package handlers

import (
	"fmt"

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

		searchPosts := &query.SearchPosts{
			Query: c.QueryParam("query"),
			View:  c.QueryParam("view"),
			Limit: c.QueryParam("limit"),
			Tags:  c.QueryParamAsArray("tags"),
		}
		getAllTags := &query.GetAllTags{}
		countPerStatus := &query.CountPostPerStatus{}

		if err := bus.Dispatch(c, searchPosts, getAllTags, countPerStatus); err != nil {
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
				"posts":          searchPosts.Result,
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

		getPost := &query.GetPostByNumber{Number: number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		if c.Param("slug") != getPost.Result.Slug {
			return c.Redirect(fmt.Sprintf("/posts/%d/%s", getPost.Result.Number, getPost.Result.Slug))
		}

		isSubscribed := &query.UserSubscribedTo{PostID: getPost.Result.ID}
		getComments := &query.GetCommentsByPost{Post: getPost.Result}
		getAllTags := &query.GetAllTags{}
		listVotes := &query.ListPostVotes{PostID: getPost.Result.ID, Limit: 6, IncludeEmail: false}
		getAttachments := &query.GetAttachments{Post: getPost.Result}
		if err := bus.Dispatch(c, getAllTags, getComments, listVotes, isSubscribed, getAttachments); err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:       getPost.Result.Title,
			Description: markdown.PlainText(getPost.Result.Description),
			ChunkName:   "ShowPost.page",
			Data: web.Map{
				"comments":    getComments.Result,
				"subscribed":  isSubscribed.Result,
				"post":        getPost.Result,
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

		allPosts := &query.GetAllPosts{}
		if err := bus.Dispatch(c, allPosts); err != nil {
			return c.Failure(err)
		}

		bytes, err := csv.FromPosts(allPosts.Result)
		if err != nil {
			return c.Failure(err)
		}

		return c.Attachment("posts.csv", "text/csv", bytes)
	}
}
