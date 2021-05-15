package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
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

func (u *dbUser) toModel(ctx context.Context) *entity.User {
	if u == nil {
		return nil
	}

	avatarType := enum.AvatarType(u.AvatarType.Int64)
	user := &entity.User{
		ID:            int(u.ID.Int64),
		Name:          u.Name.String,
		Email:         u.Email.String,
		Tenant:        u.Tenant.toModel(),
		Role:          enum.Role(u.Role.Int64),
		Providers:     make([]*entity.UserProvider, len(u.Providers)),
		Status:        enum.UserStatus(u.Status.Int64),
		AvatarType:    avatarType,
		AvatarBlobKey: u.AvatarBlobKey.String,
		AvatarURL:     buildAvatarURL(ctx, avatarType, int(u.ID.Int64), u.Name.String, u.AvatarBlobKey.String),
	}

	for i, p := range u.Providers {
		user.Providers[i] = &entity.UserProvider{
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
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
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
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if _, err := trx.Execute(
			"UPDATE users SET status = $3 WHERE id = $1 AND tenant_id = $2",
			c.UserID, tenant.ID, enum.UserBlocked,
		); err != nil {
			return errors.Wrap(err, "failed to block user")
		}
		return nil
	})
}

func unblockUser(ctx context.Context, c *cmd.UnblockUser) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if _, err := trx.Execute(
			"UPDATE users SET status = $3 WHERE id = $1 AND tenant_id = $2",
			c.UserID, tenant.ID, enum.UserActive,
		); err != nil {
			return errors.Wrap(err, "failed to unblock user")
		}
		return nil
	})
}

func deleteCurrentUser(ctx context.Context, c *cmd.DeleteCurrentUser) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if _, err := trx.Execute(
			"UPDATE users SET role = $3, status = $4, name = '', email = '', api_key = null, api_key_date = null WHERE id = $1 AND tenant_id = $2",
			user.ID, tenant.ID, enum.RoleVisitor, enum.UserDeleted,
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
			if _, err := trx.Execute(
				fmt.Sprintf("DELETE FROM %s WHERE %s = $1 AND tenant_id = $2", table.name, table.userColumn),
				user.ID, tenant.ID,
			); err != nil {
				return errors.Wrap(err, "failed to delete current user's %s records", table)
			}
		}

		return nil
	})
}

func regenerateAPIKey(ctx context.Context, c *cmd.RegenerateAPIKey) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		apiKey := entity.GenerateEmailVerificationKey()

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
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		result, err := queryUser(ctx, trx, "api_key = $1 AND tenant_id = $2", q.APIKey, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get user with API Key '%s'", q.APIKey)
		}
		q.Result = result
		return nil
	})
}

func userSubscribedTo(ctx context.Context, q *query.UserSubscribedTo) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if user == nil {
			q.Result = false
			return nil
		}

		var status int
		err := trx.Scalar(&status, "SELECT status FROM post_subscribers WHERE user_id = $1 AND post_id = $2", user.ID, q.PostID)
		if err != nil && errors.Cause(err) != app.ErrNotFound {
			return errors.Wrap(err, "failed to get subscription status")
		}

		if errors.Cause(err) == app.ErrNotFound {
			for _, e := range enum.AllNotificationEvents {
				for _, r := range e.RequiresSubscriptionUserRoles {
					if r == user.Role {
						q.Result = false
						return nil
					}
				}
			}
			q.Result = true
			return nil
		}

		if status == 1 {
			q.Result = true
			return nil
		}

		q.Result = false
		return nil
	})
}

func changeUserRole(ctx context.Context, c *cmd.ChangeUserRole) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		cmd := "UPDATE users SET role = $3 WHERE id = $1 AND tenant_id = $2"
		_, err := trx.Execute(cmd, c.UserID, tenant.ID, c.Role)
		if err != nil {
			return errors.Wrap(err, "failed to change user's role")
		}
		return nil
	})
}

func changeUserEmail(ctx context.Context, c *cmd.ChangeUserEmail) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		cmd := "UPDATE users SET email = $3 WHERE id = $1 AND tenant_id = $2"
		_, err := trx.Execute(cmd, c.UserID, tenant.ID, strings.ToLower(c.Email))
		if err != nil {
			return errors.Wrap(err, "failed to update user's email")
		}
		return nil
	})
}

func updateCurrentUserSettings(ctx context.Context, c *cmd.UpdateCurrentUserSettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if user != nil && c.Settings != nil && len(c.Settings) > 0 {
			query := `
			INSERT INTO user_settings (tenant_id, user_id, key, value)
			VALUES ($1, $2, $3, $4) ON CONFLICT (user_id, key) DO UPDATE SET value = $4
			`

			for key, value := range c.Settings {
				_, err := trx.Execute(query, tenant.ID, user.ID, key, value)
				if err != nil {
					return errors.Wrap(err, "failed to update user settings")
				}
			}
		}

		return nil
	})
}

