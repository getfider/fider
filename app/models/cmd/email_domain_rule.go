package cmd

import "github.com/getfider/fider/app/models/entity"

type AddEmailDomainRule struct {
	Domain   string
	RuleType string
	Result   *entity.EmailDomainRule
}

type DeleteEmailDomainRule struct {
	ID int
}
