package router_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/router"
	. "github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
	RegisterTestingT(t)

	r := router.GetMainEngine(nil)
	Expect(r).ToNot(BeNil())
}
