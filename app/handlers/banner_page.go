package handlers

import (
	"net/http"

	"github.com/getfider/fider/app/pkg/web"
)

// ManageBanner renders the admin page where tenants configure the
// site-wide banner shown above the page header.
func ManageBanner() web.HandlerFunc {
	return func(c *web.Context) error {
		tenant := c.Tenant()
		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/ManageBanner.page",
			Title: "Site Banner · Site Settings",
			Data: web.Map{
				"siteBannerEnabled": tenant.SiteBannerEnabled,
				"siteBannerMessage": tenant.SiteBannerMessage,
				"siteBannerVariant": tenant.SiteBannerVariant,
			},
		})
	}
}
