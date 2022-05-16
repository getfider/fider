package dbx_test

import (
	"reflect"
	"sync"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
)

func createScanner(values ...any) func(dest ...any) error {
	return func(dest ...any) error {
		for i := 0; i < len(values); i++ {
			reflect.Indirect(reflect.ValueOf(dest[i])).Set(reflect.ValueOf(values[i]))
		}
		return nil
	}
}

func TestTypeMapper_NonStruct(t *testing.T) {
	RegisterT(t)
	var id int
	mapping := dbx.NewTypeMapper(reflect.TypeOf(id))
	Expect(mapping.Fields).HasLen(0)
}

func TestTypeMapper_SimpleStruct(t *testing.T) {
	RegisterT(t)
	u := user{}
	mapping := dbx.NewTypeMapper(reflect.TypeOf(u))
	Expect(mapping.Fields).HasLen(2)
	Expect(mapping.Fields["id"].FieldName).Equals([]string{"ID"})
	Expect(mapping.Fields["name"].FieldName).Equals([]string{"Name"})
}

func TestTypeMapper_StringPointer(t *testing.T) {
	RegisterT(t)
	m := tenant{}
	mapping := dbx.NewTypeMapper(reflect.TypeOf(m))
	Expect(mapping.Fields).HasLen(2)
	Expect(mapping.Fields["id"].FieldName).Equals([]string{"ID"})
	Expect(mapping.Fields["name"].FieldName).Equals([]string{"Name"})
}

func TestTypeMapper_NestedStruct(t *testing.T) {
	RegisterT(t)
	u := userWithTenant{}
	mapping := dbx.NewTypeMapper(reflect.TypeOf(u))
	Expect(mapping.Fields).HasLen(4)
	Expect(mapping.Fields["id"].FieldName).Equals([]string{"ID"})
	Expect(mapping.Fields["name"].FieldName).Equals([]string{"Name"})
	Expect(mapping.Fields["tenant_id"].FieldName).Equals([]string{"Tenant", "ID"})
	Expect(mapping.Fields["tenant_name"].FieldName).Equals([]string{"Tenant", "Name"})
}

func TestTypeMapper_DeepNestedStruct(t *testing.T) {
	RegisterT(t)
	u := userProvider{}
	mapping := dbx.NewTypeMapper(reflect.TypeOf(u))
	Expect(mapping.Fields).HasLen(5)
	Expect(mapping.Fields["provider"].FieldName).Equals([]string{"Provider"})
	Expect(mapping.Fields["user_id"].FieldName).Equals([]string{"User", "ID"})
	Expect(mapping.Fields["user_name"].FieldName).Equals([]string{"User", "Name"})
	Expect(mapping.Fields["user_tenant_id"].FieldName).Equals([]string{"User", "Tenant", "ID"})
	Expect(mapping.Fields["user_tenant_name"].FieldName).Equals([]string{"User", "Tenant", "Name"})
}

func TestRowMapper_NonStruct(t *testing.T) {
	RegisterT(t)
	mapper := dbx.NewRowMapper()
	var id int
	err := mapper.Map(&id, []string{"id"}, createScanner(5))
	Expect(err).IsNil()
	Expect(id).Equals(5)
}

func TestRowMapper_Concurrent(t *testing.T) {
	RegisterT(t)
	var wg sync.WaitGroup
	mapper := dbx.NewRowMapper()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			var id int
			err := mapper.Map(&id, []string{"id"}, createScanner(5))
			Expect(err).IsNil()
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestRowMapper_SimpleStruct(t *testing.T) {
	RegisterT(t)
	mapper := dbx.NewRowMapper()
	u := user{}
	err := mapper.Map(&u, []string{"id", "name"}, createScanner(1, "Hello World"))
	Expect(err).IsNil()
	Expect(u.ID).Equals(1)
	Expect(u.Name).Equals("Hello World")
}

func TestRowMapper_NestedStruct(t *testing.T) {
	RegisterT(t)
	mapper := dbx.NewRowMapper()
	u := userWithTenant{}
	err := mapper.Map(&u, []string{"id", "name", "tenant_id", "tenant_name"}, createScanner(1, "Hello World", 2, "Demonstration"))
	Expect(err).IsNil()
	Expect(u.ID).Equals(1)
	Expect(u.Name).Equals("Hello World")
	Expect(u.Tenant.ID).Equals(2)
	Expect(u.Tenant.Name).Equals("Demonstration")
}

func TestRowMapper_DeepNestedStruct(t *testing.T) {
	RegisterT(t)
	mapper := dbx.NewRowMapper()
	u := userProvider{}
	err := mapper.Map(&u, []string{"provider", "user_id", "user_name", "user_tenant_id", "user_tenant_name"}, createScanner("google", 1, "Jon Snow", 2, "The Tenant"))
	Expect(err).IsNil()
	Expect(u.Provider).Equals("google")
	Expect(u.User.ID).Equals(1)
	Expect(u.User.Name).Equals("Jon Snow")
	Expect(u.User.Tenant.ID).Equals(2)
	Expect(u.User.Tenant.Name).Equals("The Tenant")
}
