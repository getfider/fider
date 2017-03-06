package app_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/app"
	. "github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
	RegisterTestingT(t)

	ctx := &app.WchyServices{}
	r := app.GetMainEngine(ctx)
	Expect(r).ToNot(BeNil())
}
