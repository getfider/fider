package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestTagStorage_AddAndGet(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	tags := postgres.NewTagStorage(demoTenant(tenants), trx)
	tag, err := tags.Add("Feature Request", "FF0000", true)
	Expect(err).To(BeNil())
	Expect(tag.ID).To(Equal(1))

	dbTag, err := tags.GetBySlug(tag.Slug)

	Expect(err).To(BeNil())
	Expect(dbTag.ID).To(Equal(1))
	Expect(dbTag.Name).To(Equal("Feature Request"))
	Expect(dbTag.Slug).To(Equal("feature-request"))
	Expect(dbTag.Color).To(Equal("FF0000"))
	Expect(dbTag.IsPublic).To(BeTrue())
}
