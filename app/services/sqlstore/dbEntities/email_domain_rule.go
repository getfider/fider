package dbEntities

import (
	"time"

	"github.com/getfider/fider/app/models/entity"
)

type EmailDomainRule struct {
	ID        int       `db:"id"`
	Domain    string    `db:"domain"`
	RuleType  string    `db:"rule_type"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy *int      `db:"created_by"`
}

func (r *EmailDomainRule) ToModel() *entity.EmailDomainRule {
	if r == nil {
		return nil
	}
	return &entity.EmailDomainRule{
		ID:        r.ID,
		Domain:    r.Domain,
		RuleType:  r.RuleType,
		CreatedAt: r.CreatedAt,
		CreatedBy: r.CreatedBy,
	}
}
