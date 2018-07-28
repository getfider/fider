package apiv1

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/web"
)

// AssignTag to existing dea
func AssignTag() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.AssignUnassignTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tags.AssignTag(input.Tag, input.Post)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UnassignTag from existing dea
func UnassignTag() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.AssignUnassignTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tags.UnassignTag(input.Tag, input.Post)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
