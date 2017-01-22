package context

import "github.com/WeCanHearYou/wchy-api/services"

// WchyContext is an application-wide context
type WchyContext struct {
	Health services.HealthCheckService
}
