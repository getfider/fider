package handlers_test

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage/inmemory"
)

var demoTenant *models.Tenant
var orangeTenant *models.Tenant

var jonSnow *models.User
var aryaStark *models.User

func getEmptyServices() *app.Services {
	return &app.Services{
		Tenants: &inmemory.TenantStorage{},
		Users:   &inmemory.UserStorage{},
		OAuth:   &mock.OAuthService{},
	}
}

func getServices() *app.Services {
	services := &app.Services{
		Tenants: &inmemory.TenantStorage{},
		Users:   &inmemory.UserStorage{},
		OAuth:   &mock.OAuthService{},
	}

	demoTenant, _ = services.Tenants.Add("Demonstration", "demo")
	orangeTenant, _ = services.Tenants.Add("Orange Inc.", "orange")

	return services
}

func setupServer() (*mock.Server, *app.Services) {

	tenants := &inmemory.TenantStorage{}
	demoTenant, _ = tenants.Add("Demonstration", "demo")
	orangeTenant, _ = tenants.Add("Orange Inc.", "orange")

	services := &app.Services{
		Tenants: tenants,
		Users:   &inmemory.UserStorage{},
		Ideas:   &inmemory.IdeaStorage{},
		OAuth:   &mock.OAuthService{},
	}

	jonSnow = &models.User{
		Name:   "Jon Snow",
		Email:  "jon.snow@got.com",
		Tenant: demoTenant,
		Role:   models.RoleAdministrator,
		Providers: []*models.UserProvider{
			{UID: "FB1234", Name: oauth.FacebookProvider},
		},
	}
	services.Users.Register(jonSnow)

	aryaStark = &models.User{Name: "Arya Stark", Email: "arya.stark@got.com", Tenant: demoTenant, Role: models.RoleVisitor}
	services.Users.Register(aryaStark)

	server := mock.NewServer()
	server.Context.SetServices(services)

	return server, services
}
