package util_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/app/util"
	. "github.com/onsi/gomega"
)

func TestMD5(t *testing.T) {
	RegisterTestingT(t)

	encrypted := util.MD5("jonsnow")

	Expect(encrypted).To(Equal("5a665206e6374a1b3b95e05d1ae9ecd8"))
}
