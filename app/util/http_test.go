package util_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/app/util"
	. "github.com/onsi/gomega"
)

func TestStripPort(t *testing.T) {
	RegisterTestingT(t)

	Expect(util.StripPort("mydomain.com:5000")).To(Equal("mydomain.com"))
	Expect(util.StripPort("mydomain.com")).To(Equal("mydomain.com"))
}
