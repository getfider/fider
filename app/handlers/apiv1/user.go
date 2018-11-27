package apiv1

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"
)

// ListUsers returns all registered users
func ListUsers() web.HandlerFunc {
	return func(c web.Context) error {
		users, err := c.Services().Users.GetAll()
		if err != nil {
			return c.Failure(err)
		}
		return c.Ok(users)
	}
}

// CreateUser is used to create new users
func CreateUser() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.CreateUser)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		user, err := c.Services().Users.GetByProvider("reference", input.Model.Reference)
		if err != nil && errors.Cause(err) == app.ErrNotFound {
			if input.Model.Email != "" {
				user, err = c.Services().Users.GetByEmail(input.Model.Email)
			}
			if err != nil && errors.Cause(err) == app.ErrNotFound {
				user = &models.User{
					Tenant: c.Tenant(),
					Name:   input.Model.Name,
					Email:  input.Model.Email,
				}
				err = c.Services().Users.Register(user)
			}
		}

		if err != nil {
			return c.Failure(err)
		}

		if input.Model.Reference != "" && !user.HasProvider("reference") {
			if err := c.Services().Users.RegisterProvider(user.ID, &models.UserProvider{
				Name: "reference",
				UID:  input.Model.Reference,
			}); err != nil {
				return c.Failure(err)
			}
		}

		return c.Ok(web.Map{
			"id": user.ID,
		})
	}
}
