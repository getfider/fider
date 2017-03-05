package service

import "database/sql"

// HealthCheckService is a simple general purpose health check service
type HealthCheckService interface {
	IsDatabaseOnline() bool
}

// PostgresHealthCheckService is a simple general purpose health check service
type PostgresHealthCheckService struct {
	DB *sql.DB
}

// IsDatabaseOnline checks if database is online
func (svc PostgresHealthCheckService) IsDatabaseOnline() bool {
	_, err := svc.DB.Query("SELECT now()")
	if err != nil {
		return false
	}

	return true
}
