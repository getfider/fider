package services

// HealthCheckService a simple general purpose health check service
type HealthCheckService interface {
	IsDatabaseOnline() bool
}
