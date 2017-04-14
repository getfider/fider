package postgres

import (
	"time"

	"database/sql"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
)

// UserStorage is used for user operations using a Postgres database
type UserStorage struct {
	DB *dbx.Database
}

// GetByID returns a user based on given id
func (s *UserStorage) GetByID(userID int) (*models.User, error) {
	return getUser(s.DB, "id = $1", userID)
}

// GetByEmail returns a user based on given email
func (s *UserStorage) GetByEmail(tenantID int, email string) (*models.User, error) {
	return getUser(s.DB, "email = $1 AND tenant_id = $2", email, tenantID)
}

// Register creates a new user based on given information
func (s *UserStorage) Register(user *models.User) error {
	tx, err := s.DB.Begin()
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
func (s *UserStorage) RegisterProvider(userID int, provider *models.UserProvider) error {
	now := time.Now()
	return s.DB.Execute("INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES ($1, $2, $3, $4)", userID, provider.Name, provider.UID, now)
}

// GetByID returns a user based on given id
func getUser(db *dbx.Database, filter string, args ...interface{}) (*models.User, error) {
	user := &models.User{}
	user.Tenant = &models.Tenant{}
	row := db.QueryRow("SELECT id, name, email, tenant_id, role FROM users WHERE "+filter, args...)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Tenant.ID, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, app.ErrNotFound
		}
		return nil, err
	}

	rows, err := db.Query("SELECT provider_uid, provider FROM user_providers WHERE user_id = $1", user.ID)
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
