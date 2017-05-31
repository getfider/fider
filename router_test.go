package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetMainEngine(t *testing.T) {
	RegisterTestingT(t)

	r := GetMainEngine(nil, nil)
	Expect(r).ToNot(BeNil())
}
