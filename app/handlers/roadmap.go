package handlers

import (
	"net/http"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// roadmapColumn is one lane on the board: a tenant status plus the posts
// currently sitting in that status. Marshalled to JSON for the React client.
type roadmapColumn struct {
	Status *entity.Status `json:"status"`
	Posts  []*entity.Post `json:"posts"`
}

// RoadmapPage renders the roadmap board. Self-hosted installations and Pro
// tenants get the full data; everyone else renders the page empty so the
// client shows the upgrade call-to-action.
//
// Lanes are derived from tenant.statuses (feedback.fider.io/posts/111) —
// every active, show-on-roadmap status becomes its own column, ordered by
// sort_order. Admins flip the per-status "Show on roadmap" toggle in the
// Manage Statuses page to publish any status (built-in or custom) here.
func RoadmapPage() web.HandlerFunc {
	return func(c *web.Context) error {
		props := web.Props{
			Page:  "Roadmap/Roadmap.page",
			Title: "Roadmap",
		}

		if env.IsSingleHostMode() || c.Tenant().IsPro {
			tenantStatuses := c.Tenant().Statuses
			columns := make([]roadmapColumn, 0, len(tenantStatuses))

			for _, s := range tenantStatuses {
				if !s.IsActive || !s.ShowOnRoadmap {
					continue
				}

				// Keep in sync with ROADMAP_DEFAULT_LIMIT on the client; the
				// "Show more" link uses posts.length >= limit as its
				// heuristic and needs both ends to agree on the initial cap.
				posts := &query.SearchPosts{Statuses: []string{s.Slug}, Limit: "10"}
				if err := bus.Dispatch(c, posts); err != nil {
					return c.Failure(err)
				}

				columns = append(columns, roadmapColumn{Status: s, Posts: posts.Result})
			}

			getAllTags := &query.GetAllTags{}
			if err := bus.Dispatch(c, getAllTags); err != nil {
				return c.Failure(err)
			}

			props.Data = web.Map{
				"columns": columns,
				"tags":    getAllTags.Result,
			}
		}

		return c.Page(http.StatusOK, props)
	}
}
