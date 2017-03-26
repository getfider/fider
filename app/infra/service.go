package infra

// HealthCheckService is a simple general purpose health check service
type HealthCheckService interface {
	IsDatabaseOnline() bool
}
