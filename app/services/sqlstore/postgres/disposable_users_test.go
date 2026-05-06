package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestGetDisposableUsers_BundledMatch(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// mailinator.com is in the bundled disposable list.
	_, err := trx.Execute(
		`INSERT INTO users (name, email, tenant_id, role, status, avatar_type, avatar_bkey) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		"Fake User", "fake@mailinator.com", demoTenant.ID, 1, 1, 2, "",
	)
	Expect(err).IsNil()

	q := &query.GetDisposableUsers{Limit: 200}
	Expect(bus.Dispatch(demoTenantCtx, q)).IsNil()
	Expect(q.Result.Total >= 1).IsTrue()

	var found bool
	for _, u := range q.Result.Users {
		if u.Email == "fake@mailinator.com" {
			found = true
			break
		}
	}
	Expect(found).IsTrue()
}

func TestGetDisposableUsers_TenantDeny(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	_, err := trx.Execute(
		`INSERT INTO users (name, email, tenant_id, role, status, avatar_type, avatar_bkey) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		"Spam User", "user@spam.example", demoTenant.ID, 1, 1, 2, "",
	)
	Expect(err).IsNil()

	Expect(bus.Dispatch(demoTenantCtx, &cmd.AddEmailDomainRule{
		Domain: "spam.example", RuleType: entity.EmailDomainRuleDeny,
	})).IsNil()

	q := &query.GetDisposableUsers{Limit: 200}
	Expect(bus.Dispatch(demoTenantCtx, q)).IsNil()

	var found bool
	for _, u := range q.Result.Users {
		if u.Email == "user@spam.example" {
			found = true
			break
		}
	}
	Expect(found).IsTrue()
}

func TestGetDisposableUsers_AllowOverridesBundled(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	_, err := trx.Execute(
		`INSERT INTO users (name, email, tenant_id, role, status, avatar_type, avatar_bkey) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		"Allowed", "ok@mailinator.com", demoTenant.ID, 1, 1, 2, "",
	)
	Expect(err).IsNil()

	Expect(bus.Dispatch(demoTenantCtx, &cmd.AddEmailDomainRule{
		Domain: "mailinator.com", RuleType: entity.EmailDomainRuleAllow,
	})).IsNil()

	q := &query.GetDisposableUsers{Limit: 200}
	Expect(bus.Dispatch(demoTenantCtx, q)).IsNil()
	for _, u := range q.Result.Users {
		Expect(u.Email).NotEquals("ok@mailinator.com")
	}
}

func TestBulkDeleteUsersByID(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	var id1, id2 int
	err := trx.Get(&id1,
		`INSERT INTO users (name, email, tenant_id, role, status, avatar_type, avatar_bkey) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		"Bulk Target One", "bulk1@mailinator.com", demoTenant.ID, 1, 1, 2, "",
	)
	Expect(err).IsNil()

	err = trx.Get(&id2,
		`INSERT INTO users (name, email, tenant_id, role, status, avatar_type, avatar_bkey) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		"Bulk Target Two", "bulk2@mailinator.com", demoTenant.ID, 1, 1, 2, "",
	)
	Expect(err).IsNil()

	c := &cmd.BulkDeleteUsersByID{UserIDs: []int{id1, id2}}
	Expect(bus.Dispatch(demoTenantCtx, c)).IsNil()
	Expect(c.Result).Equals(2)

	type minimalUser struct {
		Status int    `db:"status"`
		Email  string `db:"email"`
	}
	var u1, u2 minimalUser
	Expect(trx.Get(&u1, "SELECT status, email FROM users WHERE id = $1", id1)).IsNil()
	Expect(trx.Get(&u2, "SELECT status, email FROM users WHERE id = $1", id2)).IsNil()
	Expect(u1.Email).Equals("")
	Expect(u2.Email).Equals("")
}
