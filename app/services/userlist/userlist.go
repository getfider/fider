// Company represents a company in UserList.com
package userlist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

// Company represents a company in UserList.com. It contains information about the company,
// including its identifier, name, sign-up date, custom properties, and associated users.
type Company struct {
	Identifier string                 `json:"identifier"`
	Name       string                 `json:"name,omitempty"`
	SignedUpAt string                 `json:"signed_up_at,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Users      []UserListUser         `json:"users,omitempty"`
}

// UserListUser represents a user in UserList.com. It contains information about the user,
// including their identifier, email, sign-up date, and custom properties.
type UserListUser struct {
	Identifier string                 `json:"identifier,omitempty"`
	Company    string                 `json:"company,omitempty"`
	Email      string                 `json:"email,omitempty"`
	SignedUpAt string                 `json:"signed_up_at,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type UserListUserToDelete struct {
	Identifier string `json:"identifier"`
}

type CompanyCreateResponse struct {
	Identifier string `json:"id"`
}

const (
	UserListBaseUrl = "https://push.userlist.com"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "UserList"
}

func (s Service) Category() string {
	return "userlist"
}

func (s Service) Enabled() bool {
	return env.Config.UserList.Enabled
}

func (s Service) Init() {
	bus.AddHandler(createUserListCompany)
	bus.AddHandler(updateUserListCompany)
	bus.AddHandler(updateUserListUser)
	bus.AddHandler(addOrRemoveUserListUser)
}

func addOrRemoveUserListUser(ctx context.Context, u *cmd.UserListHandleRoleChange) error {
	if u.Role == enum.RoleAdministrator {
		// Get the user so we can add it to userlist.
		user := &query.GetUserByID{
			UserID: u.Id,
		}
		err := bus.Dispatch(ctx, user)
		if err != nil {
			return err
		}

		err = updateUserListUser(ctx, &cmd.UserListUpdateUser{
			Id:       u.Id,
			TenantId: user.Result.Tenant.ID,
			Email:    user.Result.Email,
			Name:     user.Result.Name,
		})
		if err != nil {
			return err
		}
		return nil
	}

	// If the user is not an administrator, we remove it from userlist.
	userToDelete := &UserListUserToDelete{
		Identifier: strconv.Itoa(u.Id),
	}

	err := pushUserListUpdate(userToDelete, ctx)
	if err != nil {
		return err
	}
	return nil

}

func updateUserListUser(ctx context.Context, u *cmd.UserListUpdateUser) error {

	user := &UserListUser{
		Identifier: strconv.Itoa(u.Id),
		Company:    strconv.Itoa(u.TenantId),
	}

	if len(u.Email) > 0 {
		user.Email = u.Email
	}

	if len(u.Name) > 0 {
		user.Properties = map[string]interface{}{
			"name": u.Name,
		}
	}

	err := pushUserListUpdate(user, ctx)
	if err != nil {
		return err
	}
	return nil

}

func updateUserListCompany(ctx context.Context, c *cmd.UserListUpdateCompany) error {
	company := &Company{
		Identifier: strconv.Itoa(c.TenantId),
		Name:       c.Name,
	}

	if c.BillingStatus > 0 {
		company.Properties = map[string]interface{}{
			"billing_status": c.BillingStatus.String(),
		}
	}

	err := pushUserListUpdate(company, ctx)
	if err != nil {
		return err
	}
	return nil

}

func createUserListCompany(ctx context.Context, c *cmd.UserListCreateCompany) error {

	company := &Company{
		Identifier: strconv.Itoa(c.TenantId),
		Name:       c.Name,
		SignedUpAt: c.SignedUpAt,
		Properties: map[string]interface{}{
			"billing_status": c.BillingStatus,
			"subdomain":      c.Subdomain,
		},
		Users: []UserListUser{
			{
				Identifier: strconv.Itoa(c.UserId),
				Email:      c.UserEmail,
				SignedUpAt: c.SignedUpAt,
				Properties: map[string]interface{}{
					"name": c.UserName,
				},
			},
		},
	}

	err := pushUserListUpdate(company, ctx)
	if err != nil {
		return err
	}
	return nil

}

func pushUserListUpdate(obj interface{}, ctx context.Context) error {
	var url string
	method := http.MethodPost
	switch obj.(type) {
	case *Company:
		url = fmt.Sprintf("%s/companies", UserListBaseUrl)
	case *UserListUser:
		url = fmt.Sprintf("%s/users", UserListBaseUrl)
	case *UserListUserToDelete:
		url = fmt.Sprintf("%s/users", UserListBaseUrl)
		method = http.MethodDelete
	default:
		return errors.New("invalid type passed to pushUserListUpdate")
	}

	jsonContent, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "failed to marshal object")
	}

	req := &cmd.HTTPRequest{
		URL:    url,
		Body:   bytes.NewBuffer(jsonContent),
		Method: method,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Push %s", env.Config.UserList.ApiKey),
		},
	}

	if err := bus.Dispatch(ctx, req); err != nil {
		return errors.Wrap(err, "Failed to send userlist update")
	}

	if req.ResponseStatusCode >= 300 {
		return errors.New("unexpected status code while updating userlist: %d", req.ResponseStatusCode)
	}

	return nil
}
