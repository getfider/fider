package errors_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/errors"

	. "github.com/onsi/gomega"
)

func TestNewError(t *testing.T) {
	RegisterTestingT(t)

	first := errors.New("document not found")
	Expect(first.Error()).To(Equal("document not found (app/pkg/errors/errors_test.go:14)"))
}

func TestWrappedError(t *testing.T) {
	RegisterTestingT(t)

	first := errors.New("document not found")
	wrapped := errors.Wrap(first, "could not create user")
	wrappedAgain := errors.Wrap(wrapped, "failed to register new user")

	Expect(first.Error()).To(Equal(`document not found (app/pkg/errors/errors_test.go:21)`))
	Expect(wrapped.Error()).To(Equal(`Error Trace: 
- could not create user (app/pkg/errors/errors_test.go:22)
- document not found (app/pkg/errors/errors_test.go:21)`))
	Expect(wrappedAgain.Error()).To(Equal(`Error Trace: 
- failed to register new user (app/pkg/errors/errors_test.go:23)
- could not create user (app/pkg/errors/errors_test.go:22)
- document not found (app/pkg/errors/errors_test.go:21)`))

	Expect(errors.Cause(wrapped)).To(Equal(first))
	Expect(errors.Cause(first)).To(Equal(first))
	Expect(errors.Cause(wrappedAgain)).To(Equal(first))
}

func TestNilErrors(t *testing.T) {
	RegisterTestingT(t)

	Expect(errors.Cause(nil)).To(BeNil())
}
