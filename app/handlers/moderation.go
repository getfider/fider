package handlers

import (
	"net/http"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// ModerationPage is the moderation administration page
func ModerationPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/ContentModeration.page",
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

// GetModerationCount returns the count of items awaiting moderation
func GetModerationCount() web.HandlerFunc {
	return func(c *web.Context) error {
		q := &query.GetModerationCount{}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"count": q.Result,
		})
	}
}
