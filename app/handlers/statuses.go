package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/pkg/web"
)

// ListStatuses returns the active status catalogue for the current tenant.
// Used both by the admin UI (Manage Statuses) and by the bootstrap path that
// hydrates the React client's PostStatus.All at page load.
func ListStatuses() web.HandlerFunc {
	return func(c *web.Context) error {
		q := &query.ListActiveStatusesForTenant{}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}
		return c.Ok(q.Result)
	}
}

// CreateStatus inserts an admin-defined custom status for the current tenant.
func CreateStatus() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.CreateStatus)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		create := &cmd.CreateStatus{
			Slug:          action.Slug,
			Label:         action.Label,
			Kind:          action.Kind,
			Color:         action.Color,
			Icon:          action.Icon,
			ShowOnHome:    action.ShowOnHome,
			ShowOnRoadmap: action.ShowOnRoadmap,
			Filterable:    action.Filterable,
			SortOrder:     action.SortOrder,
		}
		if err := bus.Dispatch(c, create); err != nil {
			return c.Failure(err)
		}
		return c.Ok(create.Result)
	}
}

// UpdateStatus mutates a status the current tenant owns.
func UpdateStatus() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.UpdateStatus)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		update := &cmd.UpdateStatus{
			ID:            action.ID,
			Label:         action.Label,
			Color:         action.Color,
			Icon:          action.Icon,
			ShowOnHome:    action.ShowOnHome,
			ShowOnRoadmap: action.ShowOnRoadmap,
			Filterable:    action.Filterable,
			SortOrder:     action.SortOrder,
			IsActive:      action.IsActive,
		}
		if err := bus.Dispatch(c, update); err != nil {
			return c.Failure(err)
		}
		return c.Ok(web.Map{})
	}
}

// DeleteStatus removes a non-system status. Refused if any post still uses it.
// Note: bypass BindTo for the action because DELETE-with-no-body has tripped
// content-type validation here in the past; we parse :id directly and run
// the IsAuthorized check inline.
func DeleteStatus() web.HandlerFunc {
	return func(c *web.Context) error {
		id, err := c.ParamAsInt("id")
		if err != nil {
			r := validate.Success()
			r.AddFieldFailure("id", "Status id must be an integer.")
			return c.HandleValidation(r)
		}
		if c.User() == nil || c.User().Role != enum.RoleAdministrator {
			return c.Forbidden()
		}

		count := &query.CountPostsByStatus{StatusID: id}
		if err := bus.Dispatch(c, count); err != nil {
			return c.Failure(err)
		}
		if count.Result > 0 {
			r := validate.Success()
			r.AddFieldFailure("status", "Cannot delete: posts are still using this status. Reassign them first.")
			return c.HandleValidation(r)
		}

		if err := bus.Dispatch(c, &cmd.DeleteStatus{ID: id}); err != nil {
			r := validate.Success()
			r.AddFieldFailure("status", err.Error())
			return c.HandleValidation(r)
		}
		return c.Ok(web.Map{})
	}
}
