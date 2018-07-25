package validate_test

import (
	"errors"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/validate"
)

func TestResult_Error(t *testing.T) {
	RegisterT(t)

	err := errors.New("Something went wrong...")
	r := validate.Error(err)
	Expect(r.Ok).IsFalse()
	Expect(r.Authorized).IsTrue()
	Expect(r.Err).Equals(err)
}

func TestResult_AddFieldFailure_Empty(t *testing.T) {
	RegisterT(t)

	r := validate.Success()
	r.AddFieldFailure("name")
	Expect(r.Ok).IsTrue()
	Expect(r.Failures["name"]).HasLen(0)

	r.AddFieldFailure("name", "This field is required")
	Expect(r.Ok).IsFalse()
	Expect(r.Failures["name"]).HasLen(1)
}
