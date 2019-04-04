package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
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

func countUsers(ctx context.Context, q *query.CountUsers) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		var count int
		err := trx.Scalar(&count, "SELECT COUNT(*) FROM users WHERE tenant_id = $1", tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to count users")
		}
		q.Result = count
		return nil
	})
}

func blockUser(ctx context.Context, c *cmd.BlockUser) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		if _, err := trx.Execute(
			"UPDATE users SET status = $3 WHERE id = $1 AND tenant_id = $2",
			c.UserID, tenant.ID, models.UserBlocked,
		); err != nil {
			return errors.Wrap(err, "failed to block user")
		}
		return nil
	})
}

func unblockUser(ctx context.Context, c *cmd.UnblockUser) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		if _, err := trx.Execute(
			"UPDATE users SET status = $3 WHERE id = $1 AND tenant_id = $2",
			c.UserID, tenant.ID, models.UserActive,
		); err != nil {
			return errors.Wrap(err, "failed to unblock user")
		}
		return nil
	})
}

func regenerateAPIKey(ctx context.Context, c *cmd.RegenerateAPIKey) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		apiKey := models.GenerateSecretKey()

		if _, err := trx.Execute(
			"UPDATE users SET api_key = $3, api_key_date = $4 WHERE id = $1 AND tenant_id = $2",
			user.ID, tenant.ID, apiKey, time.Now(),
		); err != nil {
			return errors.Wrap(err, "failed to update current user's API Key")
		}

		c.Result = apiKey
		return nil
	})
}

func getUserByAPIKey(ctx context.Context, q *query.GetUserByAPIKey) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		result, err := queryUser(ctx, trx, "api_key = $1 AND tenant_id = $2", q.APIKey, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get user with API Key '%s'", q.APIKey)
		}
		q.Result = result
		return nil
	})
}

func queryUser(ctx context.Context, trx *dbx.Trx, filter string, args ...interface{}) (*models.User, error) {
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

	return user.toModel(ctx), nil
}
