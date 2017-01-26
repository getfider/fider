package services

import (
	"database/sql"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("services/health")

// PostgresHealthCheckService is a simple general purpose health check service
type PostgresHealthCheckService struct {
	db *sql.DB
}

// NewPostgresHealthCheckService creates a new PostgresHealthCheckService
func NewPostgresHealthCheckService(db *sql.DB) *PostgresHealthCheckService {
	return &PostgresHealthCheckService{db}
}

// IsDatabaseOnline checks if database is online
func (svc PostgresHealthCheckService) IsDatabaseOnline() bool {
	_, err := svc.db.Query("SELECT now()")
	if err != nil {
		log.Error(err)
		return false
	}

	return true
}
