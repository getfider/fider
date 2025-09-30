package dbEntities

import (
	"time"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/dbx"
)

type emailVerification struct {
	ID         int                        `db:"id"`
	Name       string                     `db:"name"`
	Email      string                     `db:"email"`
	Key        string                     `db:"key"`
	Kind       enum.EmailVerificationKind `db:"kind"`
	UserID     dbx.NullInt                `db:"user_id"`
	CreatedAt  time.Time                  `db:"created_at"`
	ExpiresAt  time.Time                  `db:"expires_at"`
	VerifiedAt dbx.NullTime               `db:"verified_at"`
}

func (t *emailVerification) toModel() *entity.EmailVerification {
	model := &entity.EmailVerification{
		Name:       t.Name,
		Email:      t.Email,
		Key:        t.Key,
		Kind:       t.Kind,
		CreatedAt:  t.CreatedAt,
		ExpiresAt:  t.ExpiresAt,
		VerifiedAt: nil,
	}

	if t.VerifiedAt.Valid {
		model.VerifiedAt = &t.VerifiedAt.Time
	}

	if t.UserID.Valid {
		model.UserID = int(t.UserID.Int64)
	}

	return model
}
