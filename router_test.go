package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
	RegisterTestingT(t)

	ctx := &WechyServices{}
	r := GetMainEngine(ctx)
	Expect(r).ToNot(BeNil())
}
