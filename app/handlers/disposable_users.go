package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// ListDisposableUsers returns up to 200 users whose email matches the
// disposable-blocking rules, plus the total match count.
func ListDisposableUsers() web.HandlerFunc {
	return func(c *web.Context) error {
		q := &query.GetDisposableUsers{Limit: 200}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}
		return c.Ok(web.Map{
			"total": q.Result.Total,
			"users": q.Result.Users,
		})
	}
}

// BulkDeleteDisposableUsers soft-deletes the listed users, but only those
// that still match the current deny rules at the moment of deletion.
func BulkDeleteDisposableUsers() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.BulkDeleteDisposable)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if len(action.UserIDs) == 0 {
			return c.Ok(web.Map{"deleted": 0})
		}

		// Re-verify each ID still matches current rules.
		current := &query.GetDisposableUsers{Limit: 1000000}
		if err := bus.Dispatch(c, current); err != nil {
			return c.Failure(err)
		}
		valid := make(map[int]struct{}, len(current.Result.Users))
		for _, u := range current.Result.Users {
			valid[u.UserID] = struct{}{}
		}
		filtered := make([]int, 0, len(action.UserIDs))
		for _, id := range action.UserIDs {
			if _, ok := valid[id]; ok {
				filtered = append(filtered, id)
			}
		}

		bulk := &cmd.BulkDeleteUsersByID{UserIDs: filtered}
		if err := bus.Dispatch(c, bulk); err != nil {
			return c.Failure(err)
		}
		return c.Ok(web.Map{"deleted": bulk.Result})
	}
}
