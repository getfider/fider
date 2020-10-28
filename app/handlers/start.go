package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/passhash"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"
)

//CreateTenant creates a new tenant
func CreateTenant() web.HandlerFunc {
	return func(c *web.Context) error {
		if env.Config.SignUpDisabled {
			return c.NotFound()
		}

		input := new(actions.CreateTenant)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		socialSignUp := input.Model.Token != ""

		status := enum.TenantActive
		if socialSignUp {
			status = enum.TenantActive
		}

		createTenant := &cmd.CreateTenant{
			Name:      input.Model.TenantName,
			Subdomain: input.Model.Subdomain,
			Status:    status,
		}
		err := bus.Dispatch(c, createTenant)
		if err != nil {
			return c.Failure(err)
		}

		c.SetTenant(createTenant.Result)

		user := &models.User{
			Tenant: createTenant.Result,
			Role:   enum.RoleAdministrator,
		}

		user.Name = input.Model.Name
		user.Email = input.Model.Email
		user.Password, err = passhash.HashString(input.Model.Password)
		if err != nil {
			return c.Failure(err)
		}
		if err := bus.Dispatch(c, &cmd.RegisterUser{User: user}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

//Start is the entry point for installation /start
func Start() web.HandlerFunc {
	return func(c *web.Context) error {
		if env.Config.SignUpDisabled {
			return c.NotFound()
		}

		if env.IsSingleHostMode() {
			firstTenant := &query.GetFirstTenant{}
			err := bus.Dispatch(c, firstTenant)
			if err != nil && errors.Cause(err) != app.ErrNotFound {
				return c.Failure(err)
			}

			if firstTenant.Result != nil {
				return c.Redirect("/")
			}
		}
		return c.Page(web.Props{
			Title:       "Start",
			Description: "Start for Fider and let your customers share, vote and discuss on suggestions they have to make your product even better.",
			ChunkName:   "Start.page",
		})
	}
}
