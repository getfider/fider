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
	Email      string                 `json:"email,omitempty"`
	SignedUpAt string                 `json:"signed_up_at,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
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
}

func createUserListCompany(ctx context.Context, c *cmd.CreateUserListCompany) error {

	company := &Company{
		Identifier: strconv.Itoa(c.TenantId),
		Name:       c.Name,
		SignedUpAt: c.SignedUpAt,
		Properties: map[string]interface{}{
			"billing_status": c.BillingStatus,
			"user_count":     1,
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

	jsonContent, err := json.Marshal(company)
	if err != nil {
		return errors.Wrap(err, "failed to marshal company object")
	}

	req := &cmd.HTTPRequest{
		URL:    fmt.Sprintf("%s/companies", UserListBaseUrl),
		Body:   bytes.NewBuffer(jsonContent),
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Push %s", env.Config.UserList.ApiKey),
		},
	}

	if err := bus.Dispatch(ctx, req); err != nil {
		return errors.Wrap(err, "Failed to create userlist company")
	}

	if req.ResponseStatusCode >= 300 {
		return errors.New("unexpected status code while creating company in userlist: %d", req.ResponseStatusCode)
	}

	return nil

}
