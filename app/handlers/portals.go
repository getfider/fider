package handlers

import (
	"net/http"
	"strings"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// PortalDirectory lists the public portals hosted on this instance. It is only reachable on
// the root domain of a multi-tenant instance, and only when the operator has enabled it.
func PortalDirectory() web.HandlerFunc {
	return func(c *web.Context) error {
		if !env.IsPortalDirectoryEnabled() {
			return c.NotFound()
		}

		publicTenants := &query.GetPublicTenants{}
		if err := bus.Dispatch(c, publicTenants); err != nil {
			return c.Failure(err)
		}

		portals := make([]dto.PortalSummary, 0, len(publicTenants.Result))
		for _, tenant := range publicTenants.Result {
			url := web.TenantBaseURL(c, tenant)
			portals = append(portals, dto.PortalSummary{
				Name:    tenant.Name,
				URL:     url,
				Host:    hostOf(url),
				LogoURL: web.TenantLogoURL(c, tenant),
			})
		}

		return c.Page(http.StatusOK, web.Props{
			Page:        "PortalDirectory/PortalDirectory.page",
			Title:       "Portals",
			Description: "Browse the feedback portals hosted here.",
			Data: web.Map{
				"portals": portals,
			},
		})
	}
}

// hostOf strips the scheme from a base URL, leaving the host (and port, when present).
func hostOf(url string) string {
	if index := strings.Index(url, "://"); index >= 0 {
		return url[index+3:]
	}
	return url
}
