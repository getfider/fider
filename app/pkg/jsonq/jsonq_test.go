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

func TestContains(t *testing.T) {
	RegisterTestingT(t)

	query := jsonq.New(`{ "name": "Jon Snow" }`)
	Expect(query.Contains("name")).To(BeTrue())
	Expect(query.Contains("what")).To(BeFalse())
}

func TestContainsNested(t *testing.T) {
	RegisterTestingT(t)

	query := jsonq.New(`{ "failures": { "name": "Name is required" } }`)
	Expect(query.Contains("failures.name")).To(BeTrue())
	Expect(query.Contains("failures.what")).To(BeFalse())
}
