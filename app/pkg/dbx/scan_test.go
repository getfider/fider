package dbx

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

type tenantWithStringPointer struct {
	ID   int        `db:"id"`
	Name *string    `db:"name"`
	Time *time.Time `db:"time"`
}

type tenant struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type user struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Tenant     *tenant `db:"tenant"`
	IgnoreThis bool
}

type userProvider struct {
	Name string `db:"provider"`
	User *user  `db:"user"`
}

func TestScan(t *testing.T) {
	RegisterTestingT(t)

	u := user{}
	result := scan("", &u)
	Expect(len(result)).To(Equal(4))
	Expect(result["id"]).To(Equal(&u.ID))
	Expect(result["name"]).To(Equal(&u.Name))
	Expect(result["tenant_id"]).To(Equal(&u.Tenant.ID))
	Expect(result["tenant_name"]).To(Equal(&u.Tenant.Name))
}

func TestScan_WithStringPointer(t *testing.T) {
	RegisterTestingT(t)

	u := tenantWithStringPointer{}
	result := scan("", &u)
	Expect(len(result)).To(Equal(3))
	Expect(result["id"]).To(Equal(&u.ID))
	Expect(result["name"]).To(Equal(u.Name))
	Expect(result["time"]).To(Equal(u.Time))
}

func TestDeepScan(t *testing.T) {
	RegisterTestingT(t)

	u := userProvider{}
	result := scan("", &u)
	Expect(len(result)).To(Equal(5))
	Expect(result["provider"]).To(Equal(&u.Name))
	Expect(result["user_id"]).To(Equal(&u.User.ID))
	Expect(result["user_name"]).To(Equal(&u.User.Name))
	Expect(result["user_tenant_id"]).To(Equal(&u.User.Tenant.ID))
	Expect(result["user_tenant_name"]).To(Equal(&u.User.Tenant.Name))
}
