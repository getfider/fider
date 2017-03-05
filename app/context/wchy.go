package context

import (
	"github.com/WeCanHearYou/wchy/app/auth"
	"github.com/WeCanHearYou/wchy/app/service"
)

// WchySettings is an application-wide settings
type WchySettings struct {
	BuildTime    string
	AuthEndpoint string
}

// WchyContext is an application-wide context
type WchyContext struct {
	OAuth    auth.OAuthService
	Health   service.HealthCheckService
	User     service.UserService
	Idea     service.IdeaService
	Tenant   service.TenantService
	Settings WchySettings
}
