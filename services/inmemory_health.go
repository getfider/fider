package services

// InMemoryHealthCheckService checks for database health status
type InMemoryHealthCheckService struct {
	status bool
}

// NewInMemoryHealthCheckService creates a new InMemoryHealthCheckService
func NewInMemoryHealthCheckService(status bool) *InMemoryHealthCheckService {
	return &InMemoryHealthCheckService{status: status}
}

// IsDatabaseOnline checks if database is online
func (svc InMemoryHealthCheckService) IsDatabaseOnline() bool {
	return svc.status
}
