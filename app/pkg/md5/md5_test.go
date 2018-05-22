package md5_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/md5"
)

func TestMD5Hash(t *testing.T) {
	RegisterT(t)

	hash := md5.Hash("Fider")

	Expect(hash).Equals("3734538c8b650e4f354a55a436566bb6")
}
