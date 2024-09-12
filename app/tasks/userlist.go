package tasks

import (
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/worker"
)

func UserListCreateCompany(tenant entity.Tenant, user entity.User) worker.Task {
	return describe("Create UserList Company", func(c *worker.Context) error {
		log.Debugf(c, "Sending new tenant @{Tenant} to userlist with user email @{User}", dto.Props{
			"Tenant": tenant.Name,
			"User":   user.Email,
		})
		if err := bus.Dispatch(c, &cmd.UserListCreateCompany{
			Name:          tenant.Name,
			TenantId:      tenant.ID,
			SignedUpAt:    time.Now().Format(time.RFC3339),
			BillingStatus: enum.BillingTrial.String(),
			Subdomain:     tenant.Subdomain,
			UserId:        user.ID,
			UserEmail:     user.Email,
			UserName:      user.Name,
		}); err != nil {
			return c.Failure(err)
		}
		return nil
	})
}

func UserListUpdateCompany(action *dto.UserListUpdateCompany) worker.Task {
	return describe("Update Company in UserList", func(c *worker.Context) error {
		log.Debugf(c, "Updating company @{Tenant} in UserList", dto.Props{
			"Tenant": action.Name,
		})
		if err := bus.Dispatch(c, &cmd.UserListUpdateCompany{
			TenantId:      action.TenantID,
			Name:          action.Name,
			BillingStatus: action.BillingStatus,
		}); err != nil {
			return c.Failure(err)
		}
		return nil
	})
}

func UserListUpdateUser(id int, name string, email string) worker.Task {
	return describe("Update User in UserList", func(c *worker.Context) error {
		log.Debugf(c, "Updating user @{User} in UserList", dto.Props{
			"User": id,
		})
		if err := bus.Dispatch(c, &cmd.UserListUpdateUser{
			Id:    id,
			Email: email,
			Name:  name,
		}); err != nil {
			return c.Failure(err)
		}
		return nil
	})
}

func UserListAddOrRemoveUser(userID int, role enum.Role) worker.Task {
	return describe("Add or Remove User in UserList", func(c *worker.Context) error {
		log.Debugf(c, "Handling role change for user in UserList", dto.Props{
			"User": userID,
		})
		if err := bus.Dispatch(c, &cmd.UserListHandleRoleChange{
			Id:   userID,
			Role: role,
		}); err != nil {
			return c.Failure(err)
		}
		return nil
	})
}
