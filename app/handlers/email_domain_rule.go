package handlers

import (
	"strconv"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// ListEmailDomainRules returns the deny/allow rules for the current tenant.
func ListEmailDomainRules() web.HandlerFunc {
	return func(c *web.Context) error {
		q := &query.GetEmailDomainRules{}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}
		return c.Ok(web.Map{
			"deny":  q.Result.Deny,
			"allow": q.Result.Allow,
		})
	}
}

// AddEmailDomainRule creates a new tenant-scoped deny or allow rule.
func AddEmailDomainRule() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.AddEmailDomainRule)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}
		add := &cmd.AddEmailDomainRule{Domain: action.Domain, RuleType: action.RuleType}
		if err := bus.Dispatch(c, add); err != nil {
			return c.Failure(err)
		}
		return c.Ok(add.Result)
	}
}

// DeleteEmailDomainRule removes a tenant-scoped rule by ID.
func DeleteEmailDomainRule() web.HandlerFunc {
	return func(c *web.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.NotFound()
		}
		if err := bus.Dispatch(c, &cmd.DeleteEmailDomainRule{ID: id}); err != nil {
			return c.Failure(err)
		}
		return c.Ok(web.Map{})
	}
}
