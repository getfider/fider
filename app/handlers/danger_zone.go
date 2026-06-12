package handlers

import (
	"net/http"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// gracePeriod is how long a site deletion is delayed so the owner can change their mind.
const gracePeriod = time.Hour

// tenantOwner returns the account owner (the lowest-id active administrator) of the current
// tenant. Callers compare the returned user's ID with c.User().ID to gate owner-only actions.
func tenantOwner(c *web.Context) (*entity.User, error) {
	q := &query.GetTenantOwner{TenantID: c.Tenant().ID}
	if err := bus.Dispatch(c, q); err != nil {
		return nil, err
	}
	return q.Result, nil
}

// DangerZonePage renders the admin Danger Zone (delete entire site).
func DangerZonePage() web.HandlerFunc {
	return func(c *web.Context) error {
		owner, err := tenantOwner(c)
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/DangerZone.page",
			Title: "Danger Zone · Site Settings",
			Data: web.Map{
				"isOwner":             c.User().ID == owner.ID,
				"scheduledDeletionAt": c.Tenant().ScheduledDeletionAt,
			},
		})
	}
}

// RequestTenantDeletion schedules deletion of the whole site after a grace period. Only the
// account owner may do this, and only on hosted multi-tenant instances.
func RequestTenantDeletion() web.HandlerFunc {
	return func(c *web.Context) error {
		if env.IsSingleHostMode() {
			return c.Forbidden()
		}

		owner, err := tenantOwner(c)
		if err != nil {
			return c.Failure(err)
		}
		if c.User().ID != owner.ID {
			return c.Forbidden()
		}

		input := new(struct {
			Subdomain string `json:"subdomain"`
		})
		if err := c.Bind(input); err != nil {
			return c.BadRequest(web.Map{})
		}
		if input.Subdomain != c.Tenant().Subdomain {
			return c.BadRequest(web.Map{"message": "The subdomain you entered does not match this site."})
		}

		cancelKey := rand.String(64)
		scheduledAt := time.Now().Add(gracePeriod)

		if err := bus.Dispatch(c, &cmd.ScheduleTenantDeletion{
			TenantID:          c.Tenant().ID,
			RequestedByUserID: c.User().ID,
			CancelKey:         cancelKey,
			ScheduledAt:       scheduledAt,
		}); err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.SendDeleteAccountScheduledEmail(owner, c.Tenant().Name, scheduledAt, c.BaseURL(), cancelKey))

		return c.Ok(web.Map{"scheduledDeletionAt": scheduledAt})
	}
}

// CancelTenantDeletionByOwner cancels a scheduled deletion from the Danger Zone page. Unlike
// the emailed link, this is authorised by the logged-in owner rather than a key.
func CancelTenantDeletionByOwner() web.HandlerFunc {
	return func(c *web.Context) error {
		if env.IsSingleHostMode() {
			return c.Forbidden()
		}

		owner, err := tenantOwner(c)
		if err != nil {
			return c.Failure(err)
		}
		if c.User().ID != owner.ID {
			return c.Forbidden()
		}

		if err := bus.Dispatch(c, &cmd.CancelTenantDeletion{TenantID: c.Tenant().ID}); err != nil {
			return c.Failure(err)
		}
		return c.Ok(web.Map{})
	}
}

// CancelTenantDeletion cancels a scheduled deletion. The unguessable key in the query string
// is the sole authorisation, so the emailed cancel link works without the owner being logged
// in (it only ever restores access).
func CancelTenantDeletion() web.HandlerFunc {
	return func(c *web.Context) error {
		key := c.QueryParam("k")
		if key != "" {
			byKey := &query.GetTenantByCancelKey{Key: key}
			if err := bus.Dispatch(c, byKey); err == nil && byKey.Result != nil {
				if err := bus.Dispatch(c, &cmd.CancelTenantDeletion{TenantID: byKey.Result.ID}); err != nil {
					return c.Failure(err)
				}
			}
		}
		return c.Redirect(c.BaseURL() + "/admin/danger-zone")
	}
}
