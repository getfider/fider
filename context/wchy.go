package context

import (
	"github.com/WeCanHearYou/wchy/auth"
	"github.com/WeCanHearYou/wchy/service"
)

// WchySettings is an application-wide settings
type WchySettings struct {
	BuildTime    string
	AuthEndpoint string
}

// WchyContext is an application-wide context
type WchyContext struct {
	Health   service.HealthCheckService
	Auth     auth.Service
	Idea     service.IdeaService
	Tenant   service.TenantService
	Settings WchySettings
}
