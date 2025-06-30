package apiv1

import (
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// GetPostDetails returns all details needed to display a post
func GetPostDetails() web.HandlerFunc {
	return func(c *web.Context) error {
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.NotFound()
		}

		getPost := &query.GetPostByNumber{Number: number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		isSubscribed := &query.UserSubscribedTo{PostID: getPost.Result.ID}
		getComments := &query.GetCommentsByPost{Post: getPost.Result}
		getAllTags := &query.GetAllTags{}
		listVotes := &query.ListPostVotes{PostID: getPost.Result.ID, Limit: 24, IncludeEmail: false}
		getAttachments := &query.GetAttachments{Post: getPost.Result}
		if err := bus.Dispatch(c, getAllTags, getComments, listVotes, isSubscribed, getAttachments); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"comments":    getComments.Result,
			"subscribed":  isSubscribed.Result,
			"post":        getPost.Result,
			"tags":        getAllTags.Result,
			"votes":       listVotes.Result,
			"attachments": getAttachments.Result,
		})
	}
}
