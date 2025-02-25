package handlers

import (
	"net/http"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/query"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
)

// ManageTags is the home page for managing tags
func ManageTags() web.HandlerFunc {
	return func(c *web.Context) error {
		getAllTags := &query.GetAllTags{}
		if err := bus.Dispatch(c, getAllTags); err != nil {
			return c.Failure(err)
		}

		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/ManageTags.page",
			Title: "Manage Tags · Site Settings",
			Data: web.Map{
				"tags": getAllTags.Result,
			},
		})
	}
}
