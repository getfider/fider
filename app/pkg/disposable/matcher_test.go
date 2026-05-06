package disposable_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"

	"github.com/getfider/fider/app/pkg/disposable"
)

func TestIsBlocked_BundledOnly(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBlocked("foo@mailinator.com", nil, nil)).IsTrue()
	Expect(disposable.IsBlocked("foo@temp.mailinator.com", nil, nil)).IsTrue()
	Expect(disposable.IsBlocked("foo@gmail.com", nil, nil)).IsFalse()
}

func TestIsBlocked_TenantDeny(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBlocked("foo@example.com", []string{"example.com"}, nil)).IsTrue()
	Expect(disposable.IsBlocked("foo@x.example.com", []string{"example.com"}, nil)).IsTrue()
	Expect(disposable.IsBlocked("foo@example.org", []string{"example.com"}, nil)).IsFalse()
}

func TestIsBlocked_TenantAllowOverridesBundled(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBlocked("foo@mailinator.com", nil, []string{"mailinator.com"})).IsFalse()
	Expect(disposable.IsBlocked("foo@temp.mailinator.com", nil, []string{"mailinator.com"})).IsFalse()
}

func TestIsBlocked_TenantAllowOverridesTenantDeny(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBlocked("foo@example.com",
		[]string{"example.com"}, []string{"example.com"})).IsFalse()
}

func TestIsBlocked_CaseInsensitive(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBlocked("Foo@MAILINATOR.com", nil, nil)).IsTrue()
	Expect(disposable.IsBlocked("Foo@Example.COM", []string{"EXAMPLE.com"}, nil)).IsTrue()
}

func TestIsBlocked_MalformedEmail(t *testing.T) {
	RegisterT(t)
	Expect(disposable.IsBlocked("not-an-email", nil, nil)).IsFalse()
	Expect(disposable.IsBlocked("@mailinator.com", nil, nil)).IsTrue()
	Expect(disposable.IsBlocked("foo@", nil, nil)).IsFalse()
	Expect(disposable.IsBlocked("", nil, nil)).IsFalse()
}
