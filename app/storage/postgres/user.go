package postgres

import (
	"crypto/md5"
	"time"

	"database/sql"

	"fmt"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
)

type dbUser struct {
	ID        sql.NullInt64  `db:"id"`
	Name      sql.NullString `db:"name"`
	Email     sql.NullString `db:"email"`
	Tenant    *dbTenant      `db:"tenant"`
	Role      sql.NullInt64  `db:"role"`
	Providers []*dbUserProvider
}

type dbUserProvider struct {
	Name sql.NullString `db:"provider"`
	UID  sql.NullString `db:"provider_uid"`
}

func (u *dbUser) toModel() *models.User {
	user := &models.User{
		ID:        int(u.ID.Int64),
		Name:      u.Name.String,
		Email:     u.Email.String,
		Gravatar:  fmt.Sprintf("%x", md5.Sum([]byte(u.Email.String))),
		Tenant:    u.Tenant.toModel(),
		Role:      int(u.Role.Int64),
		Providers: make([]*models.UserProvider, len(u.Providers)),
	}

	for i, p := range u.Providers {
		user.Providers[i] = &models.UserProvider{
			Name: p.Name.String,
			UID:  p.UID.String,
		}
	}

	return user
}

// UserStorage is used for user operations using a Postgres database
type UserStorage struct {
	trx *dbx.Trx
}

// NewUserStorage creates a new UserStorage
func NewUserStorage(trx *dbx.Trx) *UserStorage {
	return &UserStorage{trx: trx}
}

// GetByID returns a user based on given id
func (s *UserStorage) GetByID(userID int) (*models.User, error) {
	return getUser(s.trx, "id = $1", userID)
}

// GetByEmail returns a user based on given email
func (s *UserStorage) GetByEmail(tenantID int, email string) (*models.User, error) {
	return getUser(s.trx, "email = $1 AND tenant_id = $2", email, tenantID)
}

// GetByProvider returns a user based on provider details
func (s *UserStorage) GetByProvider(tenantID int, provider string, uid string) (*models.User, error) {
	var userID int
	query := "SELECT user_id FROM user_providers up INNER JOIN users u ON u.id = up.user_id WHERE up.provider = $1 AND up.provider_uid = $2 AND u.tenant_id = $3"
	if err := s.trx.Scalar(&userID, query, provider, uid, tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, app.ErrNotFound
		} else if err != nil {
			return nil, err
		}
	}
	return s.GetByID(userID)
}

// Register creates a new user based on given information
func (s *UserStorage) Register(user *models.User) error {
	now := time.Now()
	if err := s.trx.QueryRow("INSERT INTO users (name, email, created_on, tenant_id, role) VALUES ($1, $2, $3, $4, $5) RETURNING id", user.Name, user.Email, now, user.Tenant.ID, user.Role).Scan(&user.ID); err != nil {
		return err
	}

	for _, provider := range user.Providers {
		if err := s.trx.Execute("INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES ($1, $2, $3, $4)", user.ID, provider.Name, provider.UID, now); err != nil {
			return err
		}
	}

	return nil
}

// RegisterProvider adds given provider to userID
func (s *UserStorage) RegisterProvider(userID int, provider *models.UserProvider) error {
	cmd := "INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES ($1, $2, $3, $4)"
	return s.trx.Execute(cmd, userID, provider.Name, provider.UID, time.Now())
}

// Update user settings
func (s *UserStorage) Update(userID int, settings *models.UpdateUserSettings) error {
	cmd := "UPDATE users SET name = $2 WHERE id = $1"
	return s.trx.Execute(cmd, userID, settings.Name)
}

// GetByID returns a user based on given id
func getUser(trx *dbx.Trx, filter string, args ...interface{}) (*models.User, error) {
	user := dbUser{}
	err := trx.Get(&user, "SELECT id, name, email, tenant_id, role FROM users WHERE "+filter, args...)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = trx.Select(&user.Providers, "SELECT provider_uid, provider FROM user_providers WHERE user_id = $1", user.ID.Int64)
	if err != nil {
		return nil, err
	}

	return user.toModel(), nil
}
