package jsonq_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/jsonq"
)

func TestGet(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "name": "Jon Snow", "age": 23 }`)
	Expect(query.String("name")).Equals("Jon Snow")
	Expect(query.Int32("age")).Equals(23)
}

func TestGetNull(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "name": null }`)
	Expect(query.String("name")).Equals("")
}

func TestGet_NestedObject(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "person": { "name": "Jon Snow", "age": 23 } }`)
	Expect(query.String("person.name")).Equals("Jon Snow")
	Expect(query.Int32("person.age")).Equals(23)
}

func TestGet_EmptyKey(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "name": "Jon Snow", "age": 23 }`)
	Expect(query.String("")).Equals("")
	Expect(query.Int32("")).Equals(0)
}

func TestGet_NumberAsString(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "age": 23 }`)
	Expect(query.String("age")).Equals("23")
}

func TestGetStringNested(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "failures": { "name": "Jon Snow" } }`)
	Expect(query.String("failures.name")).Equals("Jon Snow")
}

func TestGetWithFallback(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "name": "", "login": "jonsnow" }`)
	Expect(query.String("login")).Equals("jonsnow")
	Expect(query.String("name")).Equals("")
	Expect(query.String("name, login")).Equals("jonsnow")
}

func TestGetValueFromObjectArray(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "data": [{ "name": "Jon Snow" }, { "age": 23 }] }`)
	Expect(query.String("data[0].name")).Equals("Jon Snow")
	Expect(query.String("data[0].age")).Equals("")
	Expect(query.String("data[1].age")).Equals("23")
}

func TestGetValueFromNestedArray(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "data": [ { "numbers": [52, 6, 24] } ] }`)
	Expect(query.Int32("data[0].numbers[0]")).Equals(52)
	Expect(query.Int32("data[0].numbers[1]")).Equals(6)
	Expect(query.Int32("data[0].numbers[2]")).Equals(24)
}

func TestGetValueFromNestedObjectArray(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "data": [{ "person": { "name": "Jon Snow" } }, { "person": { "age": "23" } }] }`)
	Expect(query.String("data[0].person.name")).Equals("Jon Snow")
	Expect(query.String("data[0].person.age")).Equals("")
	Expect(query.String("data[1].person.age")).Equals("23")
}

func TestGetValueFromStringArray(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "data": ["Jon Snow", "Arya Stark"] }`)
	Expect(query.String("data[0]")).Equals("Jon Snow")
	Expect(query.String("data[1]")).Equals("Arya Stark")
	Expect(query.String("data[2]")).Equals("")
}

func TestGetValueFromEmptyArray(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "data": [] }`)
	Expect(query.String("unknown[0].name")).Equals("")
	Expect(query.String("data[0].name")).Equals("")
	Expect(query.String("data[0].age")).Equals("")
}

func TestContains(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "name": "Jon Snow" }`)
	Expect(query.IsArray()).IsFalse()
	Expect(query.Contains("name")).IsTrue()
	Expect(query.Contains("what")).IsFalse()
	Expect(query.Contains("feature.name")).IsFalse()
}

func TestIsArray(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`[0,1,2,3]`)
	Expect(query.IsArray()).IsTrue()
	Expect(query.ArrayLength()).Equals(4)
}

func TestContainsNested(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "failures": { "name": "Name is required" } }`)
	Expect(query.IsArray()).IsFalse()

	Expect(query.Contains("failures")).IsTrue()
	Expect(query.Contains("failures.name")).IsTrue()

	Expect(query.Contains("name")).IsFalse()
	Expect(query.Contains("failures.what")).IsFalse()
}

func TestStrings_ArrayOfStrings(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "roles": ["ROLE_ADMIN", "ROLE_USER"] }`)
	Expect(query.Strings("roles")).Equals([]string{"ROLE_ADMIN", "ROLE_USER"})
}

func TestStrings_SingleString(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "role": "ROLE_ADMIN" }`)
	Expect(query.Strings("role")).Equals([]string{"ROLE_ADMIN"})
}

func TestStrings_NestedPath(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "user": { "profile": { "roles": ["ROLE_ADMIN"] } } }`)
	Expect(query.Strings("user.profile.roles")).Equals([]string{"ROLE_ADMIN"})
}

func TestStrings_MissingSelectorReturnsNil(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "name": "Jon" }`)
	Expect(query.Strings("roles")).IsNil()
}

func TestStrings_EmptyArrayReturnsNil(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "roles": [] }`)
	Expect(query.Strings("roles")).IsNil()
}

func TestArrayFieldStrings_ExtractField(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "roles": [{"id": "ROLE_ADMIN", "name": "Admin"}, {"id": "ROLE_USER", "name": "User"}] }`)
	Expect(query.ArrayFieldStrings("roles", "id")).Equals([]string{"ROLE_ADMIN", "ROLE_USER"})
}

func TestArrayFieldStrings_NestedPath(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "user": { "groups": [{"name": "ROLE_ADMIN"}, {"name": "ROLE_USER"}] } }`)
	Expect(query.ArrayFieldStrings("user.groups", "name")).Equals([]string{"ROLE_ADMIN", "ROLE_USER"})
}

func TestArrayFieldStrings_MissingSelectorReturnsNil(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "name": "Jon" }`)
	Expect(query.ArrayFieldStrings("roles", "id")).IsNil()
}

func TestArrayFieldStrings_NotArrayReturnsNil(t *testing.T) {
	RegisterT(t)

	query := jsonq.New(`{ "roles": "not an array" }`)
	Expect(query.ArrayFieldStrings("roles", "id")).IsNil()
}
