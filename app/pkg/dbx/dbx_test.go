package dbx_test

import (
	"os"
	"testing"

	"github.com/getfider/fider/app/pkg/dbx"
	. "github.com/onsi/gomega"
)

type user struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	IgnoreThis string
}

type userWithTenant struct {
	ID     int     `db:"id"`
	Name   string  `db:"name"`
	Tenant *tenant `db:"tenant"`
}

type userProvider struct {
	Provider string          `db:"provider"`
	User     *userWithTenant `db:"user"`
}

type tenant struct {
	ID   int     `db:"id"`
	Name *string `db:"name"`
}

func TestMain(m *testing.M) {
	db, _ := dbx.New()
	db.Migrate()
	db.Seed()
	db.Close()

	code := m.Run()
	os.Exit(code)
}

func TestBind_SimpleStruct(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()
	u := user{}

	err := trx.Get(&u, "SELECT id, name FROM users LIMIT 1")
	Expect(err).To(BeNil())
	Expect(u.ID).To(Equal(300))
	Expect(u.Name).To(Equal("Jon Snow"))
	Expect(u.IgnoreThis).To(Equal(""))
}

func TestBind_DeepNestedStruct(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()
	u := userProvider{}

	err := trx.Get(&u, `SELECT up.provider,
														up.user_id,
														u.tenant_id AS user_tenant_id
										 FROM user_providers up
										 INNER JOIN users u
										 ON up.user_id = u.id
										 WHERE provider_uid = 'FB2222'`)
	Expect(err).To(BeNil())
	Expect(u.Provider).To(Equal("facebook"))
	Expect(u.User.ID).To(Equal(400))
	Expect(u.User.Tenant.ID).To(Equal(400))
}

func TestBind_SimpleStruct_SingleField(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()
	u := user{}

	err := trx.Get(&u, "SELECT name FROM users LIMIT 1")
	Expect(err).To(BeNil())
	Expect(u.ID).To(Equal(0))
	Expect(u.Name).To(Equal("Jon Snow"))
}

func TestBind_SimpleStruct_Multiple(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()
	u := []*user{}

	err := trx.Select(&u, "SELECT name FROM users WHERE tenant_id = 300")
	Expect(err).To(BeNil())

	Expect(len(u)).To(Equal(2))
	Expect(u[0].Name).To(Equal("Jon Snow"))
	Expect(u[1].Name).To(Equal("Arya Stark"))
}

func TestBind_NestedStruct(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	u := userWithTenant{}

	err := trx.Get(&u, `
		SELECT u.id, u.name, t.id AS tenant_id, t.name AS tenant_name
		FROM users u
		INNER JOIN tenants t
		ON t.id = u.tenant_id
		WHERE u.id = 300
		LIMIT 1
	`)
	Expect(err).To(BeNil())
	Expect(u.ID).To(Equal(300))
	Expect(u.Name).To(Equal("Jon Snow"))
	Expect(u.Tenant).NotTo(BeNil())
	Expect(u.Tenant.ID).To(Equal(300))
	Expect(*u.Tenant.Name).To(Equal("Demonstration"))
}

func TestBind_NestedStruct_Multiple(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	u := []*userWithTenant{}

	err := trx.Select(&u, `
		SELECT u.id, u.name, t.id AS tenant_id, t.name AS tenant_name
		FROM users u
		INNER JOIN tenants t
		ON t.id = u.tenant_id
		WHERE u.tenant_id = 300
	`)
	Expect(err).To(BeNil())
	Expect(len(u)).To(Equal(2))
	Expect(u[0].Name).To(Equal("Jon Snow"))
	Expect(*u[0].Tenant.Name).To(Equal("Demonstration"))
	Expect(u[1].Name).To(Equal("Arya Stark"))
	Expect(*u[1].Tenant.Name).To(Equal("Demonstration"))

}

func TestExists_True(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	exists, err := trx.Exists("SELECT 1 FROM users WHERE id = 300")
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
}

func TestExists_False(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	exists, err := trx.Exists("SELECT 1 FROM users WHERE id = 0")
	Expect(err).To(BeNil())
	Expect(exists).To(BeFalse())
}

func TestCount(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	count, err := trx.Count("SELECT 1 FROM users WHERE id = 300")
	Expect(err).To(BeNil())
	Expect(count).To(Equal(1))
}

func TestCount_Empty(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	count, err := trx.Count("SELECT 1 FROM users WHERE id = 0")
	Expect(err).To(BeNil())
	Expect(count).To(Equal(0))
}
