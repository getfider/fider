package dbEntities

import (
	"time"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/dbx"
)

type EmailDomainRule struct {
	ID        int         `db:"id"`
	Domain    string      `db:"domain"`
	RuleType  string      `db:"rule_type"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy dbx.NullInt `db:"created_by"`
}

func (r *EmailDomainRule) ToModel() *entity.EmailDomainRule {
	if r == nil {
		return nil
	}
	var createdBy *int
	if r.CreatedBy.Valid {
		v := int(r.CreatedBy.Int64)
		createdBy = &v
	}
	return &entity.EmailDomainRule{
		ID:        r.ID,
		Domain:    r.Domain,
		RuleType:  r.RuleType,
		CreatedAt: r.CreatedAt,
		CreatedBy: createdBy,
	}
}
