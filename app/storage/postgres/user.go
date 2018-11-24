package postgres

import (
	"fmt"
	"strings"
	"time"

	"database/sql"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

type dbUser struct {
	ID        sql.NullInt64  `db:"id"`
	Name      sql.NullString `db:"name"`
	Email     sql.NullString `db:"email"`
	Tenant    *dbTenant      `db:"tenant"`
	Role      sql.NullInt64  `db:"role"`
	Status    sql.NullInt64  `db:"status"`
	Providers []*dbUserProvider
}

type dbUserProvider struct {
	Name sql.NullString `db:"provider"`
	UID  sql.NullString `db:"provider_uid"`
}

func (u *dbUser) toModel() *models.User {
	if u == nil {
		return nil
	}

	user := &models.User{
		ID:        int(u.ID.Int64),
		Name:      u.Name.String,
		Email:     u.Email.String,
		Tenant:    u.Tenant.toModel(),
		Role:      models.Role(u.Role.Int64),
		Providers: make([]*models.UserProvider, len(u.Providers)),
		Status:    models.UserStatus(u.Status.Int64),
	}

	for i, p := range u.Providers {
		user.Providers[i] = &models.UserProvider{
			Name: p.Name.String,
			UID:  p.UID.String,
		}
	}

	return user
}

type dbUserSetting struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

// UserStorage is used for user operations using a Postgres database
type UserStorage struct {
	tenant *models.Tenant
	user   *models.User
	trx    *dbx.Trx
}

// NewUserStorage creates a new UserStorage
func NewUserStorage(trx *dbx.Trx) *UserStorage {
	return &UserStorage{trx: trx}
}

// SetCurrentTenant to current context
func (s *UserStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.tenant = tenant
}

// SetCurrentUser to current context
func (s *UserStorage) SetCurrentUser(user *models.User) {
	s.user = user
}

// GetByID returns a user based on given id
func (s *UserStorage) GetByID(userID int) (*models.User, error) {
	user, err := getUser(s.trx, "id = $1", userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user with id '%d'", userID)
	}
	return user, nil
}

// GetByEmail returns a user based on given email
func (s *UserStorage) GetByEmail(email string) (*models.User, error) {
	user, err := getUser(s.trx, "email = $1 AND tenant_id = $2", email, s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user with email '%s'", email)
	}
	return user, nil
}

// GetByProvider returns a user based on provider details
func (s *UserStorage) GetByProvider(provider string, uid string) (*models.User, error) {
	var userID int
	query := `
	SELECT user_id 
	FROM user_providers up 
	INNER JOIN users u 
	ON u.id = up.user_id 
	AND u.tenant_id = up.tenant_id 
	WHERE up.provider = $1 
	AND up.provider_uid = $2 
	AND u.tenant_id = $3`
	if err := s.trx.Scalar(&userID, query, provider, uid, s.tenant.ID); err != nil {
		return nil, errors.Wrap(err, "failed to get user by provider '%s' and uid '%s'", provider, uid)
	}
	return s.GetByID(userID)
}

// Register creates a new user based on given information
func (s *UserStorage) Register(user *models.User) error {
	now := time.Now()
	user.Status = models.UserActive
	user.Email = strings.TrimSpace(user.Email)
	if err := s.trx.Get(&user.ID,
		"INSERT INTO users (name, email, created_at, tenant_id, role, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Name, user.Email, now, s.tenant.ID, user.Role, models.UserActive); err != nil {
		return errors.Wrap(err, "failed to register new user")
	}

	for _, provider := range user.Providers {
		if _, err := s.trx.Execute("INSERT INTO user_providers (tenant_id, user_id, provider, provider_uid, created_at) VALUES ($1, $2, $3, $4, $5)", s.tenant.ID, user.ID, provider.Name, provider.UID, now); err != nil {
			return errors.Wrap(err, "failed to add provider to new user")
		}
	}

	return nil
}

// RegisterProvider adds given provider to userID
func (s *UserStorage) RegisterProvider(userID int, provider *models.UserProvider) error {
	cmd := "INSERT INTO user_providers (tenant_id, user_id, provider, provider_uid, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := s.trx.Execute(cmd, s.tenant.ID, userID, provider.Name, provider.UID, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to add provider '%s' to user with id '%d'", provider.Name, userID)
	}
	return nil
}

// Update user profile
func (s *UserStorage) Update(settings *models.UpdateUserSettings) error {
	cmd := "UPDATE users SET name = $2 WHERE id = $1 AND tenant_id = $3"
	_, err := s.trx.Execute(cmd, s.user.ID, settings.Name, s.tenant.ID)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}
	return nil
}

// UpdateSettings of given user
func (s *UserStorage) UpdateSettings(settings map[string]string) error {
	if s.user != nil && settings != nil && len(settings) > 0 {
		query := `
		INSERT INTO user_settings (tenant_id, user_id, key, value)
		VALUES ($1, $2, $3, $4) ON CONFLICT (user_id, key) DO UPDATE SET value = $4
		`

		for key, value := range settings {
			_, err := s.trx.Execute(query, s.tenant.ID, s.user.ID, key, value)
			if err != nil {
				return errors.Wrap(err, "failed to update user settings")
			}
		}
	}

	return nil
}

