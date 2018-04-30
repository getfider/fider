package cmd

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetMainEngine(t *testing.T) {
	RegisterTestingT(t)

	r := getMainEngine(nil)
	Expect(r).ToNot(BeNil())
}
