package cmd

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func TestGetMainEngine(t *testing.T) {
	RegisterT(t)

	r := routes(web.New())
	Expect(r).IsNotNil()
}
