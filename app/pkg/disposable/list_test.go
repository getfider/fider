package disposable_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"

	"github.com/getfider/fider/app/pkg/disposable"
)

func TestIsBundled_ExactMatch(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBundled("mailinator.com")).IsTrue()
	Expect(disposable.IsBundled("gmail.com")).IsFalse()
}

func TestIsBundled_SubdomainMatch(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBundled("foo.mailinator.com")).IsTrue()
	Expect(disposable.IsBundled("a.b.c.mailinator.com")).IsTrue()
}

func TestIsBundled_CaseInsensitive(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBundled("MAILINATOR.COM")).IsTrue()
	Expect(disposable.IsBundled("Foo.Mailinator.Com")).IsTrue()
}

func TestIsBundled_EmptyAndMalformed(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBundled("")).IsFalse()
	Expect(disposable.IsBundled(".")).IsFalse()
	Expect(disposable.IsBundled(".com")).IsFalse()
	Expect(disposable.IsBundled("foo..mailinator.com")).IsFalse()
}

func TestIsBundled_TrailingDot(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBundled("mailinator.com.")).IsTrue()
	Expect(disposable.IsBundled("foo.mailinator.com.")).IsTrue()
}
