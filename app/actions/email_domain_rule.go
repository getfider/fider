package actions

import (
	"context"
	"regexp"
	"strings"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/validate"
)

var domainRuleRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)+$`)

// AddEmailDomainRule is the input model for adding a deny/allow rule
type AddEmailDomainRule struct {
	Domain   string `json:"domain"`
	RuleType string `json:"ruleType"`
}

// IsAuthorized returns true if current user is administrator
func (action *AddEmailDomainRule) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *AddEmailDomainRule) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()
	action.Domain = strings.ToLower(strings.TrimSpace(action.Domain))

	if action.Domain == "" {
		result.AddFieldFailure("domain", "Domain is required.")
		return result
	}
	if len(action.Domain) > 255 {
		result.AddFieldFailure("domain", "Domain must have less than 255 characters.")
	}
	if !domainRuleRegex.MatchString(action.Domain) {
		result.AddFieldFailure("domain", "Domain is not valid.")
	}
	if action.RuleType != entity.EmailDomainRuleDeny && action.RuleType != entity.EmailDomainRuleAllow {
		result.AddFieldFailure("ruleType", "Invalid rule type.")
	}
	return result
}
