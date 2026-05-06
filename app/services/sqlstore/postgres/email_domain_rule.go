package postgres

import (
	"context"
	"strings"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

func addEmailDomainRule(ctx context.Context, c *cmd.AddEmailDomainRule) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		domain := strings.ToLower(strings.TrimSpace(c.Domain))
		var createdBy *int
		if user != nil {
			id := user.ID
			createdBy = &id
		}

		row := dbEntities.EmailDomainRule{}
		err := trx.Get(&row,
			`INSERT INTO email_domain_rules (tenant_id, domain, rule_type, created_by)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id, domain, rule_type, created_at, created_by`,
			tenant.ID, domain, c.RuleType, createdBy,
		)
		if err != nil {
			return errors.Wrap(err, "failed to add email domain rule")
		}
		c.Result = row.ToModel()
		return nil
	})
}

func deleteEmailDomainRule(ctx context.Context, c *cmd.DeleteEmailDomainRule) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(
			"DELETE FROM email_domain_rules WHERE id = $1 AND tenant_id = $2",
			c.ID, tenant.ID,
		)
		if err != nil {
			return errors.Wrap(err, "failed to delete email domain rule")
		}
		return nil
	})
}

func getEmailDomainRules(ctx context.Context, q *query.GetEmailDomainRules) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		rows := []*dbEntities.EmailDomainRule{}
		err := trx.Select(&rows,
			`SELECT id, domain, rule_type, created_at, created_by
			 FROM email_domain_rules
			 WHERE tenant_id = $1
			 ORDER BY created_at DESC`,
			tenant.ID,
		)
		if err != nil {
			return errors.Wrap(err, "failed to list email domain rules")
		}
		for _, r := range rows {
			m := r.ToModel()
			if m.RuleType == entity.EmailDomainRuleDeny {
				q.Result.Deny = append(q.Result.Deny, m)
			} else {
				q.Result.Allow = append(q.Result.Allow, m)
			}
		}
		return nil
	})
}
