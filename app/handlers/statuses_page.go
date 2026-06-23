package handlers

import (
	"net/http"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// ManageStatuses renders the admin page where tenants configure their
// post status catalogue.
func ManageStatuses() web.HandlerFunc {
	return func(c *web.Context) error {
		listStatuses := &query.ListActiveStatusesForTenant{}
		if err := bus.Dispatch(c, listStatuses); err != nil {
			return c.Failure(err)
		}

		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/ManageStatuses.page",
			Title: "Manage Statuses · Site Settings",
			Data: web.Map{
				"statuses": listStatuses.Result,
			},
		})
	}
}
