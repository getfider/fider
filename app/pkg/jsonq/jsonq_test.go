package jsonq_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/jsonq"
	. "github.com/onsi/gomega"
)

func TestGetString(t *testing.T) {
	RegisterTestingT(t)

	query := jsonq.New(`{ "name": "Jon Snow" }`)
	Expect(query.String("name")).To(Equal("Jon Snow"))
}

func TestGetStringNested(t *testing.T) {
	RegisterTestingT(t)

	query := jsonq.New(`{ "failures": { "name": "Jon Snow" } }`)
	Expect(query.String("failures.name")).To(Equal("Jon Snow"))
}

func TestContains(t *testing.T) {
	RegisterTestingT(t)

	query := jsonq.New(`{ "name": "Jon Snow" }`)
	Expect(query.IsArray()).To(BeFalse())
	Expect(query.Contains("name")).To(BeTrue())
	Expect(query.Contains("what")).To(BeFalse())
	Expect(query.Contains("feature.name")).To(BeFalse())
}

func TestIsArray(t *testing.T) {
	RegisterTestingT(t)

	query := jsonq.New(`[0,1,2,3]`)
	Expect(query.IsArray()).To(BeTrue())
}

func TestContainsNested(t *testing.T) {
	RegisterTestingT(t)

	query := jsonq.New(`{ "failures": { "name": "Name is required" } }`)
	Expect(query.IsArray()).To(BeFalse())

	Expect(query.Contains("failures")).To(BeTrue())
	Expect(query.Contains("failures.name")).To(BeTrue())

	Expect(query.Contains("name")).To(BeFalse())
	Expect(query.Contains("failures.what")).To(BeFalse())
}
