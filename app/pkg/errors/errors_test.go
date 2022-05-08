package errors_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/errors"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestNewError(t *testing.T) {
	RegisterT(t)

	first := errors.New("document not found")
	Expect(first.Error()).Equals("document not found (app/pkg/errors/errors_test.go:14)")
}

func TestWrappedError(t *testing.T) {
	RegisterT(t)

	first := errors.New("document not found")
	wrapped := errors.Wrap(first, "could not create user '%s'", "John")
	wrappedAgain := errors.Wrap(wrapped, "failed to register new user")

	Expect(first.Error()).Equals(`document not found (app/pkg/errors/errors_test.go:21)`)
	Expect(wrapped.Error()).Equals(`Error Trace: 
- could not create user 'John' (app/pkg/errors/errors_test.go:22)
- document not found (app/pkg/errors/errors_test.go:21)`)
	Expect(wrappedAgain.Error()).Equals(`Error Trace: 
- failed to register new user (app/pkg/errors/errors_test.go:23)
- could not create user 'John' (app/pkg/errors/errors_test.go:22)
- document not found (app/pkg/errors/errors_test.go:21)`)

	Expect(errors.Cause(wrapped)).Equals(first)
	Expect(errors.Cause(first)).Equals(first)
	Expect(errors.Cause(wrappedAgain)).Equals(first)
}

func TestWrapAndStack(t *testing.T) {
	RegisterT(t)

	first := errors.New("document not found")
	wrapped := errors.Wrap(first, "could not create user")
	stacked := errors.Stack(wrapped)
	stackedTwice := errors.Stack(stacked)

	Expect(stackedTwice.Error()).Equals(`Error Trace: 
- app/pkg/errors/errors_test.go:45
- app/pkg/errors/errors_test.go:44
- could not create user (app/pkg/errors/errors_test.go:43)
- document not found (app/pkg/errors/errors_test.go:42)`)

	Expect(errors.Cause(stacked)).Equals(first)
	Expect(errors.Cause(stackedTwice)).Equals(first)
}

func doSomething(err error) error {
	return errors.StackN(err, 1)
}

func TestStackN(t *testing.T) {
	RegisterT(t)

	first := errors.New("document not found")
	err := doSomething(first)
	Expect(err.Error()).Equals(`Error Trace: 
- app/pkg/errors/errors_test.go:65
- document not found (app/pkg/errors/errors_test.go:64)`)

	Expect(errors.Cause(err)).Equals(first)
}

func TestNilErrors(t *testing.T) {
	RegisterT(t)

	Expect(errors.Cause(nil)).IsNil()
	Expect(errors.Stack(nil)).IsNil()
	Expect(errors.Wrap(nil, "")).IsNil()
}

func TestPanicked(t *testing.T) {
	RegisterT(t)

	defer func() {
		if r := recover(); r != nil {
			err := errors.Panicked(r)
			Expect(err.Error()).ContainsSubstring(`runtime/debug.Stack`)
			Expect(err.Error()).ContainsSubstring(`github.com/getfider/fider/app/pkg/errors.Panicked`)
			Expect(err.Error()).ContainsSubstring(`- github.com/getfider/fider/app/pkg/errors_test.TestPanicked:`)
			Expect(err.Error()).ContainsSubstring(`(app/pkg/errors/errors.go:39)
- Boom!`)
		}
	}()

	panic("Boom!")
}
