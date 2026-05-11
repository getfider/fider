package handlers

import (
	"net/http"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// RoadmapPage renders the roadmap board. Pro tenants and self-hosted
// installations get the full data; other tenants render the page with no data
// so the client shows the upgrade call-to-action.
func RoadmapPage() web.HandlerFunc {
	return func(c *web.Context) error {
		props := web.Props{
			Page:  "Roadmap/Roadmap.page",
			Title: "Roadmap",
		}

		if env.IsSingleHostMode() || c.Tenant().IsPro {
			plannedPosts := &query.SearchPosts{View: "planned"}
			startedPosts := &query.SearchPosts{View: "started"}
			completedPosts := &query.SearchPosts{View: "completed"}
			getAllTags := &query.GetAllTags{}

			if err := bus.Dispatch(c, plannedPosts, startedPosts, completedPosts, getAllTags); err != nil {
				return c.Failure(err)
			}

			props.Data = web.Map{
				"plannedPosts":   plannedPosts.Result,
				"startedPosts":   startedPosts.Result,
				"completedPosts": completedPosts.Result,
				"tags":           getAllTags.Result,
			}
		}

		return c.Page(http.StatusOK, props)
	}
}
