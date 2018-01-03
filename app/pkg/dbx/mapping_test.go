package dbx_test

import (
	"reflect"
	"testing"

	"github.com/getfider/fider/app/pkg/dbx"
	. "github.com/onsi/gomega"
)

func createScanner(values ...interface{}) func(dest ...interface{}) error {
	return func(dest ...interface{}) error {
		for i := 0; i < len(values); i++ {
			reflect.Indirect(reflect.ValueOf(dest[i])).Set(reflect.ValueOf(values[i]))
		}
		return nil
	}
}

func TestTypeMapper_SimpleStruct(t *testing.T) {
	RegisterTestingT(t)
	u := user{}
	mapping := dbx.NewTypeMapper(reflect.TypeOf(u))
	Expect(len(mapping.Fields)).To(Equal(2))
	Expect(mapping.Fields["id"].FieldName).To(Equal([]string{"ID"}))
	Expect(mapping.Fields["name"].FieldName).To(Equal([]string{"Name"}))
}

func TestTypeMapper_StringPointer(t *testing.T) {
	RegisterTestingT(t)
	m := tenant{}
	mapping := dbx.NewTypeMapper(reflect.TypeOf(m))
	Expect(len(mapping.Fields)).To(Equal(2))
	Expect(mapping.Fields["id"].FieldName).To(Equal([]string{"ID"}))
	Expect(mapping.Fields["name"].FieldName).To(Equal([]string{"Name"}))
}

func TestTypeMapper_NestedStruct(t *testing.T) {
	RegisterTestingT(t)
	u := userWithTenant{}
	mapping := dbx.NewTypeMapper(reflect.TypeOf(u))
	Expect(len(mapping.Fields)).To(Equal(4))
	Expect(mapping.Fields["id"].FieldName).To(Equal([]string{"ID"}))
	Expect(mapping.Fields["name"].FieldName).To(Equal([]string{"Name"}))
	Expect(mapping.Fields["tenant_id"].FieldName).To(Equal([]string{"Tenant", "ID"}))
	Expect(mapping.Fields["tenant_name"].FieldName).To(Equal([]string{"Tenant", "Name"}))
}

func TestTypeMapper_DeepNestedStruct(t *testing.T) {
	RegisterTestingT(t)
	u := userProvider{}
	mapping := dbx.NewTypeMapper(reflect.TypeOf(u))
	Expect(len(mapping.Fields)).To(Equal(5))
	Expect(mapping.Fields["provider"].FieldName).To(Equal([]string{"Provider"}))
	Expect(mapping.Fields["user_id"].FieldName).To(Equal([]string{"User", "ID"}))
	Expect(mapping.Fields["user_name"].FieldName).To(Equal([]string{"User", "Name"}))
	Expect(mapping.Fields["user_tenant_id"].FieldName).To(Equal([]string{"User", "Tenant", "ID"}))
	Expect(mapping.Fields["user_tenant_name"].FieldName).To(Equal([]string{"User", "Tenant", "Name"}))
}

func TestRowMapper_SimpleStruct(t *testing.T) {
	RegisterTestingT(t)
	mapper := dbx.NewRowMapper()
	u := user{}
	mapper.Map(&u, []string{"id", "name"}, createScanner(1, "Hello World"))
	Expect(u.ID).To(Equal(1))
	Expect(u.Name).To(Equal("Hello World"))
}

func TestRowMapper_NestedStruct(t *testing.T) {
	RegisterTestingT(t)
	mapper := dbx.NewRowMapper()
	u := userWithTenant{}
	mapper.Map(&u, []string{"id", "name", "tenant_id", "tenant_name"}, createScanner(1, "Hello World", 2, "Demonstration"))
	Expect(u.ID).To(Equal(1))
	Expect(u.Name).To(Equal("Hello World"))
	Expect(u.Tenant.ID).To(Equal(2))
	Expect(u.Tenant.Name).To(Equal("Demonstration"))
}

func TestRowMapper_DeepNestedStruct(t *testing.T) {
	RegisterTestingT(t)
	mapper := dbx.NewRowMapper()
	u := userProvider{}
	mapper.Map(&u, []string{"provider", "user_id", "user_name", "user_tenant_id", "user_tenant_name"}, createScanner("google", 1, "Jon Snow", 2, "The Tenant"))
	Expect(u.Provider).To(Equal("google"))
	Expect(u.User.ID).To(Equal(1))
	Expect(u.User.Name).To(Equal("Jon Snow"))
	Expect(u.User.Tenant.ID).To(Equal(2))
	Expect(u.User.Tenant.Name).To(Equal("The Tenant"))
}
