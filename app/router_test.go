package app_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app"
	. "github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
	RegisterTestingT(t)

	ctx := &app.WechyServices{}
	r := app.GetMainEngine(ctx)
	Expect(r).ToNot(BeNil())
}
