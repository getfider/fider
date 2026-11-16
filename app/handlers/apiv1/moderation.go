package apiv1

import (
	"strconv"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// ApprovePost approves a post
func ApprovePost() web.HandlerFunc {
	return func(c *web.Context) error {
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid post ID"})
		}

		if err := approvePost(c, postID); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// DeclinePost declines (deletes) a post
func DeclinePost() web.HandlerFunc {
	return func(c *web.Context) error {
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid post ID"})
		}

		if err := bus.Dispatch(c, &cmd.DeclinePost{PostID: postID}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ApproveComment approves a comment
func ApproveComment() web.HandlerFunc {
	return func(c *web.Context) error {
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid comment ID"})
		}

		if err := approveComment(c, commentID); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// DeclineComment declines (deletes) a comment
func DeclineComment() web.HandlerFunc {
	return func(c *web.Context) error {
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid comment ID"})
		}

		if err := bus.Dispatch(c, &cmd.DeclineComment{CommentID: commentID}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// DeclinePostAndBlock declines (deletes) a post and blocks the user
func DeclinePostAndBlock() web.HandlerFunc {
	return func(c *web.Context) error {
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid post ID"})
		}

		getPost := &query.GetPostByID{PostID: postID}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.BlockUser{UserID: getPost.Result.User.ID})
		if err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.DeclinePost{PostID: postID})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// DeclineCommentAndBlock declines (deletes) a comment and blocks the user
func DeclineCommentAndBlock() web.HandlerFunc {
	return func(c *web.Context) error {
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid comment ID"})
		}

		getComment := &query.GetCommentByID{CommentID: commentID}
		if err := bus.Dispatch(c, getComment); err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.BlockUser{UserID: getComment.Result.User.ID})
		if err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.DeclineComment{CommentID: commentID})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ApprovePostAndVerify approves a post and trusts the user
func ApprovePostAndVerify() web.HandlerFunc {
	return func(c *web.Context) error {
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid post ID"})
		}

		getPost := &query.GetPostByID{PostID: postID}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		if err := approvePost(c, postID); err != nil {
			return c.Failure(err)
		}

		if err := bus.Dispatch(c, &cmd.TrustUser{UserID: getPost.Result.User.ID}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ApproveCommentAndVerify approves a comment and trusts the user
func ApproveCommentAndVerify() web.HandlerFunc {
	return func(c *web.Context) error {
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid comment ID"})
		}

		getComment := &query.GetCommentByID{CommentID: commentID}
		if err := bus.Dispatch(c, getComment); err != nil {
			return c.Failure(err)
		}

		if err := approveComment(c, commentID); err != nil {
			return c.Failure(err)
		}

		if err := bus.Dispatch(c, &cmd.TrustUser{UserID: getComment.Result.User.ID}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// approvePost approves a pending post and notifies the subscribers that were held back while it was pending
func approvePost(c *web.Context, postID int) error {
	approve := &cmd.ApprovePost{PostID: postID}
	if err := bus.Dispatch(c, approve); err != nil || !approve.Result {
		return err
	}

	getPost := &query.GetPostByID{PostID: postID}
	if err := bus.Dispatch(c, getPost); err != nil {
		return err
	}

	c.Enqueue(tasks.NotifyAboutApprovedPost(getPost.Result))
	return nil
}

// approveComment approves a pending comment and notifies the subscribers that were held back while it was pending
func approveComment(c *web.Context, commentID int) error {
	approve := &cmd.ApproveComment{CommentID: commentID}
	if err := bus.Dispatch(c, approve); err != nil || approve.Result == 0 {
		return err
	}

	getPost := &query.GetPostByID{PostID: approve.Result}
	getComment := &query.GetCommentByID{CommentID: commentID}
	if err := bus.Dispatch(c, getPost, getComment); err != nil {
		return err
	}

	// A comment on a still-pending post stays staff-only
	if getPost.Result.IsApproved {
		c.Enqueue(tasks.NotifyAboutApprovedComment(getComment.Result, getPost.Result))
	}
	return nil
}
