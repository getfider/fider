package handlers_test

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage/inmemory"
)

func getServices() *app.Services {
	return &app.Services{
		Tenants: &inmemory.TenantStorage{},
		Users:   &inmemory.UserStorage{},
		OAuth:   &oauth.MockOAuthService{},
	}
}
