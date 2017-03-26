package postgres

import "database/sql"

// HealthCheckService is a simple general purpose health check service
type HealthCheckService struct {
	DB *sql.DB
}

// IsDatabaseOnline checks if database is online
func (svc HealthCheckService) IsDatabaseOnline() bool {
	_, err := svc.DB.Query("SELECT now()")
	if err != nil {
		return false
	}

	return true
}
