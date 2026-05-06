package postgres

import (
	"context"
	"strings"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/disposable"
	"github.com/getfider/fider/app/pkg/errors"
)

type candidateUser struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func getDisposableUsers(ctx context.Context, q *query.GetDisposableUsers) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		rules := &query.GetEmailDomainRules{}
		if err := getEmailDomainRules(ctx, rules); err != nil {
			return errors.Wrap(err, "failed to load tenant rules")
		}
		deny := domainList(rules.Result.Deny)
		allow := domainList(rules.Result.Allow)

		candidates := []*candidateUser{}
		err := trx.Select(&candidates,
			`SELECT id, name, email FROM users
			 WHERE tenant_id = $1 AND status != $2 AND email != ''`,
			tenant.ID, enum.UserDeleted,
		)
		if err != nil {
			return errors.Wrap(err, "failed to query users")
		}

		matched := []*query.GetDisposableUsersRow{}
		for _, u := range candidates {
			if !disposable.IsBlocked(u.Email, deny, allow) {
				continue
			}
			matched = append(matched, &query.GetDisposableUsersRow{
				UserID: u.ID, Name: u.Name, Email: u.Email,
			})
		}
		q.Result.Total = len(matched)

		limit := q.Limit
		if limit <= 0 || limit > 200 {
			limit = 200
		}
		if len(matched) > limit {
			matched = matched[:limit]
		}

		for _, m := range matched {
			if err := trx.Get(&m.VoteCount,
				"SELECT COUNT(*) FROM post_votes WHERE user_id = $1 AND tenant_id = $2",
				m.UserID, tenant.ID); err != nil {
				return errors.Wrap(err, "failed to count votes")
			}
			if err := trx.Get(&m.PostCount,
				"SELECT COUNT(*) FROM posts WHERE user_id = $1 AND tenant_id = $2",
				m.UserID, tenant.ID); err != nil {
				return errors.Wrap(err, "failed to count posts")
			}
			if err := trx.Get(&m.CommentCount,
				"SELECT COUNT(*) FROM comments WHERE user_id = $1 AND tenant_id = $2",
				m.UserID, tenant.ID); err != nil {
				return errors.Wrap(err, "failed to count comments")
			}
		}

		q.Result.Users = matched
		return nil
	})
}

func domainList(rules []*entity.EmailDomainRule) []string {
	out := make([]string, 0, len(rules))
	for _, r := range rules {
		out = append(out, strings.ToLower(r.Domain))
	}
	return out
}

func bulkDeleteUsersByID(ctx context.Context, c *cmd.BulkDeleteUsersByID) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		c.Result = 0
		for _, id := range c.UserIDs {
			if err := deleteUserInternal(trx, tenant.ID, id); err != nil {
				return err
			}
			c.Result++
		}
		return nil
	})
}
