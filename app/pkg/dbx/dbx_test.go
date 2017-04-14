package dbx_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	. "github.com/onsi/gomega"
)

type user struct {
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

	u := []user{}

	err := db.Select(&u, "SELECT name FROM users WHERE tenant_id = 300")
	Expect(err).To(BeNil())

	Expect(len(u)).To(Equal(2))
	Expect(u[0].Name).To(Equal("Jon Snow"))
	Expect(u[1].Name).To(Equal("Arya Stark"))
}
