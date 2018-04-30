package cmd

import (
	"testing"

	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

func TestGetMainEngine(t *testing.T) {
	RegisterTestingT(t)

	r := routes(web.New(nil))
	Expect(r).ToNot(BeNil())
}
