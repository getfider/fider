package userlist_test

import (
	"context"
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/httpclient/httpclientmock"
	"github.com/getfider/fider/app/services/userlist"

	. "github.com/getfider/fider/app/pkg/assert"
)

var ctx context.Context

func reset() {
	ctx = context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{
		Subdomain: "got",
	})
	bus.Init(userlist.Service{}, httpclientmock.Service{})
}

func TestCreatTenant_Success(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	createCompanyCmd := &cmd.UserListCreateCompany{
		Name:          "Fider",
		UserId:        1,
		UserEmail:     "jon.snow@got.com",
		UserName:      "Jon Snow",
		TenantId:      1,
		SignedUpAt:    time.Now().Format(time.UnixDate),
		BillingStatus: "active",
		Subdomain:     "got",
	}

	bus.Dispatch(ctx, createCompanyCmd)

	Expect(httpclientmock.RequestsHistory).HasLen(1)
	Expect(httpclientmock.RequestsHistory[0].URL.String()).Equals("https://push.userlist.com/companies")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Authorization")).Equals("Push abcdefg")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Content-Type")).Equals("application/json")
}

func TestUdpateTenant_Success(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	updateCompanyCmd := &cmd.UserListUpdateCompany{
		Name:          "Fider",
		Id:            1,
		BillingStatus: "active",
		Subdomain:     "got",
	}

	bus.Dispatch(ctx, updateCompanyCmd)

	Expect(httpclientmock.RequestsHistory).HasLen(1)
	Expect(httpclientmock.RequestsHistory[0].URL.String()).Equals("https://push.userlist.com/companies")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Authorization")).Equals("Push abcdefg")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Content-Type")).Equals("application/json")
}
