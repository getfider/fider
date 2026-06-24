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
			// Keep in sync with ROADMAP_DEFAULT_LIMIT on the client; the "Show more"
			// link uses posts.length >= limit as its heuristic and needs both ends
			// to agree on the initial cap.
			plannedPosts := &query.SearchPosts{View: "planned", Limit: "10"}
			startedPosts := &query.SearchPosts{View: "started", Limit: "10"}
			completedPosts := &query.SearchPosts{View: "completed", Limit: "10"}
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
