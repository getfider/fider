package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestEmailDomainRule_AddListDelete(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	addDeny := &cmd.AddEmailDomainRule{Domain: "example.com", RuleType: entity.EmailDomainRuleDeny}
	Expect(bus.Dispatch(demoTenantCtx, addDeny)).IsNil()
	Expect(addDeny.Result).IsNotNil()
	Expect(addDeny.Result.ID > 0).IsTrue()
	Expect(addDeny.Result.Domain).Equals("example.com")
	Expect(addDeny.Result.RuleType).Equals(entity.EmailDomainRuleDeny)

	addAllow := &cmd.AddEmailDomainRule{Domain: "trusted.com", RuleType: entity.EmailDomainRuleAllow}
	Expect(bus.Dispatch(demoTenantCtx, addAllow)).IsNil()

	list := &query.GetEmailDomainRules{}
	Expect(bus.Dispatch(demoTenantCtx, list)).IsNil()
	Expect(len(list.Result.Deny)).Equals(1)
	Expect(len(list.Result.Allow)).Equals(1)
	Expect(list.Result.Deny[0].Domain).Equals("example.com")
	Expect(list.Result.Allow[0].Domain).Equals("trusted.com")

	Expect(bus.Dispatch(demoTenantCtx, &cmd.DeleteEmailDomainRule{ID: addDeny.Result.ID})).IsNil()

	list2 := &query.GetEmailDomainRules{}
	Expect(bus.Dispatch(demoTenantCtx, list2)).IsNil()
	Expect(len(list2.Result.Deny)).Equals(0)
	Expect(len(list2.Result.Allow)).Equals(1)
}

func TestEmailDomainRule_DuplicateRejected(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	first := &cmd.AddEmailDomainRule{Domain: "dup.com", RuleType: entity.EmailDomainRuleDeny}
	Expect(bus.Dispatch(demoTenantCtx, first)).IsNil()

	dup := &cmd.AddEmailDomainRule{Domain: "dup.com", RuleType: entity.EmailDomainRuleDeny}
	err := bus.Dispatch(demoTenantCtx, dup)
	Expect(err).IsNotNil() // unique constraint violation
}

func TestEmailDomainRule_TenantIsolation(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	addForDemo := &cmd.AddEmailDomainRule{Domain: "isolated.com", RuleType: entity.EmailDomainRuleDeny}
	Expect(bus.Dispatch(demoTenantCtx, addForDemo)).IsNil()

	list := &query.GetEmailDomainRules{}
	Expect(bus.Dispatch(avengersTenantCtx, list)).IsNil()
	Expect(len(list.Result.Deny)).Equals(0)
	Expect(len(list.Result.Allow)).Equals(0)
}
