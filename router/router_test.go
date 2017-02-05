package router

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
	RegisterTestingT(t)

	r := GetMainEngine(nil)
	Expect(r).ToNot(BeNil())
}

func TestStripPort(t *testing.T) {
	RegisterTestingT(t)

	Expect(stripPort("localhost:3000")).To(Equal("localhost"))
	Expect(stripPort("mypage.com:3000")).To(Equal("mypage.com"))
	Expect(stripPort("sub.mypage.com:3000")).To(Equal("sub.mypage.com"))
	Expect(stripPort("mypage.com")).To(Equal("mypage.com"))
}
