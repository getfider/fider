package apiv1

import (
	"time"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// AllComments returns a list of all comments regardless of the post
func AllComments() web.HandlerFunc {
	return func(c *web.Context) error {
		var since time.Time

		if sinceParam, err := time.Parse(time.RFC3339, c.QueryParam("since")); err == nil {
			since = sinceParam
		}

		getComments := &query.GetCommentRefs{Since: since}
		if err := bus.Dispatch(c, getComments); err != nil {
			return c.Failure(err)
		}

		return c.Ok(getComments.Result)
	}
}