func getCurrentUserSettings(ctx context.Context, q *query.GetCurrentUserSettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make(map[string]string)

		var settings []*dbUserSetting
		err := trx.Select(&settings, "SELECT key, value FROM user_settings WHERE user_id = $1 AND tenant_id = $2", user.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get user settings")
		}

		for _, e := range enum.AllNotificationEvents {
			for _, r := range e.DefaultEnabledUserRoles {
				if r == user.Role {
					q.Result[e.UserSettingsKeyName] = e.DefaultSettingValue
				}
			}
		}

		for _, s := range settings {
			q.Result[s.Key] = s.Value
		}

		return nil
	})
}

func registerUser(ctx context.Context, c *cmd.RegisterUser) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, _ *entity.User) error {
		now := time.Now()
		c.User.Status = enum.UserActive
		c.User.Email = strings.ToLower(strings.TrimSpace(c.User.Email))
		if err := trx.Get(&c.User.ID,
			"INSERT INTO users (name, email, created_at, tenant_id, role, status, avatar_type, avatar_bkey) VALUES ($1, $2, $3, $4, $5, $6, $7, '') RETURNING id",
			c.User.Name, c.User.Email, now, tenant.ID, c.User.Role, enum.UserActive, enum.AvatarTypeGravatar); err != nil {
			return errors.Wrap(err, "failed to register new user")
		}

		for _, provider := range c.User.Providers {
			cmd := "INSERT INTO user_providers (tenant_id, user_id, provider, provider_uid, created_at) VALUES ($1, $2, $3, $4, $5)"
			if _, err := trx.Execute(cmd, tenant.ID, c.User.ID, provider.Name, provider.UID, now); err != nil {
				return errors.Wrap(err, "failed to add provider to new user")
			}
		}

		return nil
	})
}

func registerUserProvider(ctx context.Context, c *cmd.RegisterUserProvider) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		cmd := "INSERT INTO user_providers (tenant_id, user_id, provider, provider_uid, created_at) VALUES ($1, $2, $3, $4, $5)"
		_, err := trx.Execute(cmd, tenant.ID, c.UserID, c.ProviderName, c.ProviderUID, time.Now())
		if err != nil {
			return errors.Wrap(err, "failed to add provider '%s:%s' to user with id '%d'", c.ProviderName, c.ProviderUID, c.UserID)
		}
		return nil
	})
}

func updateCurrentUser(ctx context.Context, c *cmd.UpdateCurrentUser) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if c.Avatar.Remove {
			c.Avatar.BlobKey = ""
		}
		cmd := "UPDATE users SET name = $3, avatar_type = $4, avatar_bkey = $5 WHERE id = $1 AND tenant_id = $2"
		_, err := trx.Execute(cmd, user.ID, tenant.ID, c.Name, c.AvatarType, c.Avatar.BlobKey)
		if err != nil {
			return errors.Wrap(err, "failed to update user")
		}
		return nil
	})
}

func getUserByID(ctx context.Context, q *query.GetUserByID) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		u, err := queryUser(ctx, trx, "id = $1", q.UserID)
		if err != nil {
			return errors.Wrap(err, "failed to get user with id '%d'", q.UserID)
		}
		q.Result = u
		return nil
	})
}

func getUserByEmail(ctx context.Context, q *query.GetUserByEmail) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		email := strings.ToLower(q.Email)
		u, err := queryUser(ctx, trx, "email = $1 AND tenant_id = $2", email, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get user with email '%s'", email)
		}
		q.Result = u
		return nil
	})
}

func getUserByProvider(ctx context.Context, q *query.GetUserByProvider) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var userID int
		if err := trx.Scalar(&userID, `
			SELECT user_id 
			FROM user_providers up 
			INNER JOIN users u 
			ON u.id = up.user_id 
			AND u.tenant_id = up.tenant_id 
			WHERE up.provider = $1 
			AND up.provider_uid = $2 
			AND u.tenant_id = $3`, q.Provider, q.UID, tenant.ID); err != nil {
			return errors.Wrap(err, "failed to get user by provider '%s' and uid '%s'", q.Provider, q.UID)
		}

		byID := &query.GetUserByID{UserID: userID}
		err := getUserByID(ctx, byID)
		q.Result = byID.Result
		return err
	})
}

func getAllUsers(ctx context.Context, q *query.GetAllUsers) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var users []*dbUser
		err := trx.Select(&users, `
			SELECT id, name, email, tenant_id, role, status, avatar_type, avatar_bkey
			FROM users 
			WHERE tenant_id = $1 
			AND status != $2
			ORDER BY id`, tenant.ID, enum.UserDeleted)
		if err != nil {
			return errors.Wrap(err, "failed to get all users")
		}

		q.Result = make([]*entity.User, len(users))
		for i, user := range users {
			q.Result[i] = user.toModel(ctx)
		}
		return nil
	})
}

func queryUser(ctx context.Context, trx *dbx.Trx, filter string, args ...interface{}) (*entity.User, error) {
	user := dbUser{}
	sql := fmt.Sprintf("SELECT id, name, email, tenant_id, role, status, avatar_type, avatar_bkey FROM users WHERE status != %d AND ", enum.UserDeleted)
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
