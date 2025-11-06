package handlers

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/web"
)

// ModerationPage is the moderation administration page (stub - requires commercial license)
func ModerationPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Failure(app.ErrCommercialLicenseRequired)
	}
}

// GetModerationItems returns all unmoderated posts and comments (stub - requires commercial license)
func GetModerationItems() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Failure(app.ErrCommercialLicenseRequired)
	}
}

// GetModerationCount returns the count of items awaiting moderation (stub - requires commercial license)
func GetModerationCount() web.HandlerFunc {
	return func(c *web.Context) error {
		// Return 0 instead of error to allow UI to function normally
		return c.Ok(web.Map{
			"count": 0,
		})
	}
}
