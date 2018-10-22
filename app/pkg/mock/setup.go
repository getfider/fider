package mock

import (
	"os"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email/noop"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage/inmemory"
)

// DemoTenant is a mocked tenant
var DemoTenant *models.Tenant

// AvengersTenant is a mocked tenant
var AvengersTenant *models.Tenant

// JonSnow is a mocked user
var JonSnow *models.User

// AryaStark is a mocked user
var AryaStark *models.User

// NewSingleTenantServer creates a new multitenant test server
func NewSingleTenantServer() (*Server, *app.Services) {
	services := createServices(false)
	server := createServer(services)
	os.Setenv("HOST_MODE", "single")
	return server, services
}

// NewServer creates a new server and services for HTTP testing
func NewServer() (*Server, *app.Services) {
	services := createServices(true)
	server := createServer(services)
	os.Setenv("HOST_MODE", "multi")
	return server, services
}

// NewWorker creates a new worker and services for worker testing
func NewWorker() (*Worker, *app.Services) {
	services := createServices(true)
	worker := createWorker(services)
	return worker, services
}

func createServices(seed bool) *app.Services {
	services := &app.Services{
		Tenants:       inmemory.NewTenantStorage(),
		Users:         &inmemory.UserStorage{},
		Tags:          inmemory.NewTagStorage(),
		Notifications: inmemory.NewNotificationStorage(),
		Posts:         inmemory.NewPostStorage(),
		OAuth:         &OAuthService{},
		Emailer:       noop.NewSender(),
	}

	if seed {
		DemoTenant, _ = services.Tenants.Add("Demonstration", "demo", models.TenantActive)
		AvengersTenant, _ = services.Tenants.Add("Avengers", "avengers", models.TenantActive)
		AvengersTenant.CNAME = "feedback.theavengers.com"

		JonSnow = &models.User{
			Name:   "Jon Snow",
			Email:  "jon.snow@got.com",
			Tenant: DemoTenant,
			Status: models.UserActive,
			Role:   models.RoleAdministrator,
			Providers: []*models.UserProvider{
				{UID: "FB1234", Name: oauth.FacebookProvider},
			},
		}
		services.Users.Register(JonSnow)

		AryaStark = &models.User{
			Name:   "Arya Stark",
			Email:  "arya.stark@got.com",
			Tenant: DemoTenant,
			Status: models.UserActive,
			Role:   models.RoleVisitor,
		}
		services.Users.Register(AryaStark)
	}

	return services
}
