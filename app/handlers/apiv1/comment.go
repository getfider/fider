package apiv1

import (
	"time"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

type commentsResponse struct {
	Data       []*entity.CommentRef `json:"data"`
	Pagination *paginationInfo      `json:"pagination"`
}

type paginationInfo struct {
	HasNext    bool   `json:"hasNext"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// AllComments returns a list of all comments regardless of the post
func AllComments() web.HandlerFunc {
	return func(c *web.Context) error {
		var since time.Time

		if sinceParam, err := time.Parse(time.RFC3339, c.QueryParam("since")); err == nil {
			since = sinceParam
		}

		// Default limit to 50, max 100
		limit := 50
		if limitParam, err := c.QueryParamAsInt("limit"); err == nil && limitParam > 0 {
			limit = limitParam
			if limit > 100 {
				limit = 100
			}
		}

		getComments := &query.GetCommentRefs{
			Since: since,
			Limit: limit,
		}
		if err := bus.Dispatch(c, getComments); err != nil {
			return c.Failure(err)
		}

		// Determine if there are more results by checking if we got a full page
		hasNext := len(getComments.Result) == limit
		var nextCursor string
		if hasNext && len(getComments.Result) > 0 {
			lastComment := getComments.Result[len(getComments.Result)-1]
			// Use the created_at time as the cursor, but prefer edited_at if available
			cursorTime := lastComment.CreatedAt
			if lastComment.EditedAt != nil {
				cursorTime = *lastComment.EditedAt
			}
			nextCursor = cursorTime.Format(time.RFC3339)
		}

		response := commentsResponse{
			Data: getComments.Result,
			Pagination: &paginationInfo{
				HasNext:    hasNext,
				NextCursor: nextCursor,
			},
		}

		return c.Ok(response)
	}
}
