package dbEntities

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/web"
)

// User is the database mapping for users table
type User struct {
	ID            sql.NullInt64  `db:"id"`
	Name          sql.NullString `db:"name"`
	Email         sql.NullString `db:"email"`
	Tenant        *Tenant        `db:"tenant"`
	Role          sql.NullInt64  `db:"role"`
	Status        sql.NullInt64  `db:"status"`
	AvatarType    sql.NullInt64  `db:"avatar_type"`
	AvatarBlobKey sql.NullString `db:"avatar_bkey"`
	IsVerified    sql.NullBool   `db:"is_verified"`
	Providers     []*UserProvider
}

type UserProvider struct {
	Name sql.NullString `db:"provider"`
	UID  sql.NullString `db:"provider_uid"`
}

type userSetting struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

func (u *User) ToModel(ctx context.Context) *entity.User {
	if u == nil {
		return nil
	}

	tenant, ok := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	if !ok || tenant == nil {
		tenant = u.Tenant.ToModel()
	}

	avatarType := enum.AvatarType(u.AvatarType.Int64)
	avatarURL := ""
	if u.AvatarType.Valid {
		avatarURL = buildAvatarURL(ctx, avatarType, int(u.ID.Int64), u.Name.String, u.AvatarBlobKey.String)
	}

	user := &entity.User{
		ID:            int(u.ID.Int64),
		Name:          u.Name.String,
		Email:         u.Email.String,
		Tenant:        tenant,
		Role:          enum.Role(u.Role.Int64),
		Status:        enum.UserStatus(u.Status.Int64),
		AvatarType:    avatarType,
		AvatarBlobKey: u.AvatarBlobKey.String,
		AvatarURL:     avatarURL,
		IsVerified:    u.IsVerified.Bool,
	}

	if u.Providers != nil {
		user.Providers = make([]*entity.UserProvider, len(u.Providers))
		for i, p := range u.Providers {
			user.Providers[i] = &entity.UserProvider{
				Name: p.Name.String,
				UID:  p.UID.String,
			}
		}
	}

	return user
}

func buildAvatarURL(ctx context.Context, avatarType enum.AvatarType, id int, name, avatarBlobKey string) string {
	if name == "" {
		name = "-"
	}

	if avatarType == enum.AvatarTypeCustom {
		return web.AssetsURL(ctx, "/static/images/%s", avatarBlobKey)
	}
	return web.AssetsURL(ctx, "/static/avatars/%s/%d/%s", avatarType.String(), id, url.PathEscape(name))
}
