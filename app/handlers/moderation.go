package handlers

import (
	"net/http"

	"github.com/getfider/fider/app/pkg/web"
)

// ModerationPage is the moderation administration page (stub - requires commercial license)
func ModerationPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/ContentModeration.page",
			Title: "Moderation · Site Settings",
		})
	}
}

// GetModerationItems returns all unmoderated posts and comments (stub - requires commercial license)
func GetModerationItems() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Ok(web.Map{
			"items": []interface{}{},
		})
	}
}

// GetModerationCount returns the count of items awaiting moderation (stub - requires commercial license)
func GetModerationCount() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Ok(web.Map{
			"count": 0,
		})
	}
}
