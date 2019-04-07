package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"database/sql"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

type dbUser struct {
	ID            sql.NullInt64  `db:"id"`
	Name          sql.NullString `db:"name"`
	Email         sql.NullString `db:"email"`
	Tenant        *dbTenant      `db:"tenant"`
	Role          sql.NullInt64  `db:"role"`
	Status        sql.NullInt64  `db:"status"`
	AvatarType    sql.NullInt64  `db:"avatar_type"`
	AvatarBlobKey sql.NullString `db:"avatar_bkey"`
	Providers     []*dbUserProvider
}

type dbUserProvider struct {
	Name sql.NullString `db:"provider"`
	UID  sql.NullString `db:"provider_uid"`
}

func (u *dbUser) toModel(ctx context.Context) *models.User {
	if u == nil {
		return nil
	}

	avatarType := models.AvatarType(u.AvatarType.Int64)
	user := &models.User{
		ID:            int(u.ID.Int64),
		Name:          u.Name.String,
		Email:         u.Email.String,
		Tenant:        u.Tenant.toModel(),
		Role:          models.Role(u.Role.Int64),
		Providers:     make([]*models.UserProvider, len(u.Providers)),
		Status:        models.UserStatus(u.Status.Int64),
		AvatarType:    avatarType,
		AvatarBlobKey: u.AvatarBlobKey.String,
		AvatarURL:     buildAvatarURL(ctx, avatarType, int(u.ID.Int64), u.Name.String, u.AvatarBlobKey.String),
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
	ctx    context.Context
}

// NewUserStorage creates a new UserStorage
func NewUserStorage(trx *dbx.Trx, ctx context.Context) *UserStorage {
	return &UserStorage{trx: trx, ctx: ctx}
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
	user, err := s.getUser(s.trx, "id = $1", userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user with id '%d'", userID)
	}
	return user, nil
}

// GetByEmail returns a user based on given email
func (s *UserStorage) GetByEmail(email string) (*models.User, error) {
	email = strings.ToLower(email)
	user, err := s.getUser(s.trx, "email = $1 AND tenant_id = $2", email, s.tenant.ID)
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
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if err := s.trx.Get(&user.ID,
		"INSERT INTO users (name, email, created_at, tenant_id, role, status, avatar_type, avatar_bkey) VALUES ($1, $2, $3, $4, $5, $6, $7, '') RETURNING id",
		user.Name, user.Email, now, s.tenant.ID, user.Role, models.UserActive, models.AvatarTypeGravatar); err != nil {
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
	if settings.Avatar.Remove {
		settings.Avatar.BlobKey = ""
	}
	cmd := "UPDATE users SET name = $3, avatar_type = $4, avatar_bkey = $5 WHERE id = $1 AND tenant_id = $2"
	_, err := s.trx.Execute(cmd, s.user.ID, s.tenant.ID, settings.Name, settings.AvatarType, settings.Avatar.BlobKey)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}
	return nil
}

// GetByID returns a user based on given id
func (s *UserStorage) getUser(trx *dbx.Trx, filter string, args ...interface{}) (*models.User, error) {
	user := dbUser{}
	sql := fmt.Sprintf("SELECT id, name, email, tenant_id, role, status, avatar_type, avatar_bkey FROM users WHERE status != %d AND ", models.UserDeleted)
	err := trx.Get(&user, sql+filter, args...)
	if err != nil {
		return nil, err
	}

	err = trx.Select(&user.Providers, "SELECT provider_uid, provider FROM user_providers WHERE user_id = $1", user.ID.Int64)
	if err != nil {
		return nil, err
	}

	return user.toModel(s.ctx), nil
}

// GetAll return all users of current tenant
func (s *UserStorage) GetAll() ([]*models.User, error) {
	var users []*dbUser
	err := s.trx.Select(&users, `
		SELECT id, name, email, tenant_id, role, status, avatar_type, avatar_bkey
		FROM users 
		WHERE tenant_id = $1 
		AND status != $2
		ORDER BY id`, s.tenant.ID, models.UserDeleted)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all users")
	}

	var result = make([]*models.User, len(users))
	for i, user := range users {
		result[i] = user.toModel(s.ctx)
	}
	return result, nil
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