// GetUserSettings returns current user's settings
func (s *UserStorage) GetUserSettings() (map[string]string, error) {
	if s.user == nil {
		return make(map[string]string, 0), nil
	}

	var settings []*dbUserSetting
	err := s.trx.Select(&settings, "SELECT key, value FROM user_settings WHERE user_id = $1 AND tenant_id = $2", s.user.ID, s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user settings")
	}

	var result = make(map[string]string, len(settings))

	for _, e := range models.AllNotificationEvents {
		for _, r := range e.DefaultEnabledUserRoles {
			if r == s.user.Role {
				result[e.UserSettingsKeyName] = e.DefaultSettingValue
			}
		}
	}

	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

// ChangeRole of given user
func (s *UserStorage) ChangeRole(userID int, role models.Role) error {
	cmd := "UPDATE users SET role = $3 WHERE id = $1 AND tenant_id = $2"
	_, err := s.trx.Execute(cmd, userID, s.tenant.ID, role)
	if err != nil {
		return errors.Wrap(err, "failed to change user's role")
	}
	return nil
}

// ChangeEmail of given user
func (s *UserStorage) ChangeEmail(userID int, email string) error {
	cmd := "UPDATE users SET email = $3 WHERE id = $1 AND tenant_id = $2"
	_, err := s.trx.Execute(cmd, userID, s.tenant.ID, email)
	if err != nil {
		return errors.Wrap(err, "failed to update user's email")
	}
	return nil
}

// GetByID returns a user based on given id
func getUser(trx *dbx.Trx, filter string, args ...interface{}) (*models.User, error) {
	user := dbUser{}
	sql := fmt.Sprintf("SELECT id, name, email, tenant_id, role, status FROM users WHERE status != %d AND ", models.UserDeleted)
	err := trx.Get(&user, sql+filter, args...)
	if err != nil {
		return nil, err
	}

	err = trx.Select(&user.Providers, "SELECT provider_uid, provider FROM user_providers WHERE user_id = $1", user.ID.Int64)
	if err != nil {
		return nil, err
	}

	return user.toModel(), nil
}

// GetAll return all users of current tenant
func (s *UserStorage) GetAll() ([]*models.User, error) {
	var users []*dbUser
	err := s.trx.Select(&users, `
		SELECT id, name, email, tenant_id, role, status 
		FROM users 
		WHERE tenant_id = $1 
		AND status != $2
		ORDER BY id`, s.tenant.ID, models.UserDeleted)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all users")
	}

	var result = make([]*models.User, len(users))
	for i, user := range users {
		result[i] = user.toModel()
	}
	return result, nil
}

// HasSubscribedTo returns true if current user is receiving notification from specific post
func (s *UserStorage) HasSubscribedTo(postID int) (bool, error) {
	if s.user == nil {
		return false, nil
	}

	var status int
	err := s.trx.Scalar(&status, "SELECT status FROM post_subscribers WHERE user_id = $1 AND post_id = $2", s.user.ID, postID)
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return false, errors.Wrap(err, "failed to get subscription status")
	}

	if errors.Cause(err) == app.ErrNotFound {
		for _, e := range models.AllNotificationEvents {
			for _, r := range e.RequiresSubscriptionUserRoles {
				if r == s.user.Role {
					return false, nil
				}
			}
		}
		return true, nil
	}

	if status == 1 {
		return true, nil
	}

	return false, nil
}

// Delete removes current user personal data and mark it as deleted
func (s *UserStorage) Delete() error {
	if _, err := s.trx.Execute(
		"UPDATE users SET role = $3, status = $4, name = '', email = '', api_key = null, api_key_date = null WHERE id = $1 AND tenant_id = $2",
		s.user.ID, s.tenant.ID, models.RoleVisitor, models.UserDeleted,
	); err != nil {
		return errors.Wrap(err, "failed to delete current user")
	}

	var tables = []struct {
		name       string
		userColumn string
	}{
		{"user_providers", "user_id"},
		{"user_settings", "user_id"},
		{"notifications", "user_id"},
		{"notifications", "author_id"},
		{"post_votes", "user_id"},
		{"post_subscribers", "user_id"},
		{"email_verifications", "user_id"},
	}

	for _, table := range tables {
		if _, err := s.trx.Execute(
			fmt.Sprintf("DELETE FROM %s WHERE %s = $1 AND tenant_id = $2", table.name, table.userColumn),
			s.user.ID, s.tenant.ID,
		); err != nil {
			return errors.Wrap(err, "failed to delete current user's %s records", table)
		}
	}

	return nil
}

// RegenerateAPIKey generates a new API Key and returns it
func (s *UserStorage) RegenerateAPIKey() (string, error) {
	apiKey := models.GenerateSecretKey()

	if _, err := s.trx.Execute(
		"UPDATE users SET api_key = $3, api_key_date = $4 WHERE id = $1 AND tenant_id = $2",
		s.user.ID, s.tenant.ID, apiKey, time.Now(),
	); err != nil {
		return "", errors.Wrap(err, "failed to update current user's API Key")
	}

	return apiKey, nil
}

// GetByAPIKey returns a user based on its API key
func (s *UserStorage) GetByAPIKey(apiKey string) (*models.User, error) {
	user, err := getUser(s.trx, "api_key = $1 AND tenant_id = $2", apiKey, s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user with API Key '%s'", apiKey)
	}
	return user, nil
}

// Block a given user from using Fider
func (s *UserStorage) Block(userID int) error {
	if _, err := s.trx.Execute(
		"UPDATE users SET status = $3 WHERE id = $1 AND tenant_id = $2",
		userID, s.tenant.ID, models.UserBlocked,
	); err != nil {
		return errors.Wrap(err, "failed to block user")
	}
	return nil
}

// Unblock a given user so that they can use Fider again
func (s *UserStorage) Unblock(userID int) error {
	if _, err := s.trx.Execute(
		"UPDATE users SET status = $3 WHERE id = $1 AND tenant_id = $2",
		userID, s.tenant.ID, models.UserActive,
	); err != nil {
		return errors.Wrap(err, "failed to unblock user")
	}
	return nil
}
