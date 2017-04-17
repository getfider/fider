package dbx_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
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

type tenant struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func TestBind_SimpleStruct(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()

	u := user{}

	err := db.Get(&u, "SELECT id, name FROM users LIMIT 1")
	Expect(err).To(BeNil())
	Expect(u.ID).To(Equal(300))
	Expect(u.Name).To(Equal("Jon Snow"))
	Expect(u.IgnoreThis).To(Equal(""))
}

func TestBind_SimpleStruct_SingleField(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()

	u := user{}

	err := db.Get(&u, "SELECT name FROM users LIMIT 1")
	Expect(err).To(BeNil())
	Expect(u.ID).To(Equal(0))
	Expect(u.Name).To(Equal("Jon Snow"))
}

func TestBind_SimpleStruct_Multiple(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()

	u := []*user{}

	err := db.Select(&u, "SELECT name FROM users WHERE tenant_id = 300")
	Expect(err).To(BeNil())

	Expect(len(u)).To(Equal(2))
	Expect(u[0].Name).To(Equal("Jon Snow"))
	Expect(u[1].Name).To(Equal("Arya Stark"))
}

func TestBind_NestedStruct(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()

	u := userWithTenant{}

	err := db.Get(&u, `
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
	Expect(u.Tenant.Name).To(Equal("Demonstration"))
}

func TestBind_NestedStruct_Multiple(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()

	u := []*userWithTenant{}

	err := db.Select(&u, `
		SELECT u.id, u.name, t.id AS tenant_id, t.name AS tenant_name
		FROM users u
		INNER JOIN tenants t
		ON t.id = u.tenant_id
		WHERE u.tenant_id = 300
	`)
	Expect(err).To(BeNil())
	Expect(len(u)).To(Equal(2))
	Expect(u[0].Name).To(Equal("Jon Snow"))
	Expect(u[0].Tenant.Name).To(Equal("Demonstration"))
	Expect(u[1].Name).To(Equal("Arya Stark"))
	Expect(u[1].Tenant.Name).To(Equal("Demonstration"))
}
