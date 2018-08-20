package rand_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/rand"
)

func TestRandomString(t *testing.T) {
	RegisterT(t)

	Expect(rand.String(10000)).HasLen(10000)
	Expect(rand.String(10)).HasLen(10)
	Expect(rand.String(0)).HasLen(0)
	Expect(rand.String(-1)).HasLen(0)
}
