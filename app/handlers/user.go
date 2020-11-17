package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
	webutil "github.com/getfider/fider/app/pkg/web/util"
)

// BlockUser is used to block an existing user from using Fider
func BlockUser() web.HandlerFunc {
	return func(c *web.Context) error {
		userID, err := c.ParamAsInt("userID")
		if err != nil {
			return c.NotFound()
		}

		err = bus.Dispatch(c, &cmd.BlockUser{UserID: userID})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UnblockUser is used to unblock an existing user so they can use Fider again
func UnblockUser() web.HandlerFunc {
	return func(c *web.Context) error {
		userID, err := c.ParamAsInt("userID")
		if err != nil {
			return c.NotFound()
		}

		err = bus.Dispatch(c, &cmd.UnblockUser{UserID: userID})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// CreateUser is used to create a user with email and password
func CreateUser() web.HandlerFunc {
	return func(c *web.Context) error {
		if env.Config.SignUpDisabled {
			return c.NotFound()
		}
		input := new(actions.CreateUser)
		if result := c.BindToWithoutUser(input); !result.Ok {
			return c.HandleValidation(result)
		}

		user := &models.User{
			Tenant: c.Tenant(),
			Role:   enum.RoleCollaborator,
		}
		user.Email = input.Model.Email
		user.Password = input.Model.Password
		user.Name = input.Model.Name
		if err := bus.Dispatch(c, &cmd.RegisterUser{User: user}); err != nil {
			return c.Failure(err)
		}
		webutil.AddAuthUserCookie(c, user)

		return c.Ok(web.Map{})
	}
}
