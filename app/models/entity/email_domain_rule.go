package entity

import "time"

const (
	EmailDomainRuleDeny  = "deny"
	EmailDomainRuleAllow = "allow"
)

// EmailDomainRule is a per-tenant deny or allow rule for email domains.
type EmailDomainRule struct {
	ID        int       `json:"id"`
	Domain    string    `json:"domain"`
	RuleType  string    `json:"ruleType"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy *int      `json:"createdBy,omitempty"`
}
