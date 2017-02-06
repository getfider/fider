package context

import "github.com/WeCanHearYou/wchy/service"

// WchySettings is an application-wide settings
type WchySettings struct {
	BuildTime string
}

// WchyContext is an application-wide context
type WchyContext struct {
	Health   service.HealthCheckService
	Tenant   service.TenantService
	Settings WchySettings
}
