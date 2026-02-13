package handlers

import (
	"fmt"
	"net/http"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/csv"
	"github.com/getfider/fider/app/pkg/env"
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

		if myVotesOnly, err := c.QueryParamAsBool("myvotes"); err == nil {
			searchPosts.MyVotesOnly = myVotesOnly
		}

		if noTagsOnly, err := c.QueryParamAsBool("notags"); err == nil {
			searchPosts.NoTagsOnly = noTagsOnly
		}

		if myPostsOnly, err := c.QueryParamAsBool("myposts"); err == nil {
			searchPosts.MyPostsOnly = myPostsOnly
		}

		// Handle "pending" pseudo-status for moderation filtering
		statusesParam := c.QueryParamAsArray("statuses")
		hasPending := false
		actualStatuses := []string{}
		for _, status := range statusesParam {
			if status == "pending" {
				hasPending = true
			} else {
				actualStatuses = append(actualStatuses, status)
			}
		}

		// Set moderation filter based on pending status
		if hasPending {
			searchPosts.ModerationFilter = "pending"
		}

		searchPosts.SetStatusesFromStrings(actualStatuses)
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

		data := web.Map{
			"searchNoiseWords": env.SearchNoiseWords(),
			"posts":            searchPosts.Result,
			"tags":             getAllTags.Result,
			"countPerStatus":   countPerStatus.Result,
		}

		return c.Page(http.StatusOK, web.Props{
			Page:        "Home/Home.page",
			Description: description,
			// Header:      c.Tenant().WelcomeHeader,
			Data: data,
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
			return c.Redirect(c.BaseURL() + fmt.Sprintf("/posts/%d/%s", getPost.Result.Number, getPost.Result.Slug))
		}

		isSubscribed := &query.UserSubscribedTo{PostID: getPost.Result.ID}
		getComments := &query.GetCommentsByPost{Post: getPost.Result}
		getAllTags := &query.GetAllTags{}
		listVotes := &query.ListPostVotes{PostID: getPost.Result.ID, Limit: 24, IncludeEmail: false}
		getAttachments := &query.GetAttachments{Post: getPost.Result}
		if err := bus.Dispatch(c, getAllTags, getComments, listVotes, isSubscribed, getAttachments); err != nil {
			return c.Failure(err)
		}

		return c.Page(http.StatusOK, web.Props{
			Page:        "ShowPost/ShowPost.page",
			Title:       getPost.Result.Title,
			Description: markdown.PlainText(getPost.Result.Description),
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
