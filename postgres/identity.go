package postgres

import (
	"database/sql"
	"strings"

	"github.com/WeCanHearYou/wechy/identity"
	"github.com/WeCanHearYou/wechy/toolbox/env"
)

// UserService is used for user operations using a Postgres database
type UserService struct {
	DB *sql.DB
}

// GetByEmail returns a user based on given email
func (svc UserService) GetByEmail(email string) (*identity.User, error) {
	user := &identity.User{}
	row := svc.DB.QueryRow("SELECT id, name, email FROM users WHERE email = $1", email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, identity.ErrNotFound
	}

	return user, nil
}

// Register creates a new user based on given information
func (svc UserService) Register(user *identity.User) error {
	tx, err := svc.DB.Begin()
	if err != nil {
		return err
	}

	if err = tx.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&user.ID); err != nil {
		tx.Rollback()
		return err
	}

	for _, provider := range user.Providers {
		if _, err = tx.Exec("INSERT INTO user_providers (user_id, provider, provider_uid) VALUES ($1, $2, $3)", user.ID, provider.Name, provider.UID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// TenantService contains read and write operations for tenants
type TenantService struct {
	DB *sql.DB
}

// GetByDomain returns a tenant based on its domain
func (svc TenantService) GetByDomain(domain string) (*identity.Tenant, error) {
	tenant := &identity.Tenant{}

	row := svc.DB.QueryRow("SELECT id, name, subdomain FROM tenants WHERE subdomain = $1 OR cname = $2", extractSubdomain(domain), domain)
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Domain)
	if err != nil {
		return nil, identity.ErrNotFound
	}

	tenant.Domain = tenant.Domain + "." + env.GetCurrentDomain()
	return tenant, nil
}

func extractSubdomain(domain string) string {
	if idx := strings.Index(domain, "."); idx != -1 {
		return domain[:idx]
	}
	return domain
}
