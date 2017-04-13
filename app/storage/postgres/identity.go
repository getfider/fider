package postgres

import (
	"strings"
	"time"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/pkg/env"
)

// UserStorage is used for user operations using a Postgres database
type UserStorage struct {
	DB *dbx.Database
}

// GetByID returns a user based on given id
func (svc *UserStorage) GetByID(userID int) (*models.User, error) {
	user := &models.User{}
	user.Tenant = &models.Tenant{}
	row := svc.DB.QueryRow("SELECT id, name, email, tenant_id, role FROM users WHERE id = $1", userID)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Tenant.ID, &user.Role); err != nil {
		return nil, app.ErrNotFound
	}

	rows, err := svc.DB.Query("SELECT provider_uid, provider FROM user_providers WHERE user_id = $1", user.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		p := &models.UserProvider{}
		rows.Scan(&p.UID, &p.Name)
		user.Providers = append(user.Providers, p)
	}

	return user, nil
}

// GetByEmail returns a user based on given email
func (svc *UserStorage) GetByEmail(tenantID int, email string) (*models.User, error) {
	user := &models.User{}
	user.Tenant = &models.Tenant{}
	row := svc.DB.QueryRow("SELECT id, name, email, tenant_id, role FROM users WHERE email = $1 AND tenant_id = $2", email, tenantID)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Tenant.ID, &user.Role); err != nil {
		return nil, app.ErrNotFound
	}

	rows, err := svc.DB.Query("SELECT provider_uid, provider FROM user_providers WHERE user_id = $1", user.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		p := &models.UserProvider{}
		rows.Scan(&p.UID, &p.Name)
		user.Providers = append(user.Providers, p)
	}

	return user, nil
}

// Register creates a new user based on given information
func (svc *UserStorage) Register(user *models.User) error {
	tx, err := svc.DB.Begin()
	if err != nil {
		return err
	}

	now := time.Now()

	if err = tx.QueryRow("INSERT INTO users (name, email, created_on, tenant_id, role) VALUES ($1, $2, $3, $4, $5) RETURNING id", user.Name, user.Email, now, user.Tenant.ID, user.Role).Scan(&user.ID); err != nil {
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

// RegisterProvider adds given provider to userID
func (svc *UserStorage) RegisterProvider(userID int, provider *models.UserProvider) error {
	now := time.Now()
	return svc.DB.Execute("INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES ($1, $2, $3, $4)", userID, provider.Name, provider.UID, now)
}

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	DB *dbx.Database
}

// GetByDomain returns a tenant based on its domain
func (svc *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := &models.Tenant{}

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
