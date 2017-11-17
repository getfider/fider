package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

// CreateEditTag creates a new tag on current tenant
func CreateEditTag() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.CreateEditTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		var (
			tag *models.Tag
			err error
		)

		if input.Model.Slug != "" {
			tag, err = c.Services().Tags.Update(input.Tag.ID, input.Model.Name, input.Model.Color, input.Model.IsPublic)
		} else {
			tag, err = c.Services().Tags.Add(input.Model.Name, input.Model.Color, input.Model.IsPublic)
		}

		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(tag)
	}
}
