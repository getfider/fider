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
		input := new(actions.AssignUnassignTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.AssignTag{Tag: input.Tag, Post: input.Post})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UnassignTag from existing dea
func UnassignTag() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.AssignUnassignTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.UnassignTag{Tag: input.Tag, Post: input.Post})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// CreateEditTag creates a new tag on current tenant
func CreateEditTag() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.CreateEditTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		if input.Model.Slug != "" {
			updateTag := &cmd.UpdateTag{
				TagID:    input.Tag.ID,
				Name:     input.Model.Name,
				Color:    input.Model.Color,
				IsPublic: input.Model.IsPublic,
			}
			if err := bus.Dispatch(c, updateTag); err != nil {
				return c.Failure(err)
			}
			return c.Ok(updateTag.Result)
		}

		addNewTag := &cmd.AddNewTag{
			Name:     input.Model.Name,
			Color:    input.Model.Color,
			IsPublic: input.Model.IsPublic,
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
		input := new(actions.DeleteTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.DeleteTag{Tag: input.Tag})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
