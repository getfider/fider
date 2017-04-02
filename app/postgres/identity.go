package postgres

import (
	"strings"
	"time"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/dbx"
	"github.com/WeCanHearYou/wechy/app/toolbox/env"
)

// UserService is used for user operations using a Postgres database
type UserService struct {
	DB *dbx.Database
}

// GetByEmail returns a user based on given email
func (svc UserService) GetByEmail(email string) (*app.User, error) {
	user := &app.User{}
	row := svc.DB.QueryRow("SELECT id, name, email FROM users WHERE email = $1", email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, app.ErrNotFound
	}

	return user, nil
}

// Register creates a new user based on given information
func (svc UserService) Register(user *app.User) error {
	tx, err := svc.DB.Begin()
	if err != nil {
		return err
	}

	now := time.Now()

	if err = tx.QueryRow("INSERT INTO users (name, email, created_on, tenant_id) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Email, now, user.Tenant.ID).Scan(&user.ID); err != nil {
		tx.Rollback()
		return err
	}

	for _, provider := range user.Providers {
		if _, err = tx.Exec("INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES ($1, $2, $3, $4)", user.ID, provider.Name, provider.UID, now); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// TenantService contains read and write operations for tenants
type TenantService struct {
	DB *dbx.Database
}

// GetByDomain returns a tenant based on its domain
func (svc TenantService) GetByDomain(domain string) (*app.Tenant, error) {
	tenant := &app.Tenant{}

	row := svc.DB.QueryRow("SELECT id, name, subdomain FROM tenants WHERE subdomain = $1 OR cname = $2", extractSubdomain(domain), domain)
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Domain)
	if err != nil {
		return nil, app.ErrNotFound
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
