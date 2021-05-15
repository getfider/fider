package apiv1

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// ListTags returns all tags
func ListTags() web.HandlerFunc {
	return func(c *web.Context) error {
		q := &query.GetAllTags{}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}

		return c.Ok(q.Result)
	}
}

// AssignTag to existing dea
func AssignTag() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.AssignUnassignTag)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.AssignTag{Tag: action.Tag, Post: action.Post})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UnassignTag from existing dea
func UnassignTag() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.AssignUnassignTag)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.UnassignTag{Tag: action.Tag, Post: action.Post})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// CreateEditTag creates a new tag on current tenant
func CreateEditTag() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.CreateEditTag)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if action.Slug != "" {
			updateTag := &cmd.UpdateTag{
				TagID:    action.Tag.ID,
				Name:     action.Name,
				Color:    action.Color,
				IsPublic: action.IsPublic,
			}
			if err := bus.Dispatch(c, updateTag); err != nil {
				return c.Failure(err)
			}
			return c.Ok(updateTag.Result)
		}

		addNewTag := &cmd.AddNewTag{
			Name:     action.Name,
			Color:    action.Color,
			IsPublic: action.IsPublic,
		}
		if err := bus.Dispatch(c, addNewTag); err != nil {
			return c.Failure(err)
		}
		return c.Ok(addNewTag.Result)
	}
}

// DeleteTag deletes an existing tag
func DeleteTag() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.DeleteTag)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.DeleteTag{Tag: action.Tag})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
