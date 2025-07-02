package handlers

import (
	"net/http"
	"strconv"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// ModerationPage is the moderation administration page
func ModerationPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/Moderation.page",
			Title: "Moderation · Site Settings",
		})
	}
}

// GetModerationItems returns all unmoderated posts and comments
func GetModerationItems() web.HandlerFunc {
	return func(c *web.Context) error {
		q := &query.GetModerationItems{}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"items": q.Result,
		})
	}
}

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

// BulkApproveItems approves multiple posts and comments
func BulkApproveItems() web.HandlerFunc {
	return func(c *web.Context) error {
		var request struct {
			PostIDs    []string `json:"postIds"`
			CommentIDs []string `json:"commentIds"`
		}

		if err := c.Bind(&request); err != nil {
			return c.BadRequest(web.Map{"error": "Invalid request format"})
		}

		postIDs := make([]int, 0, len(request.PostIDs))
		for _, idStr := range request.PostIDs {
			if id, err := strconv.Atoi(idStr); err == nil {
				postIDs = append(postIDs, id)
			}
		}

		commentIDs := make([]int, 0, len(request.CommentIDs))
		for _, idStr := range request.CommentIDs {
			if id, err := strconv.Atoi(idStr); err == nil {
				commentIDs = append(commentIDs, id)
			}
		}

		if err := bus.Dispatch(c, &cmd.BulkApproveItems{
			PostIDs:    postIDs,
			CommentIDs: commentIDs,
		}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// BulkDeclineItems declines (deletes) multiple posts and comments
func BulkDeclineItems() web.HandlerFunc {
	return func(c *web.Context) error {
		var request struct {
			PostIDs    []string `json:"postIds"`
			CommentIDs []string `json:"commentIds"`
		}

		if err := c.Bind(&request); err != nil {
			return c.BadRequest(web.Map{"error": "Invalid request format"})
		}

		postIDs := make([]int, 0, len(request.PostIDs))
		for _, idStr := range request.PostIDs {
			if id, err := strconv.Atoi(idStr); err == nil {
				postIDs = append(postIDs, id)
			}
		}

		commentIDs := make([]int, 0, len(request.CommentIDs))
		for _, idStr := range request.CommentIDs {
			if id, err := strconv.Atoi(idStr); err == nil {
				commentIDs = append(commentIDs, id)
			}
		}

		if err := bus.Dispatch(c, &cmd.BulkDeclineItems{
			PostIDs:    postIDs,
			CommentIDs: commentIDs,
		}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}