package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/web"
)

// CreateTag creates a new tag on current tenant
func CreateTag() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.CreateNewTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		tag, err := c.Services().Tags.Add(input.Model.Name, input.Model.Color, input.Model.IsPublic)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(tag)
	}
}
