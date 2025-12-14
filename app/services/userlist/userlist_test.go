package userlist_test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/httpclient/httpclientmock"
	"github.com/getfider/fider/app/services/userlist"
	userlist_mock "github.com/getfider/fider/app/services/userlist/mocks"

	. "github.com/getfider/fider/app/pkg/assert"
)

var ctx context.Context

func reset() {
	ctx = context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{
		Subdomain: "got",
	})
	bus.Init(userlist.Service{}, httpclientmock.Service{}, userlist_mock.Service{})
}

func TestCreateTenant_Success(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	createCompanyCmd := &cmd.UserListCreateCompany{
		Name:       "Fider",
		UserId:     1,
		UserEmail:  "jon.snow@got.com",
		UserName:   "Jon Snow",
		TenantId:   1,
		SignedUpAt: time.Now().Format(time.UnixDate),
		Plan:       enum.PlanFree,
		Subdomain:  "got",
	}

	err := bus.Dispatch(ctx, createCompanyCmd)
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(1)
	Expect(httpclientmock.RequestsHistory[0].URL.String()).Equals("https://push.userlist.com/companies")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Authorization")).Equals("Push abcdefg")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Content-Type")).Equals("application/json")
}

func TestUpdateTenant_Success(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	updateCompanyCmd := &cmd.UserListUpdateCompany{
		Name:     "Fider",
		TenantId: 1,
		Plan:     enum.PlanPro,
	}

	err := bus.Dispatch(ctx, updateCompanyCmd)
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(1)
	Expect(httpclientmock.RequestsHistory[0].URL.String()).Equals("https://push.userlist.com/companies")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Authorization")).Equals("Push abcdefg")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Content-Type")).Equals("application/json")
}

func TestUpdateTenant_PlanUpdatedIfSet(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	updateCompanyCmd := &cmd.UserListUpdateCompany{
		Name:     "Fider",
		TenantId: 1,
		Plan:     enum.PlanPro,
	}

	err := bus.Dispatch(ctx, updateCompanyCmd)
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(1)

	body, _ := io.ReadAll(httpclientmock.RequestsHistory[0].Body)
	containsPlan := strings.Contains(string(body), "plan")
	Expect(containsPlan).IsTrue()

	// Check we're sending the plan value
	Expect(strings.Contains(string(body), "plan\":\"pro\"")).IsTrue()

}

func TestUpdateTenant_PlanNotUpdatedIfNotSet(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	updateCompanyCmd := &cmd.UserListUpdateCompany{
		Name:     "Fider",
		TenantId: 1,
	}

	err := bus.Dispatch(ctx, updateCompanyCmd)
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(1)

	body, _ := io.ReadAll(httpclientmock.RequestsHistory[0].Body)
	containsPlan := strings.Contains(string(body), "plan")
	Expect(containsPlan).IsFalse()
}

func TestUpdateTenant_NameShouldUpdateIfSet(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	updateCompanyCmd := &cmd.UserListUpdateCompany{
		Name:     "Fider",
		TenantId: 1,
	}

	err := bus.Dispatch(ctx, updateCompanyCmd)
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(1)

	body, _ := io.ReadAll(httpclientmock.RequestsHistory[0].Body)
	containsName := strings.Contains(string(body), "name")
	Expect(containsName).IsTrue()
}

func TestUpdateTenant_NameShouldNotUpdateIfNotSet(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	updateCompanyCmd := &cmd.UserListUpdateCompany{
		TenantId: 1,
	}

	err := bus.Dispatch(ctx, updateCompanyCmd)
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(1)

	body, _ := io.ReadAll(httpclientmock.RequestsHistory[0].Body)
	containsName := strings.Contains(string(body), "name")
	Expect(containsName).IsFalse()
}

func TestUpdateUser_NameOnly(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	err := bus.Dispatch(ctx, &cmd.UserListUpdateUser{
		Id:       1,
		TenantId: 1,
		Name:     "Freddy",
	})
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(1)

	body, _ := io.ReadAll(httpclientmock.RequestsHistory[0].Body)
	containsName := strings.Contains(string(body), "name")
	Expect(containsName).IsTrue()

	containsEmail := strings.Contains(string(body), "email")
	Expect(containsEmail).IsFalse()
}

func TestUpdateUser_EmailOnly(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	err := bus.Dispatch(ctx, &cmd.UserListUpdateUser{
		Id:       1,
		TenantId: 1,
		Email:    "Freddy@example.com",
	})
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(1)

	body, _ := io.ReadAll(httpclientmock.RequestsHistory[0].Body)
	containsName := strings.Contains(string(body), "name")
	Expect(containsName).IsFalse()

	containsEmail := strings.Contains(string(body), "email")
	Expect(containsEmail).IsTrue()
}

func TestDoNothingIfUserNotAdmin(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	// Change the DB hit to return a visitor
	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		q.Result = &entity.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john.doe@example.com",
			Role:  enum.RoleVisitor,
		}
		return nil
	})

	err := bus.Dispatch(ctx, &cmd.UserListUpdateUser{
		Id:       1,
		TenantId: 1,
		Email:    "Freddy@example.com",
	})
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(0)

}

func TestMakeUserAdministrator(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	err := bus.Dispatch(ctx, &cmd.UserListHandleRoleChange{
		Id:   1,
		Role: enum.RoleAdministrator,
	})
	Expect(err).IsNil()

	Expect(httpclientmock.RequestsHistory).HasLen(1)

	Expect(httpclientmock.RequestsHistory[0].URL.String()).Equals("https://push.userlist.com/users")
	Expect(httpclientmock.RequestsHistory[0].Method).Equals("POST")

	body, _ := io.ReadAll(httpclientmock.RequestsHistory[0].Body)
	Expect(strings.Contains(string(body), "\"email\":\"john.doe@example.com\"")).IsTrue()
	Expect(strings.Contains(string(body), "\"name\":\"John Doe\"")).IsTrue()
	Expect(strings.Contains(string(body), "\"identifier\":\"1\"")).IsTrue()

}
