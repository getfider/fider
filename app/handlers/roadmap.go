package handlers

import (
	"net/http"

	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// Roadmap displays posts grouped by status in a kanban-style view
func Roadmap() web.HandlerFunc {
	return func(c *web.Context) error {
		if !c.Tenant().IsRoadmapEnabled {
			return c.NotFound()
		}

		postsByStatus := &query.GetPostsByStatuses{
			Statuses: []enum.PostStatus{
				enum.PostPlanned,
				enum.PostStarted,
				enum.PostCompleted,
			},
		}

		if err := bus.Dispatch(c, postsByStatus); err != nil {
			return c.Failure(err)
		}

		return c.Page(http.StatusOK, web.Props{
			Page:  "Roadmap/Roadmap.page",
			Title: "Roadmap",
			Data: web.Map{
				"planned":   postsByStatus.Result[enum.PostPlanned],
				"started":   postsByStatus.Result[enum.PostStarted],
				"completed": postsByStatus.Result[enum.PostCompleted],
			},
		})
	}
}
