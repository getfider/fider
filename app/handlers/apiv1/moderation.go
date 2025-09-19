package apiv1

import (
	"strconv"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// ApprovePost approves a post
func ApprovePost() web.HandlerFunc {
	return func(c *web.Context) error {
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid post ID"})
		}

		if err := bus.Dispatch(c, &cmd.ApprovePost{PostID: postID}); err != nil {
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

		if err := bus.Dispatch(c, &cmd.ApproveComment{CommentID: commentID}); err != nil {
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

		// Get the post that is being declined
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid post ID"})
		}

		getPost := &query.GetPostByID{PostID: postID}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		// Call the existing BlockUser command inside the DeclinePostAndBlock handler
		err = bus.Dispatch(c, &cmd.BlockUser{UserID: getPost.Result.User.ID})
		if err != nil {
			return c.Failure(err)
		}

		// Finally call the existing DeclinePost command
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
		// Get the comment that is being declined
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.BadRequest(web.Map{"error": "Invalid comment ID"})
		}

		getComment := &query.GetCommentByID{CommentID: commentID}
		if err := bus.Dispatch(c, getComment); err != nil {
			return c.Failure(err)
		}

		// Call the existing BlockUser command inside the DeclinePostAndBlock handler
		err = bus.Dispatch(c, &cmd.BlockUser{UserID: getComment.Result.User.ID})
		if err != nil {
			return c.Failure(err)
		}

		// Finally call the existing DeclinePost command
		err = bus.Dispatch(c, &cmd.DeclineComment{CommentID: commentID})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
